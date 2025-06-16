// hooks/useAdmin.ts
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  fetchAdmins,
  uploadProfileImage,
  fetchAllAdminDetails,
  createAdmin,
  deleteAdmin,
  assignAdminRole,
  fetchUsersWithoutRole,
  fetchStudents,
} from "../services/admin.service";
import { toast } from "react-hot-toast";

// ✅ Fetch simplified list (e.g. for dropdowns)
export const useAdmins = () => {
  return useQuery({
    queryKey: ["admins"],
    queryFn: fetchAdmins,
    staleTime: 1000 * 60 * 5, // ✅ Cache for 5 minutes
  });
};

// ✅ Fetch full admin details
export const useAdminDetails = () => {
  return useQuery({
    queryKey: ["admin-details"],
    queryFn: fetchAllAdminDetails,
    staleTime: 1000 * 60 * 5, // ✅ Avoid repeated calls
  });
};

// ✅ Upload profile image
export const useUploadProfileImage = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: uploadProfileImage,
    onSuccess: () => {
      toast.success("Profile image updated 🎉");
      queryClient.invalidateQueries({ queryKey: ["admin-details"] });
    },
    onError: () => {
      toast.error("Failed to upload profile image ❌");
    },
  });
};

// ✅ Create admin user (update cache manually)
export const useCreateAdmin = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: createAdmin,
    onSuccess: (newAdmin) => {
      toast.success("Admin created ✅");
      queryClient.setQueryData<any[]>(["admin-details"], (prev) =>
        prev ? [...prev, newAdmin] : [newAdmin]
      );
    },
    onError: () => {
      toast.error("Failed to create admin ❌");
    },
  });
};

// ✅ Assign role to admin (no need to invalidate entire list)
export const useAssignAdminRole = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ userId, roleId }: { userId: string; roleId: string }) =>
      assignAdminRole(userId, roleId),
    onSuccess: () => {
      toast.success("Role assigned successfully ✅");

      // Optionally update relevant cache keys
      queryClient.invalidateQueries({ queryKey: ["admin-details"] });
      queryClient.invalidateQueries({ queryKey: ["users-without-role"] });
    },
    onError: () => {
      toast.error("Failed to assign role ❌");
    },
  });
};

// ✅ Delete admin user (update cache manually)
export const useDeleteAdmin = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (userId: string) => deleteAdmin(userId),
    onSuccess: (_, userId) => {
      toast.success("Admin deleted 🗑️");
      queryClient.setQueryData<any[]>(["admin-details"], (prev) =>
        prev ? prev.filter((admin) => admin.id !== userId) : []
      );
    },
    onError: () => {
      toast.error("Failed to delete admin ❌");
    },
  });
};

// ✅ Fetch users without any assigned role
export const useUsersWithoutRole = () => {
  return useQuery({
    queryKey: ["users-without-role"],
    queryFn: fetchUsersWithoutRole,
    staleTime: 1000 * 60 * 5, // ✅ Reduce repeated calls
  });
};

// ✅ Fetch users with the Student role
export const useStudents = () => {
  return useQuery({
    queryKey: ["students"],
    queryFn: fetchStudents,
    staleTime: 1000 * 60 * 10, // ✅ Longer caching
  });
};
