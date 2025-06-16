// src/services/enrollment.service.ts
import axios from "axios";
import { getAuthHeaders, getCurrentUser } from "./authService";

// ✅ Interfaces
export interface ProgramDetails {
  id: string;
  programType: string;
  programName: string;
  languageOfInstruction: string;
}
export interface StudentInfo {
  firstName: string;
  lastName: string;
  email: string;
}

export interface EnrollmentResponse {
  userId: string;
  programId: string;
  enrolledAt: string;
  program: ProgramDetails;
  university: UniversityDetails;
  student: StudentInfo; // ✅ new field
}

export interface UniversityDetails {
  id: string;
  name: string;
  country: string;
  city: string;
  universityType: string;
}


// ✅ Assign student to a program
export const assignStudentToProgram = async (
  userId: string,
  programId: string
): Promise<EnrollmentResponse> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) {
    throw new Error("No valid user or workCompanyId found");
  }

  const response = await axios.post<{ data: EnrollmentResponse }>(
    `http://localhost:8080/v1/enrollments/${user.workCompanyId}/users/${userId}/programs/${programId}`,
    {},
    { headers: getAuthHeaders() }
  );

  return response.data.data;
};

// ✅ Unassign student from a program
export const unassignStudentFromProgram = async (
  userId: string
): Promise<string> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) {
    throw new Error("No valid user or workCompanyId found");
  }

  const response = await axios.delete<{ data: string }>(
    `http://localhost:8080/v1/enrollments/${user.workCompanyId}/users/${userId}`,
    { headers: getAuthHeaders() }
  );

  return response.data.data;
};

// ✅ Get student’s own enrollment (for student dashboard)
export const getStudentEnrollment = async (): Promise<EnrollmentResponse> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) {
    throw new Error("No valid user or workCompanyId found");
  }

  const response = await axios.get<{ data: EnrollmentResponse }>(
    `http://localhost:8080/v1/enrollments/${user.workCompanyId}/view`,
    { headers: getAuthHeaders() }
  );

  return response.data.data;
};

// ✅ NEW: Get all enrollments for a specific program
export const getEnrollmentsByProgram = async (
  programId: string
): Promise<EnrollmentResponse[]> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) {
    throw new Error("No valid user or workCompanyId found");
  }

  const response = await axios.get<{ data: EnrollmentResponse[] }>(
    `http://localhost:8080/v1/enrollments/${user.workCompanyId}/programs/${programId}`,
    { headers: getAuthHeaders() }
  );

  return response.data.data;
};

// ✅ NEW: Get all student enrollments across all programs
export const getAllStudentEnrollments = async (): Promise<EnrollmentResponse[]> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) {
    throw new Error("No valid user or workCompanyId found");
  }

  const response = await axios.get<{ data: EnrollmentResponse[] }>(
    `http://localhost:8080/v1/enrollments/${user.workCompanyId}/view/all`,
    { headers: getAuthHeaders() }
  );

  return response.data.data;
};
