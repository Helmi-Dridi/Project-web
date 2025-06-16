// src/services/admin.service.ts
import axios from "axios";
import { getAuthHeaders, getCurrentUser } from "./authService";

// ✅ Updated with full admin details structure
export interface AdminDetails {
  id: string; // ✅ required for role assignment
  firstName: string;
  lastName: string;
  email: string;
  country: string;
  status: boolean;
}
export interface CreateAdminInput {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
  country?: string;
  companyID:string;
}

export interface CreateAdminResponse {
  success: boolean;
  message: string;
}

export const fetchAllAdminDetails = async (): Promise<AdminDetails[]> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) {
    throw new Error("No valid user or workCompanyId found");
  }

  const response = await axios.get<{ data: AdminDetails[] }>(
    `http://localhost:8080/api/users/${user.workCompanyId}/admins/list/all`,
    { headers: getAuthHeaders() }
  );

  return response.data.data;
};

// ✅ Existing simplified admin list for dropdowns etc.
export interface AdminUser {
  id: string;
  name: string;
}

export const fetchAdmins = async (): Promise<AdminUser[]> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) {
    throw new Error("No valid user or workCompanyId found");
  }

  const response = await axios.get<{ data: AdminUser[] }>(
    `http://localhost:8080/api/users/${user.workCompanyId}/admins/list`,
    { headers: getAuthHeaders() }
  );

  return response.data.data;
};

// ✅ Upload user profile image
export const uploadProfileImage = async (file: File): Promise<string> => {
  const user = getCurrentUser();
  if (!user?.ID || !user?.workCompanyId) throw new Error("Invalid user session");

  const formData = new FormData();
  formData.append("file", file);

  const response = await axios.post(
    `http://localhost:8080/api/users/${user.workCompanyId}/${user.ID}/upload-profile-image`,
    formData,
    {
      headers: {
        ...getAuthHeaders(),
        "Content-Type": "multipart/form-data",
      },
    }
  );

  return response.data?.data?.profile_picture;
};
export const createAdmin = async (adminData: CreateAdminInput): Promise<CreateAdminResponse> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) throw new Error("No company ID");
  adminData.companyID = user.workCompanyId
  const response = await axios.post(
    `http://localhost:8080/api/users/${user.workCompanyId}`,
    adminData,
    {
      headers: getAuthHeaders(),
    }
  );

  return response.data;
};
export const assignAdminRole = async (userId: string, roleId: string): Promise<void> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) throw new Error("No company ID");

  await axios.post(
    `http://localhost:8080/api/users/${user.workCompanyId}/${userId}/roles/${roleId}`,
    {},
    { headers: getAuthHeaders() }
  );
};
export const deleteAdmin = async (userId: string): Promise<void> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) throw new Error("No company ID");

  await axios.delete(
    `http://localhost:8080/api/users/${user.workCompanyId}/${userId}`,
    { headers: getAuthHeaders() }
  );
};
export interface BasicUser {
  id: string; // ✅ required for role assignment
  firstName: string;
  lastName: string;
  email: string;
  country: string;
  status: boolean;
}
export const fetchUsersWithoutRole = async (): Promise<BasicUser[]> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) throw new Error("No company ID");

  const response = await axios.get<{ data: BasicUser[] }>(
    `http://localhost:8080/api/users/${user.workCompanyId}/norole`,
    { headers: getAuthHeaders() }
  );

  return response.data.data;
};

export const fetchStudents = async (): Promise<BasicUser[]> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) throw new Error("No company ID");

  const response = await axios.get<{ data: BasicUser[] }>(
    `http://localhost:8080/api/users/${user.workCompanyId}/students`,
    { headers: getAuthHeaders() }
  );

  return response.data.data;
};
