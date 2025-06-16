// src/hooks/hookDocument.ts
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  getDocuments,
  getAllDocuments,
  getDocumentById,
  uploadDocument,
  deleteDocument,
  finalizeDocument,
  getDocumentStatus,
  downloadDocument,
  updateDocument,
  type StudentDocumentTable,
  type ApiResponse,
  getDocumentsByUserId,
  verifyDocument,
} from "../services/documentService";
import { toast } from "react-hot-toast";

// 🧾 Get all documents (non-paginated)
export function useDocumentList() {
  return useQuery({
    queryKey: ["documents-list"],
    queryFn: getAllDocuments,
  });
}

export const useDocuments = (page: number, limit: number) => {
  return useQuery<{ items: StudentDocumentTable[]; totalCount: number }, Error>({
    queryKey: ['documents', page, limit],
    queryFn: () => getDocuments(page, limit),
    placeholderData: (prev) => prev, // ✅ Keeps previous data
  });
};

// 🔍 Get single document by ID
export function useDocumentById(id: string) {
  return useQuery({
    queryKey: ["document", id],
    queryFn: () => getDocumentById(id),
    enabled: !!id,
  });
}

export function useUploadDocument() {
  const queryClient = useQueryClient();

  return useMutation<ApiResponse, Error, FormData>({
    mutationFn: uploadDocument,
    onSuccess: () => {
      toast.success("✅ Document uploaded successfully!");
      queryClient.invalidateQueries({ queryKey: ["documents"] });
    },
    onError: (error) => {
      toast.error(`❌ Upload failed: ${error.message}`);
    },
  });
}
// ❌ Delete document
export function useDeleteDocument() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteDocument,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["documents"] });
    },
  });
}

// ✅ Finalize document
export function useFinalizeDocument() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: finalizeDocument,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["documents"] });
    },
  });
}

// 📊 Get document completion status
export function useDocumentStatus() {
  return useQuery({
    queryKey: ["documents-status"],
    queryFn: getDocumentStatus,
  });
}

// 📥 Download document
export function useDownloadDocument(id: string) {
  return useQuery({
    queryKey: ["document-download", id],
    queryFn: () => downloadDocument(id),
    enabled: !!id,
  });
}

// ✏️ Update document
export function useUpdateDocument() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: updateDocument,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["documents"] });
    },
  });
}
export function useDocumentsByUserId(userId: string) {
  return useQuery({
    queryKey: ["documents-by-user", userId],
    queryFn: () => getDocumentsByUserId(userId),
    enabled: !!userId,
  });
}
export function useVerifyDocument() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: string;
      payload: { verification: string; remarks?: string };
    }) => verifyDocument(id, payload),
    onSuccess: (_data, _vars, _ctx) => {
      toast.success("✅ Document status updated!");
      queryClient.invalidateQueries({ queryKey: ["documents-by-user"] }); // refetch document list
    },
    onError: () => {
      toast.error("❌ Verification failed");
    },
  });
}
