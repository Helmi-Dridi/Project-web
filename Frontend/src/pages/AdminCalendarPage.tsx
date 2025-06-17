import { useState, useEffect } from "react";
import {
  Calendar as BigCalendar,
  dateFnsLocalizer,
  type View,
  type ToolbarProps,
} from "react-big-calendar";
import { format, parse, startOfWeek, getDay } from "date-fns";
import { enUS } from "date-fns/locale";
import "react-big-calendar/lib/css/react-big-calendar.css";

import { ChevronLeft, ChevronRight } from "lucide-react";

import DailyAppointmentsModal from "../components/DailyAppointmentsModal";
import UpdateAppointmentModal from "../components/UpdateAppointmentModal";

// Summary:
// - Fixed typing by introducing a CalendarEvent interface for the calendar.
// - Added view and date state to enable navigation between time ranges.
// - Implemented a custom toolbar with accessible buttons and tooltips.
// - Calendar controls (Today, Back, Next, Month, Week, Day, Agenda) now work
//   properly and the layout is more responsive.

import {
  useUpcomingAppointments,
  useUpdateAppointmentStatus,
  useAppointmentStats,
  useSendReminder,
  useRescheduleAppointment,
} from "../hooks/useBooking";
import type { Appointment } from "../services/booking.service";

// Localization
const locales = { "en-US": enUS };
const localizer = dateFnsLocalizer({
  format,
  parse,
  startOfWeek,
  getDay,
  locales,
});



interface CalendarEvent {
  id: string;
  title: string;
  start: Date;
  end: Date;
  resource: Appointment;
}

function CalendarToolbar({ label, onNavigate, onView, view }: ToolbarProps<CalendarEvent, Appointment>) {
  const viewNames: View[] = ["month", "week", "day", "agenda"];
  return (
    <div className="flex flex-wrap items-center justify-between gap-2 mb-4">
      <div className="flex items-center gap-2">
        <button
          title="Today"
          aria-label="Go to today"
          onClick={() => onNavigate("TODAY")}
          className="px-3 py-1 rounded bg-gray-100 hover:bg-gray-200"
        >
          Today
        </button>
        <button
          title="Back"
          aria-label="Previous"
          onClick={() => onNavigate("PREV")}
          className="p-2 rounded bg-gray-100 hover:bg-gray-200"
        >
          <ChevronLeft className="w-4 h-4" />
        </button>
        <button
          title="Next"
          aria-label="Next"
          onClick={() => onNavigate("NEXT")}
          className="p-2 rounded bg-gray-100 hover:bg-gray-200"
        >
          <ChevronRight className="w-4 h-4" />
        </button>
        <span className="font-semibold ml-2">{label}</span>
      </div>
      <div className="flex items-center gap-2">
        {viewNames.map((name) => (
          <button
            key={name}
            title={name.charAt(0).toUpperCase() + name.slice(1)}
            onClick={() => onView(name)}
            className={`px-2 py-1 rounded capitalize ${view === name ? "bg-blue-600 text-white" : "bg-gray-100 hover:bg-gray-200"}`}
          >
            {name}
          </button>
        ))}
      </div>
    </div>
  );
}

