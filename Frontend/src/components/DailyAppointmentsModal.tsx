import { useState } from "react";
import Modal from "./Modal";
import AppointmentModal from "./AppointmentModal";

type Appointment = {
  id: string;
  date: string;
  timeSlot: string;
  status: string;
};

type Props = {
  isOpen: boolean;
  onClose: () => void;
  appointments: Appointment[];
  selectedDate: string;
  onSelectAppointment?: (appt: Appointment) => void;
  isLoading?: boolean;
};

export default function DailyAppointmentsModal({
  isOpen,
  onClose,
  appointments,
  selectedDate,
  onSelectAppointment,
  isLoading = false,
}: Props) {
  const [activeAppointment, setActiveAppointment] = useState<Appointment | null>(null);

  const handleCloseInner = () => {
    setActiveAppointment(null);
  };

  return (
    <>
      <Modal isOpen={isOpen} onClose={onClose}>
        <h2 className="text-lg font-bold mb-3">Appointments on {selectedDate}</h2>

        {isLoading ? (
          <p>Loading appointments...</p>
        ) : appointments.length === 0 ? (
          <p>No appointments.</p>
        ) : (
          <ul className="space-y-3">
            {appointments.map((appt) => (
              <li
                key={appt.id}
                className="border p-3 rounded shadow-sm cursor-pointer hover:bg-gray-50"
                onClick={() => {
                  setActiveAppointment(appt);
                  if (onSelectAppointment) onSelectAppointment(appt);
                }}
              >
                <p>
                  <strong>Time:</strong> {appt.timeSlot}
                </p>
                <p>
                  <strong>Status:</strong> {appt.status}
                </p>
              </li>
            ))}
          </ul>
        )}
      </Modal>

      {activeAppointment && (
        <AppointmentModal
          isOpen={true}
          onClose={handleCloseInner}
          appointment={activeAppointment}
          // Optional callbacks can be added here, e.g.:
          // onStatusChange={(status) => ...}
          // onSendReminder={() => ...}
          // onReschedule={(newDate, newSlot) => ...}
        />
      )}
    </>
  );
}
