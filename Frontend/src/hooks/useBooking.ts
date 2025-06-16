import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "react-hot-toast";
import {
  createAppointment,
  getAvailableAppointments,
  getMyAppointments,
  cancelAppointment,
  updateAppointmentStatus,
  getReceivedAppointmentsByDate,
  getAppointmentById,
  getAppointmentStats,
  getUpcomingAppointments,
  getAppointmentHistory,
  sendReminder,
  rescheduleAppointment,
} from "../services/booking.service";

import { getCurrentUser, getAuthHeaders } from "../services/authService";
import axios from "axios";

import type { CreateAppointmentPayload } from "../services/booking.service";
import { useAuth } from "../context/AuthContext";

const NO_REFETCH = {
  refetchOnWindowFocus: false,
  refetchOnReconnect: false,
};

// 🔔 Internal helper for sending a notification
const createNotification = async ({
  userID,
  type,
  content,
}: {
  userID: string;
  type: string;
  content: string;
}) => {
  await axios.post(
    `http://localhost:8080/api/notifications/${userID}`,
    {
      type,
      content,
      seen: false,
    },
    { headers: getAuthHeaders() }
  );
};

export const useCreateAppointment = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateAppointmentPayload) => createAppointment(data),
    onSuccess: () => {
      toast.success("✅ Booking confirmed!");
      queryClient.invalidateQueries({ queryKey: ["appointments"] });
      queryClient.invalidateQueries({ queryKey: ["availableAppointments"] });
    },
    onError: (error: any) => {
      const status = error?.response?.status;
      if (status === 409) {
        toast.error("❌ This time slot is already booked.");
      } else {
        toast.error("⚠️ Booking failed.");
      }
    },
  });
};

export const useAvailableAppointments = (date: string) =>
  useQuery({
    queryKey: ["availableAppointments", date],
    queryFn: () => getAvailableAppointments(date),
    enabled: !!date,
    staleTime: 1000 * 60 * 10,
    ...NO_REFETCH,
  });
export const useMyAppointments = () => {
  const { user } = useAuth();

  return useQuery({
    queryKey: ["appointments", user?.ID], // include ID in key for caching per user
    queryFn: async () => {
      console.log("📅 Fetching my appointments for user:", user?.ID);
      return await getMyAppointments();
    },
    enabled: !!user?.ID, // only run when user is loaded
    staleTime: 1000 * 60 * 5,
    refetchOnWindowFocus: false,
    refetchOnReconnect: false,
  });
};

export const useCancelAppointment = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: cancelAppointment,
    onSuccess: () => {
      toast.success("✅ Appointment canceled");
      queryClient.invalidateQueries({ queryKey: ["appointments"] });
    },
    onError: () => {
      toast.error("❌ Failed to cancel appointment");
    },
  });
};

export const useReceivedAppointmentsByDate = (date: string) =>
  useQuery({
    queryKey: ["received-appointments", date],
    queryFn: () => getReceivedAppointmentsByDate(date),
    enabled: !!date,
    ...NO_REFETCH,
  });

export const useUpdateAppointmentStatus = () => {
  const queryClient = useQueryClient();
  const user = getCurrentUser();

  return useMutation({
    mutationFn: ({ id, status }: { id: string; status: string }) =>
      updateAppointmentStatus(id, status),
    onSuccess: async (_, { status }) => {
      toast.success("✅ Status updated");
      queryClient.invalidateQueries({ queryKey: ["upcomingAppointments"] });

      if (user?.ID) {
        await createNotification({
          userID: user.ID,
          type: "Appointment",
          content: `You changed the status of appointment  to ${status}`,
        });
      }
    },
    onError: () => {
      toast.error("❌ Failed to update status");
    },
  });
};

export const useAppointmentById = (id: string) =>
  useQuery({
    queryKey: ["appointment", id],
    queryFn: () => getAppointmentById(id),
    enabled: !!id,
    ...NO_REFETCH,
  });

export const useAppointmentStats = () =>
  useQuery({
    queryKey: ["appointmentStats"],
    queryFn: getAppointmentStats,
    staleTime: 1000 * 60 * 10,
    ...NO_REFETCH,
  });

export const useUpcomingAppointments = () =>
  useQuery({
    queryKey: ["upcomingAppointments"],
    queryFn: getUpcomingAppointments,
    ...NO_REFETCH,
  });

export const useAppointmentHistory = () =>
  useQuery({
    queryKey: ["appointmentHistory"],
    queryFn: getAppointmentHistory,
    ...NO_REFETCH,
  });

export const useSendReminder = () => {
  const user = getCurrentUser();

  return useMutation({
    mutationFn: sendReminder,
    onSuccess: async (_, ) => {
      toast.success("✅ Reminder sent!");
      if (user?.ID) {
        await createNotification({
          userID: user.ID,
          type: "Reminder",
          content: `You sent a reminder for appointment`,
        });
      }
    },
    onError: () => {
      toast.error("❌ Reminder failed");
    },
  });
};

export const useRescheduleAppointment = () => {
  const queryClient = useQueryClient();
  const user = getCurrentUser();

  return useMutation({
    mutationFn: (data: { id: string; date: string; timeSlot: string }) =>
      rescheduleAppointment(data.id, data.date, data.timeSlot),
    onSuccess: async (_, {  date, timeSlot }) => {
      toast.success("✅ Appointment rescheduled");
      queryClient.invalidateQueries({ queryKey: ["upcomingAppointments"] });

      if (user?.ID) {
        await createNotification({
          userID: user.ID,
          type: "Reschedule",
          content: `You rescheduled appointment  to ${date} at ${timeSlot}`,
        });
      }
    },
    onError: () => {
      toast.error("❌ Reschedule failed");
    },
  });
};