export default function AdminCalendarPage() {
  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const [selectedItem, setSelectedItem] = useState<Appointment | null>(null);
  const [currentDate, setCurrentDate] = useState<Date>(new Date());
  const [currentView, setCurrentView] = useState<View>("month");

  const { data: appointments = [], isLoading } = useUpcomingAppointments();
  const updateStatus = useUpdateAppointmentStatus();
  const reminder = useSendReminder();
  const reschedule = useRescheduleAppointment();
  const stats = useAppointmentStats();

  // Debug log
  useEffect(() => {
    console.log("📦 Raw appointments from API:", appointments);
    appointments.forEach((appt, index) => {
      console.log(`🔍 [${index}]`, appt);
    });
  }, [appointments]);

  const events: CalendarEvent[] = appointments.flatMap((appt, index) => {
    if (!appt.date || !appt.timeSlot || !appt.timeSlot.includes(":")) {
      console.warn(`⛔ Skipping appointment [${index}] due to invalid timeSlot:`, appt.timeSlot);
      return [];
    }

    const [hourStr, minutePart] = appt.timeSlot.split(":");
    if (!minutePart) return [];

    const hour = parseInt(hourStr, 10);
    let minutes = 0;
    let hour24 = hour;
    const lower = appt.timeSlot.toLowerCase();

    if (minutePart.includes("30")) minutes = 30;
    if (lower.includes("pm") && hour !== 12) hour24 = hour + 12;
    if (lower.includes("am") && hour === 12) hour24 = 0;

    const start = new Date(appt.date);
    start.setHours(hour24, minutes, 0);
    const end = new Date(start.getTime() + 30 * 60 * 1000);

    return [{
      id: appt.id,
      title: `Status: ${appt.status}`,
      start,
      end,
      resource: appt,
    }];
  });

  const filteredAppointments = selectedDate
    ? appointments.filter((a) => a.date === selectedDate)
    : [];

  return (
    <div className="p-6">
      <h1 className="text-2xl font-semibold mb-2">📆 Appointments Calendar</h1>

      {stats.data && (
        <div className="mb-4 flex gap-4 text-sm text-gray-700">
          <span>Total: {stats.data.total}</span>
          <span className="text-green-600">Confirmed: {stats.data.confirmed}</span>
          <span className="text-red-500">Canceled: {stats.data.canceled}</span>
          <span className="text-blue-600">Completed: {stats.data.completed}</span>
        </div>
      )}

      {!isLoading && events.length === 0 && (
        <p className="text-center text-gray-500 mb-2">No events to display.</p>
      )}
      <BigCalendar<CalendarEvent, Appointment>
        localizer={localizer}
        events={events}
        startAccessor={(event) => event.start}
        endAccessor={(event) => event.end}
        style={{ height: "60vh" }}
        date={currentDate}
        onNavigate={(d) => setCurrentDate(d)}
        view={currentView}
        onView={(v) => setCurrentView(v)}
        views={["month", "week", "day", "agenda"]}
        components={{ toolbar: CalendarToolbar }}
        selectable
        onSelectSlot={(slot) => {
          const formatted = format(slot.start, "yyyy-MM-dd");
          setSelectedDate(formatted);
        }}
        onSelectEvent={(event) => {
          setSelectedItem(event.resource);
        }}
        eventPropGetter={(event: CalendarEvent) => {
          let bgColor = "#3182ce";
          if (event.resource.status === "confirmed") bgColor = "#38a169";
          if (event.resource.status === "canceled") bgColor = "#e53e3e";
          if (event.resource.status === "completed") bgColor = "#4c51bf";

          return { style: { backgroundColor: bgColor, color: "white" } };
        }}
      />

      {selectedDate && (
        <DailyAppointmentsModal
          isOpen={true}
          onClose={() => setSelectedDate(null)}
          appointments={filteredAppointments}
          selectedDate={selectedDate}
          onSelectAppointment={(appt) => {
            setSelectedItem(appt);
            setSelectedDate(null);
          }}
          isLoading={isLoading}
        />
      )}

      {selectedItem && selectedItem.id && (
        <UpdateAppointmentModal
          isOpen={true}
          onClose={() => setSelectedItem(null)}
          appointment={selectedItem}
          onStatusChange={(status) => {
            if (!selectedItem?.id) return;
            updateStatus.mutate({ id: selectedItem.id, status });
          }}
          onSendReminder={() => reminder.mutate(selectedItem.id)}
          onReschedule={(newDate, newTimeSlot) =>
            reschedule.mutate({
              id: selectedItem.id,
              date: newDate,
              timeSlot: newTimeSlot,
            })
          }
        />
      )}
    </div>
  );
}
