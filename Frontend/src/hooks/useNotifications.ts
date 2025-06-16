import {
  fetchNotifications,
  fetchNotificationById,
  fetchNotificationCount,
  updateNotification,
  deleteNotification,
  type Notification,
  type NotificationInput,
  type NotificationCount,
} from "../services/notification.service";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";

// 🔁 Prevents refetching when switching tabs
const queryOptions = {
  refetchOnWindowFocus: false,
};

export const useNotifications = () =>
  useQuery<Notification[]>({
    queryKey: ["notifications"],
    queryFn: fetchNotifications,
    refetchInterval: 3000, // ⏱️ refetch every 5 seconds
    ...queryOptions,
  });


export const useNotificationCount = () =>
  useQuery<NotificationCount>({
    queryKey: ["notifications", "count"],
    queryFn: fetchNotificationCount,
    ...queryOptions,
  });

export const useNotificationById = (notificationId: string) =>
  useQuery<Notification>({
    queryKey: ["notifications", notificationId],
    queryFn: () => fetchNotificationById(notificationId),
    enabled: !!notificationId,
    ...queryOptions,
  });



export const useUpdateNotification = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: string;
      payload: NotificationInput;
    }) => updateNotification(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      queryClient.invalidateQueries({ queryKey: ["notifications", "count"] });
    },
  });
};

export const useDeleteNotification = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => deleteNotification(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] });
      queryClient.invalidateQueries({ queryKey: ["notifications", "count"] });
    },
  });
};
