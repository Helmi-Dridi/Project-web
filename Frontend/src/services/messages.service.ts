import axios from "axios";
import { getAuthHeaders, getCurrentUser } from "./authService";

// ✅ Interfaces
export interface Message {
  id: string;
  sender_id: string;
  receiver_id: string;
  content: string;
  attachment?: string | null;
  company_id: string;
  created_at: string;
  read: boolean;
}

// ✅ Utility
const getCompanyIdOrThrow = (): string => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) throw new Error("No company ID found in user.");
  return user.workCompanyId;
};

const baseURL = "http://localhost:8080/v1/messages";

// ✅ Send a new message
export const sendMessage = async (
  receiverId: string,
  content: string,
  attachment?: string
): Promise<Message> => {
  const companyId = getCompanyIdOrThrow();
  const body = { receiver_id: receiverId, content, attachment };

  const res = await axios.post<{ data: Message }>(
    `${baseURL}/${companyId}/send`,
    body,
    { headers: getAuthHeaders() }
  );
  return res.data.data;
};

export const getMessagesWithUser = async (
  userId: string
): Promise<Message[]> => {
  const companyId = getCompanyIdOrThrow();

  if (!userId) {
    throw new Error("getMessagesWithUser: userId is required");
  }

  const res = await axios.get<{ data: Message[] }>(
    `${baseURL}/${companyId}/user/${userId}`,
    { headers: getAuthHeaders() }
  );
  return res.data.data;
};


// ✅ Mark message as read
export const markMessageAsRead = async (messageId: string): Promise<string> => {
  const companyId = getCompanyIdOrThrow();

  const res = await axios.post<{ data: string }>(
    `${baseURL}/${companyId}/${messageId}/read`,
    {},
    { headers: getAuthHeaders() }
  );
  return res.data.data;
};

// ✅ Get conversation partners
export const getConversationPartners = async (): Promise<string[]> => {
  const companyId = getCompanyIdOrThrow();

  const res = await axios.get<{ data: string[] }>(
    `${baseURL}/${companyId}/partners`,
    { headers: getAuthHeaders() }
  );
  return res.data.data;
};

// ✅ Get unread count
export const getUnreadMessageCount = async (): Promise<number> => {
  const companyId = getCompanyIdOrThrow();

  const res = await axios.get<{ data: { unreadCount: number } }>(
    `${baseURL}/${companyId}/unread-count`,
    { headers: getAuthHeaders() }
  );
  return res.data.data.unreadCount;
};

// ✅ Get paginated messages with a user
export const getPaginatedMessages = async (
  userId: string,
  limit = 20,
  offset = 0
): Promise<Message[]> => {
  const companyId = getCompanyIdOrThrow();

  const res = await axios.get<{ data: Message[] }>(
    `${baseURL}/${companyId}/user/${userId}/paginated?limit=${limit}&offset=${offset}`,
    { headers: getAuthHeaders() }
  );
  return res.data.data;
};

// ✅ Delete a message
export const deleteMessage = async (messageId: string): Promise<string> => {
  const companyId = getCompanyIdOrThrow();

  const res = await axios.delete<{ data: string }>(
    `${baseURL}/${companyId}/${messageId}`,
    { headers: getAuthHeaders() }
  );
  return res.data.data;
};

// ✅ Search messages
export const searchMessages = async (query: string): Promise<Message[]> => {
  const companyId = getCompanyIdOrThrow();

  const res = await axios.get<{ data: Message[] }>(
    `${baseURL}/${companyId}/search?q=${encodeURIComponent(query)}`,
    { headers: getAuthHeaders() }
  );
  return res.data.data;
};

// ✅ Admin inbox
export const getAdminInbox = async (): Promise<Message[]> => {
  const companyId = getCompanyIdOrThrow();

  const res = await axios.get<{ data: Message[] }>(
    `${baseURL}/${companyId}/inbox`,
    { headers: getAuthHeaders() }
  );
  return res.data.data;
};
