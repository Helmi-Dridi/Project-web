import type { Appointment } from "../services/booking.service";
import AppointmentModal from "./AppointmentModal";

type Props = {
  isOpen: boolean;
  onClose: () => void;
  appointment: Appointment;
  onStatusChange: (status: string) => void;
  onSendReminder?: () => void;
  onReschedule?: (newDate: string, newTimeSlot: string) => void;
};

export default function UpdateAppointmentModal({
  isOpen,
  onClose,
  appointment,
  onStatusChange,
  onSendReminder,
  onReschedule,
}: Props) {
  return (
    <AppointmentModal
      isOpen={isOpen}
      onClose={onClose}
      appointment={appointment}
      onStatusChange={onStatusChange}
      onSendReminder={onSendReminder}
      onReschedule={onReschedule}
    />
  );
}
