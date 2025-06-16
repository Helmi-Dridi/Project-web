// src/services/notification.service.ts

import axios from "axios";
import { getCurrentUser, getAuthHeaders } from "./authService";

// Interfaces
export interface Notification {
  id: string;
  type: string;
  content: string;
  seen: boolean;
  createdAt: string;
}

export interface NotificationInput {
  seen: boolean;
}

export interface CreateNotificationInput {
  type: string;
  content: string;
  seen?: boolean; // Optional – defaults to false in backend
}

export interface NotificationCount {
  count: number;
}

// 🔍 Get all notifications for the current user
export const fetchNotifications = async (): Promise<Notification[]> => {
  const user = getCurrentUser();
  if (!user?.ID) throw new Error("Invalid user session");

  const response = await axios.get<{ data: Notification[] }>(
    `http://localhost:8080/api/notifications/${user.ID}`,
    { headers: getAuthHeaders() }
  );

  return response.data.data;
};

// 🔢 Get count of unread notifications
export const fetchNotificationCount = async (): Promise<NotificationCount> => {
  const user = getCurrentUser();
  if (!user?.ID) throw new Error("Invalid user session");

  const response = await axios.get<{ data: NotificationCount }>(
    `http://localhost:8080/api/notifications/${user.ID}/count`,
    { headers: getAuthHeaders() }
  );

  return response.data.data;
};

// 📄 Get a specific notification by ID
export const fetchNotificationById = async (
  notificationId: string
): Promise<Notification> => {
  const user = getCurrentUser();
  if (!user?.ID) throw new Error("Invalid user session");

  const response = await axios.get<{ data: Notification }>(
    `http://localhost:8080/api/notifications/${user.ID}/${notificationId}`,
    { headers: getAuthHeaders() }
  );

  return response.data.data;
};



// ✅ Update (mark as seen)
export const updateNotification = async (
  notificationId: string,
  payload: NotificationInput
): Promise<void> => {
  const user = getCurrentUser();
  if (!user?.ID) throw new Error("Invalid user session");

  await axios.put(
    `http://localhost:8080/api/notifications/${user.ID}/${notificationId}`,
    payload,
    { headers: getAuthHeaders() }
  );
};

// ❌ Delete a notification
export const deleteNotification = async (
  notificationId: string
): Promise<void> => {
  const user = getCurrentUser();
  if (!user?.ID) throw new Error("Invalid user session");

  await axios.delete(
    `http://localhost:8080/api/notifications/${user.ID}/${notificationId}`,
    { headers: getAuthHeaders() }
  );
};
