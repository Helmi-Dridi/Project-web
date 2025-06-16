// src/services/booking.service.ts
import axios from "axios";

// -------------------------
// TYPES
// -------------------------

export type CreateAppointmentPayload = {
  date: string;
  timeSlot: string;
  receiverId: string;
};

export type Appointment = {
  id: string;
  date: string;
  timeSlot: string;
  status: string;
};

export type PaginatedAppointments = {
  items: Appointment[];
};

export type AvailabilityResponse = {
  bookedSlots: string[];
};

export type ApiResponse = {
  status: string;
  message: string;
  data: any;
};

export type AppointmentStats = {
  total: number;
  confirmed: number;
  canceled: number;
  completed: number;
};

// -------------------------
// CONFIG
// -------------------------

const BASE_URL = "http://localhost:8080/v1/appointments";

const getAuthHeaders = () => {
  const token = localStorage.getItem("token");
  return {
    Authorization: `Bearer ${token}`,
  };
};

// -------------------------
// API FUNCTIONS
// -------------------------

export async function createAppointment(data: CreateAppointmentPayload): Promise<ApiResponse> {
  const response = await axios.post(BASE_URL, data, { headers: getAuthHeaders() });
  return response.data;
}

export async function getMyAppointments(): Promise<PaginatedAppointments> {
  const response = await axios.get(`${BASE_URL}?page=1&limit=10`, {
    headers: getAuthHeaders(),
  });
  return  response.data.data;
;
}


export async function getAvailableAppointments(date: string): Promise<AvailabilityResponse> {
  const response = await axios.get(`${BASE_URL}/availability?date=${date}`, {
    headers: getAuthHeaders(),
  });
  return response.data;
}

export async function cancelAppointment(id: string): Promise<ApiResponse> {
  const response = await axios.post(`${BASE_URL}/${id}/cancel`, null, {
    headers: getAuthHeaders(),
  });
  return response.data;
}

export async function getReceivedAppointmentsByDate(date: string): Promise<PaginatedAppointments> {
  const response = await axios.get(`${BASE_URL}/received`, {
    headers: getAuthHeaders(),
    params: { date },
  });
  return response.data.data;
}

export async function updateAppointmentStatus(id: string, status: string): Promise<ApiResponse> {
  const response = await axios.post(
    `${BASE_URL}/${id}/status`,
    { status },
    { headers: getAuthHeaders() }
  );
  return response.data;
}

export async function getAppointmentById(id: string): Promise<Appointment> {
  const response = await axios.get(`${BASE_URL}/${id}`, {
    headers: getAuthHeaders(),
  });
  return response.data.data;
}

export async function getAppointmentStats(): Promise<AppointmentStats> {
  const response = await axios.get(`${BASE_URL}/stats`, {
    headers: getAuthHeaders(),
  });
  return response.data.data;
}

export async function getUpcomingAppointments(): Promise<Appointment[]> {
  const response = await axios.get(`${BASE_URL}/upcoming`, {
    headers: getAuthHeaders(),
  });
  return response.data.data;
}

export async function getAppointmentHistory(): Promise<Appointment[]> {
  const response = await axios.get(`${BASE_URL}/history`, {
    headers: getAuthHeaders(),
  });
  return response.data.data;
}

export async function sendReminder(id: string): Promise<ApiResponse> {
  const response = await axios.post(`${BASE_URL}/${id}/remind`, null, {
    headers: getAuthHeaders(),
  });
  return response.data;
}

export async function rescheduleAppointment(
  id: string,
  date: string,
  timeSlot: string
): Promise<ApiResponse> {
  const response = await axios.post(
    `${BASE_URL}/${id}/reschedule`,
    { date, timeSlot },
    { headers: getAuthHeaders() }
  );
  return response.data;
}
