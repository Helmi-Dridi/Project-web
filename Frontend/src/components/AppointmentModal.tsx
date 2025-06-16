// src/components/AppointmentModal.tsx
import Modal from "./Modal";
import { useState } from "react";
import type { Appointment } from "../services/booking.service";

type Props = {
  isOpen: boolean;
  onClose: () => void;
  appointment: Appointment;
  onStatusChange?: (status: string) => void;
  onSendReminder?: () => void;
  onReschedule?: (newDate: string, newTimeSlot: string) => void;
};

export default function AppointmentModal({
  isOpen,
  onClose,
  appointment,
  onStatusChange,
  onSendReminder,
  onReschedule,
}: Props) {
  const [newDate, setNewDate] = useState(appointment.date);
  const [newSlot, setNewSlot] = useState("");

  const handleReschedule = () => {
    if (!newSlot || !newDate) return;
    onReschedule?.(newDate, newSlot);
    setNewSlot("");
  };

  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <div className="max-h-[80vh] overflow-y-auto p-2 sm:p-4">
        <h2 className="text-xl font-semibold mb-4">📋 Appointment Details</h2>

        <div className="text-sm space-y-1 mb-4">
          <p><strong>Date:</strong> {appointment.date}</p>
          <p><strong>Time Slot:</strong> {appointment.timeSlot}</p>
          <p><strong>Status:</strong> {appointment.status}</p>
        </div>

        {onStatusChange && (
          <div className="mb-4">
            <label className="block text-sm font-medium mb-1">Change Status</label>
            <select
              defaultValue={appointment.status}
              onChange={(e) => onStatusChange(e.target.value)}
              className="p-2 border rounded w-full bg-white text-black dark:bg-gray-800 dark:text-white"
            >
              <option className="bg-white text-black dark:bg-gray-800 dark:text-white" value="booked">Booked</option>
              <option className="bg-white text-black dark:bg-gray-800 dark:text-white" value="confirmed">Confirmed</option>
              <option className="bg-white text-black dark:bg-gray-800 dark:text-white" value="canceled">Canceled</option>
              <option className="bg-white text-black dark:bg-gray-800 dark:text-white" value="completed">Completed</option>
            </select>
          </div>
        )}

       {onReschedule && (
  <>
    <div className="mb-4">
      <label className="block text-sm font-medium mb-1">New Date</label>
      <input
        type="date"
        className="p-2 border rounded w-full bg-white text-black dark:bg-gray-800 dark:text-white"
        value={newDate}
        onChange={(e) => setNewDate(e.target.value)}
        disabled={appointment.status === "confirmed"} // ❌ Disable if confirmed
      />
    </div>

    <div className="mb-4">
      <label className="block text-sm font-medium mb-1">New Time Slot</label>
      <select
        className="p-2 border rounded w-full bg-white text-black dark:bg-gray-800 dark:text-white"
        value={newSlot}
        onChange={(e) => setNewSlot(e.target.value)}
        disabled={appointment.status === "confirmed"} // ❌ Disable if confirmed
      >
        <option value="">-- Select a time slot --</option>
        <option value="9:00 AM">9:00 AM</option>
        <option value="9:30 AM">9:30 AM</option>
        <option value="10:00 AM">10:00 AM</option>
        <option value="10:30 AM">10:30 AM</option>
        <option value="11:00 AM">11:00 AM</option>
        <option value="11:30 AM">11:30 AM</option>
        <option value="12:00 PM">12:00 PM</option>
        <option value="12:30 PM">12:30 PM</option>
        <option value="1:00 PM">1:00 PM</option>
        <option value="1:30 PM">1:30 PM</option>
        <option value="2:00 PM">2:00 PM</option>
        <option value="2:30 PM">2:30 PM</option>
        <option value="3:00 PM">3:00 PM</option>
        <option value="3:30 PM">3:30 PM</option>
        <option value="4:00 PM">4:00 PM</option>
      </select>
    </div>

    <button
      onClick={handleReschedule}
      className="bg-yellow-500 text-white px-4 py-2 rounded hover:bg-yellow-600 w-full"
      disabled={
        !newSlot ||
        appointment.status === "confirmed" // ❌ Block confirm if confirmed
      }
    >
      🕓 Confirm Reschedule
    </button>
  </>
)}


        {onSendReminder && (
          <button
            onClick={onSendReminder}
            className="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600 w-full mt-4"
          >
            📩 Send Reminder
          </button>
        )}
      </div>
    </Modal>
  );
}
