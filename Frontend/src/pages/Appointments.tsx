import { useState, useEffect } from "react";
import Calendar from "react-calendar";
import "react-calendar/dist/Calendar.css";
import { format } from "date-fns";
import { useCreateAppointment, useAvailableAppointments } from "../hooks/useBooking";
import { getCurrentUser } from "../services/authService";
import { fetchAdmins, type AdminUser } from "../services/admin.service";
import toast, { Toaster } from "react-hot-toast";

const ALL_TIME_SLOTS = ["10:00 AM", "11:00 AM", "1:00 PM", "2:30 PM", "4:00 PM"];

export default function AppointmentPage() {
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);
  const [selectedTime, setSelectedTime] = useState<string | null>(null);
  const [admins, setAdmins] = useState<AdminUser[]>([]);
  const [selectedAdminId, setSelectedAdminId] = useState<string | null>(null);

  const user = getCurrentUser();
  const isoDate = selectedDate?.toISOString().split("T")[0] || "";
  const { data, isLoading: isLoadingSlots } = useAvailableAppointments(isoDate);
  const { mutate: bookAppointment, isPending } = useCreateAppointment();

  useEffect(() => {
    if (!user?.workCompanyId) return;

    const loadAdmins = async () => {
      try {
        const result = await fetchAdmins();
        setAdmins(result.map((admin) => ({ id: admin.id, name: admin.name })));
      } catch (err) {
        console.error("❌ Could not fetch admins:", err);
      }
    };

    loadAdmins();
  }, [user?.workCompanyId]);

  const handleDateChange = (value: Date | Date[] | null) => {
    if (!value || Array.isArray(value)) return;
    setSelectedDate(value);
    setSelectedTime(null);
  };

  const confirmBooking = () => {
    if (selectedDate && selectedTime && selectedAdminId) {
      bookAppointment(
        {
          date: isoDate,
          timeSlot: selectedTime,
          receiverId: selectedAdminId,
        },
        {
          onSuccess: () => toast.success("🎉 Appointment booked successfully!"),
          onError: () => toast.error("❌ Booking failed. Try again."),
        }
      );
    }
  };

  const bookedSlots = data?.bookedSlots || [];

  return (
    <div className="min-h-screen bg-gray-50 text-gray-900 py-12 px-4">
      <Toaster />
      <div className="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        {/* Profile */}
        <div className="bg-white rounded-2xl shadow-md p-6 flex flex-col items-center text-center border border-gray-200">
          <img
            src={user?.profilePicture || "https://i.pravatar.cc/100"}
            alt="User"
            className="w-24 h-24 rounded-full border-4 border-blue-600 shadow"
          />
          <h2 className="mt-4 text-xl font-semibold">{user?.name || "Student Name"}</h2>
          <p className="text-gray-500 text-sm">Student Portal</p>

          <div className="mt-6 w-full text-sm text-left space-y-3">
            <div className="flex items-center gap-2">
              <span>🕒</span> <span>30-minute consultation</span>
            </div>
            <div className="flex items-center gap-2">
              <span>📹</span> <span>Online via Zoom</span>
            </div>
            <div className="flex items-center gap-2">
              <span>🌍</span> <span>Time zone: Eastern (US & Canada)</span>
            </div>
          </div>
        </div>

        {/* Calendar */}
        <div className="bg-white rounded-2xl shadow-md p-6 border border-gray-200">
          <h3 className="text-lg font-bold mb-4">📅 Choose a Date</h3>
          <Calendar
            onChange={(value) => handleDateChange(value as Date)}
            value={selectedDate}
            minDate={new Date()}
            tileClassName={({ date }) =>
              selectedDate?.toDateString() === date.toDateString()
                ? "bg-blue-600 text-white font-semibold rounded-full"
                : "hover:bg-blue-100 rounded-full transition text-gray-800"
            }
          />
          {selectedDate && (
            <p className="mt-4 text-sm text-center text-gray-600">
              You selected: <strong>{format(selectedDate, "PPP")}</strong>
            </p>
          )}
        </div>

        {/* Time & Admin */}
        <div className="bg-white rounded-2xl shadow-md p-6 border border-gray-200">
          <h3 className="text-lg font-bold mb-4">🧑‍💼 Choose Admin</h3>
          <select
            value={selectedAdminId || ""}
            onChange={(e) => setSelectedAdminId(e.target.value || null)}
            className="w-full mb-6 px-4 py-2 border border-gray-300 rounded-lg bg-white text-gray-800 focus:ring-2 focus:ring-blue-500"
          >
            <option value="" disabled>
              -- Select Admin --
            </option>
            {admins.map((admin) => (
              <option key={admin.id} value={admin.id}>
                {admin.name}
              </option>
            ))}
          </select>

          <h3 className="text-lg font-bold mb-4">⏰ Select Time Slot</h3>
          {selectedDate ? (
            <div className="space-y-3">
              {isLoadingSlots ? (
                <p className="text-sm text-gray-500">Checking availability...</p>
              ) : (
                ALL_TIME_SLOTS.map((slot) => {
                  const isBooked = bookedSlots.includes(slot);
                  const isSelected = selectedTime === slot;

                  return (
                    <div key={slot} className="flex items-center gap-3">
                      <button
                        onClick={() => setSelectedTime(slot)}
                        disabled={isBooked}
                        className={`flex-1 px-4 py-2 rounded-md text-sm font-medium shadow-sm transition ${
                          isBooked
                            ? "bg-gray-200 text-gray-500 cursor-not-allowed"
                            : isSelected
                            ? "bg-blue-600 text-white"
                            : "bg-gray-100 text-gray-800 hover:bg-blue-100"
                        }`}
                      >
                        {slot}
                      </button>

                      {isSelected && !isBooked && (
                        <button
                          onClick={confirmBooking}
                          disabled={isPending || !selectedAdminId}
                          className="px-4 py-2 bg-blue-600 text-white rounded-md text-sm font-semibold hover:bg-blue-700 disabled:opacity-50"
                        >
                          {isPending ? "Booking..." : "Confirm"}
                        </button>
                      )}
                    </div>
                  );
                })
              )}
            </div>
          ) : (
            <p className="text-sm text-gray-500 mt-2">Select a date first to see available slots.</p>
          )}
        </div>
      </div>

      {/* Info Block */}
      <div className="mt-12 max-w-3xl mx-auto text-center text-sm text-gray-500">
        📘 This page allows you to schedule a meeting with a university advisor. Select a date, pick an available time, and confirm to receive your Zoom link. If no slots are available, try another day.
      </div>
    </div>
  );
}
