import axios from "axios";
import camelcaseKeys from "camelcase-keys";

// ENUM TYPES
export type Gender = "male" | "female" | "other";
export type PreferredContactMethod = "whatsapp" | "email" | "call";

// INTERFACES
export interface StudentProfile {
  id?: string;
  userId?: string;
  dateOfBirth?: string;
  gender: Gender;
  nationality: string;
  passportNumber?: string;
  nationalId: string;
step: boolean;

}

export interface StudentContactDetails {
  id?: string;
  userId?: string;
  phoneNumber: string;
  email: string;
  addressLine?: string;
  city?: string;
  country?: string;
  zipCode?: string;
  preferredContactMethod: PreferredContactMethod;
step: boolean;

}

export interface StudentAcademicBackground {
  id?: string;
  userId?: string;
  qualification: string;
  institutionName: string;
  graduationYear: number;
  gpaScore?: string;
  languageTestType?: string;
  languageTestScore?: string;
  certificateFilePath?: string;
step: boolean;

}

export interface StudentSettings {
  id?: string;
  userId?: string;
  languagePreference: string;
  receiveEmailNotifications: boolean;
  receiveSMSNotifications: boolean;
  accountDeleted?: boolean;
  emergencyName: string;
  emergencyRelationship: string;
  emergencyPhoneNumber: string;
  emergencyEmail: string;
step: boolean;
}

export type ApiResponse<T = any> = {
  status: string;
  message: string;
  data: T;
};

// CONFIG
const BASE = "http://localhost:8080/v1";
const authHeaders = () => ({
  Authorization: `Bearer ${localStorage.getItem("token") ?? ""}`,
});

// PROFILE API
export const getStudentProfile = async (): Promise<ApiResponse<StudentProfile>> => {
  const res = await axios.get(`${BASE}/student-profile`, { headers: authHeaders() });
  return {
    ...res.data,
    data: camelcaseKeys(res.data.data, { deep: true }),
  };
};

export const createOrUpdateStudentProfile = async (data: StudentProfile): Promise<ApiResponse> => {
  const res = await axios.post(`${BASE}/student-profile`, data, { headers: authHeaders() });
  return res.data;
};

// CONTACT API
export const getStudentContact = async (): Promise<ApiResponse<StudentContactDetails>> => {
  const res = await axios.get(`${BASE}/student-contact`, { headers: authHeaders() });
  return {
    ...res.data,
    data: camelcaseKeys(res.data.data, { deep: true }),
  };
};

export const createOrUpdateStudentContact = async (data: StudentContactDetails): Promise<ApiResponse> => {
  const res = await axios.post(`${BASE}/student-contact`, data, { headers: authHeaders() });
  return res.data;
};

// ACADEMIC API
export const getStudentAcademic = async (): Promise<ApiResponse<StudentAcademicBackground>> => {
  const res = await axios.get(`${BASE}/student-academic`, { headers: authHeaders() });
  return {
    ...res.data,
    data: camelcaseKeys(res.data.data, { deep: true }),
  };
};

export const createOrUpdateStudentAcademic = async (data: FormData): Promise<ApiResponse> => {
  const res = await axios.post(`${BASE}/student-academic`, data, {
    headers: {
      ...authHeaders(),
      "Content-Type": "multipart/form-data",
    },
  });
  return res.data;
};


// SETTINGS API
export const getStudentSettings = async (): Promise<ApiResponse<StudentSettings>> => {
  const res = await axios.get(`${BASE}/student/settings`, { headers: authHeaders() });
  return {
    ...res.data,
    data: camelcaseKeys(res.data.data, { deep: true }),
  };
};

export const createOrUpdateStudentSettings = async (data: StudentSettings): Promise<ApiResponse> => {
  const res = await axios.post(`${BASE}/student/settings`, data, { headers: authHeaders() });
  return res.data;
};
export const getStudentSettingsByUserId = async (
  userId: string
): Promise<ApiResponse<StudentSettings>> => {
  const res = await axios.get(`${BASE}/student/settings/${userId}`, {
    headers: authHeaders(),
  });

  return {
    ...res.data,
    data: camelcaseKeys(res.data.data, { deep: true }),
  };
};
export const getStudentProfileByUserId = async (
  userId: string
): Promise<ApiResponse<StudentProfile>> => {
  const res = await axios.get(`${BASE}/student-profile/${userId}`, {
    headers: authHeaders(),
  });

  return {
    ...res.data,
    data: camelcaseKeys(res.data.data, { deep: true }),
  };
};

export const getStudentContactByUserId = async (
  userId: string
): Promise<ApiResponse<StudentContactDetails>> => {
  const res = await axios.get(`${BASE}/student-contact/${userId}`, {
    headers: authHeaders(),
  });

  return {
    ...res.data,
    data: camelcaseKeys(res.data.data, { deep: true }),
  };
};
export const getStudentAcademicByUserId = async (
  userId: string
): Promise<ApiResponse<StudentAcademicBackground>> => {
  const res = await axios.get(`${BASE}/student-academic/${userId}`, {
    headers: authHeaders(),
  });

  return {
    ...res.data,
    data: camelcaseKeys(res.data.data, { deep: true }),
  };
};
