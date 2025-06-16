import { useState, useEffect } from "react";
import {
  Calendar as BigCalendar,
  dateFnsLocalizer,
} from "react-big-calendar";
import { format, parse, startOfWeek, getDay } from "date-fns";
import { enUS } from "date-fns/locale";
import "react-big-calendar/lib/css/react-big-calendar.css";

import DailyAppointmentsModal from "../components/DailyAppointmentsModal";
import UpdateAppointmentModal from "../components/UpdateAppointmentModal";

import {
  useUpcomingAppointments,
  useUpdateAppointmentStatus,
  useAppointmentStats,
  useSendReminder,
  useRescheduleAppointment,
} from "../hooks/useBooking";

// Localization
const locales = { "en-US": enUS };
const localizer = dateFnsLocalizer({
  format,
  parse,
  startOfWeek,
  getDay,
  locales,
});

type Appointment = {
  id: string;
  date: string;
  timeSlot: string;
  status: string;
};

type CalendarEvent = {
  id: string;
  title: string;
  start: Date;
  end: Date;
  resource: Appointment;
};

export default function AdminCalendarPage() {
  const [selectedDate, setSelectedDate] = useState<string | null>(null);
  const [selectedItem, setSelectedItem] = useState<Appointment | null>(null);

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

      <BigCalendar
        localizer={localizer}
        events={events}
        startAccessor={(event) => event.start}
        endAccessor={(event) => event.end}
        style={{ height: 600 }}
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
