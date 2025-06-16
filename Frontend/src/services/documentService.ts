// src/services/documentService.ts

import axios from "axios";

// Types inline (replace with external types if needed)
export type StudentDocumentIn = {
  documentType: string;
  filePath: string;
  isFinalized?: boolean;
};

export type StudentDocumentTable = {
  id: string;
  documentType: string;
  filePath: string;
  isFinalized: boolean;
  verification: string;
  remarks: string;
  uploadedAt: string;
  version: number;
};

export type StudentDocumentDetails = StudentDocumentTable & {
  userID: string;
};

export type ApiResponse = {
  status: string;
  message: string;
  data: any;
};

const BASE_URL = "http://localhost:8080/v1/documents";

// Helper: get bearer token from localStorage
const getAuthHeaders = () => {
  const token = localStorage.getItem("token");
  return {
    Authorization: `Bearer ${token}`,
  };
};

// Upload document with file (multipart)
export async function uploadDocument(formData: FormData): Promise<ApiResponse> {
  const response = await axios.post(`${BASE_URL}/upload`, formData, {
    headers: {
      ...getAuthHeaders(),
      "Content-Type": "multipart/form-data",
    },
  });
  return response.data;
}

// Create metadata-based document (JSON)
export async function createDocumentMetadata(data: StudentDocumentIn): Promise<ApiResponse> {
  const response = await axios.post(BASE_URL, data, {
    headers: getAuthHeaders(),
  });
  return response.data;
}

// Get paginated documents
export async function getDocuments(page = 1, limit = 10): Promise<{ items: StudentDocumentTable[]; totalCount: number }> {
  const response = await axios.get(`${BASE_URL}?page=${page}&limit=${limit}`, {
    headers: getAuthHeaders(),
  });
  return response.data.data;
}

// Get list of all documents (no pagination)
export async function getAllDocuments(): Promise<StudentDocumentTable[]> {
  const response = await axios.get(`${BASE_URL}/list`, {
    headers: getAuthHeaders(),
  });
  return response.data.data;
}

// Get one document by ID
export async function getDocumentById(id: string): Promise<StudentDocumentDetails> {
  const response = await axios.get(`${BASE_URL}/${id}`, {
    headers: getAuthHeaders(),
  });
  return response.data.data;
}

// Delete document by ID
export async function deleteDocument(id: string): Promise<ApiResponse> {
  const response = await axios.delete(`${BASE_URL}/${id}`, {
    headers: getAuthHeaders(),
  });
  return response.data;
}

// Finalize document
export async function finalizeDocument(id: string): Promise<ApiResponse> {
  const response = await axios.post(`${BASE_URL}/${id}/finalize`, null, {
    headers: getAuthHeaders(),
  });
  return response.data;
}

// Download document
export async function downloadDocument(id: string): Promise<Blob> {
  const response = await axios.get(`${BASE_URL}/${id}/download`, {
    headers: getAuthHeaders(),
    responseType: "blob",
  });
  return response.data;
}

export async function updateDocument(input: {
  id: string;
  payload: Partial<StudentDocumentIn>;
}): Promise<ApiResponse> {
  const { id, payload } = input;

  const response = await axios.put(`${BASE_URL}/${id}`, payload, {
    headers: getAuthHeaders(),
  });
  return response.data;
}
// Get completion status (for progress bar)
export async function getDocumentStatus(): Promise<{ completion: number; required: number; uploaded: number }> {
  const response = await axios.get(`${BASE_URL}/status`, {
    headers: getAuthHeaders(),
  });
  return response.data.data;
}


// Get all documents for a specific student (admin only)
export async function getDocumentsByUserId(userId: string): Promise<StudentDocumentTable[]> {
  const response = await axios.get(`http://localhost:8080/v1/documents/user/${userId}`, {
    headers: getAuthHeaders(),
  });
  return response.data.data;
}
export async function verifyDocument(id: string, payload: { verification: string; remarks?: string }): Promise<ApiResponse> {
  const response = await axios.post(`http://localhost:8080/v1/documents/${id}/verify`, payload, {
    headers: getAuthHeaders(),
  });
  return response.data;
}