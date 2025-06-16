import axios from "axios";
import { getAuthHeaders, getCurrentUser } from "./authService";

const API_URL = "http://localhost:8080/api";

export interface Role {
  id: string;
  name: string;
  description?: string;
  companyId: string;
}

export const fetchRoleall = async (): Promise<Role[]> => {
  const user = getCurrentUser();
  if (!user?.workCompanyId) throw new Error("Missing company ID");

  const response = await axios.get<{ data: Role[] }>(
    `${API_URL}/roles/${user.workCompanyId}/list`, // ✅ Corrected endpoint
    { headers: getAuthHeaders() }
  );

  const roles = response.data.data;
  if (!Array.isArray(roles)) {
    console.error("Expected roles to be an array, got:", roles);
    return [];
  }

  return roles;
};

