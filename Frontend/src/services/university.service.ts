// src/services/university.service.ts
import axios from "axios";
import { getAuthHeaders, getCurrentUser } from "./authService";

// ✅ Interface for basic university list
export interface UniversityTable {
  id: string;
  name: string;
  country: string;
  city: string;
  universityType: string;
  createdAt: string;
}

// ✅ Interface for university programs
export interface UniversityProgramInfo {
  id: string;
  programType: string;
  programName: string;
  languageOfInstruction: string;
  minimumGpaRequirement: string;
  languageTestRequirement: string;
  otherTestRequirements: string;
  tuitionFeesPerYear: string;
  estimatedLivingExpenses: string;
  visaRequirements: string;
  proofOfFinancialMeans: string;
  blockedAccountRequired: string;
  applicationDeadline: string;
  applicationFee: string;
  scholarshipsAvailable: string;
  intakePeriods: string;
  partTimeWorkAllowed: string;
  recognitionInTunisia: string;
  notes: string;
  comments: string;
}
export interface UniversityPagination {
  items: UniversityTable[]; // the list of universities
  page: number;             // current page
  limit: number;            // items per page
  totalCount: number;       // total universities count
}
export const fetchAllUniversities = async (page: number, limit: number): Promise<UniversityPagination> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) {
    throw new Error("No valid user or workCompanyId found");
  }

  const response = await axios.get<{ data: UniversityPagination }>(
    `http://localhost:8080/v1/universities/${user.workCompanyId}?page=${page}&limit=${limit}`,
    { headers: getAuthHeaders() }
  );

  return response.data.data;
};

// ✅ Fetch all programs for a specific university
export const fetchUniversityPrograms = async (universityId: string): Promise<UniversityProgramInfo[]> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) {
    throw new Error("No valid user or workCompanyId found");
  }

  const response = await axios.get<{ data: UniversityProgramInfo[] }>(
    `http://localhost:8080/v1/universities/${user.workCompanyId}/${universityId}/programs`,
    { headers: getAuthHeaders() }
  );

  return response.data.data;
};
