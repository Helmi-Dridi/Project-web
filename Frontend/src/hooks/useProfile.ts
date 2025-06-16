import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "react-hot-toast";

import {
  getStudentProfile,
  createOrUpdateStudentProfile,
  getStudentContact,
  createOrUpdateStudentContact,
  getStudentAcademic,
  createOrUpdateStudentAcademic,
  getStudentSettings,
  createOrUpdateStudentSettings,
  getStudentProfileByUserId,
  getStudentContactByUserId,
  getStudentAcademicByUserId,
  getStudentSettingsByUserId,
} from "../services/profile.service";

// Shared query config
const queryOptions = {
  retry: false,
  refetchOnWindowFocus: false,
  refetchOnMount: false,
  refetchOnReconnect: false,
};

// -----------------------------
// Student Profile Hook
// -----------------------------
export const useStudentProfile = () => {
  const queryClient = useQueryClient();

  const query = useQuery({
    queryKey: ["studentProfile"],
    queryFn: getStudentProfile,
    ...queryOptions,
  });

  const mutation = useMutation({
    mutationFn: createOrUpdateStudentProfile,
    onSuccess: () => {
      toast.success("Profile saved successfully ✅");
      queryClient.invalidateQueries({ queryKey: ["studentProfile"] });
    },
    onError: () => toast.error("Failed to save profile ❌"),
  });

  return {
    ...query,
    upsert: mutation.mutate,
    isSaving: mutation.isPending,
  };
};

// Admin version: Get profile by user ID
export const useStudentProfileByUserId = (userId: string) => {
  return useQuery({
    queryKey: ["studentProfile", userId],
    queryFn: () => getStudentProfileByUserId(userId),
    ...queryOptions,
  });
};

// -----------------------------
// Student Contact Hook
// -----------------------------
export const useStudentContact = () => {
  const queryClient = useQueryClient();

  const query = useQuery({
    queryKey: ["studentContact"],
    queryFn: getStudentContact,
    ...queryOptions,
  });

  const mutation = useMutation({
    mutationFn: createOrUpdateStudentContact,
    onSuccess: () => {
      toast.success("Contact info saved ✅");
      queryClient.invalidateQueries({ queryKey: ["studentContact"] });
    },
    onError: () => toast.error("Failed to save contact info ❌"),
  });

  return {
    ...query,
    upsert: mutation.mutate,
    isSaving: mutation.isPending,
  };
};

// Admin version: Get contact by user ID
export const useStudentContactByUserId = (userId: string) => {
  return useQuery({
    queryKey: ["studentContact", userId],
    queryFn: () => getStudentContactByUserId(userId),
    ...queryOptions,
  });
};

// -----------------------------
// Student Academic Hook
// -----------------------------
export const useStudentAcademic = () => {
  const queryClient = useQueryClient();

  const query = useQuery({
    queryKey: ["studentAcademic"],
    queryFn: getStudentAcademic,
    ...queryOptions,
  });

  const mutation = useMutation({
    mutationFn: createOrUpdateStudentAcademic,
    onSuccess: () => {
      toast.success("Academic info saved 🎓");
      queryClient.invalidateQueries({ queryKey: ["studentAcademic"] });
    },
    onError: () => toast.error("Failed to save academic info ❌"),
  });

  return {
    ...query,
    upsert: mutation.mutate,
    isSaving: mutation.isPending,
  };
};

// Admin version: Get academic by user ID
export const useStudentAcademicByUserId = (userId: string) => {
  return useQuery({
    queryKey: ["studentAcademic", userId],
    queryFn: () => getStudentAcademicByUserId(userId),
    ...queryOptions,
  });
};

// -----------------------------
// Student Settings Hook
// -----------------------------
export const useStudentSettings = () => {
  const queryClient = useQueryClient();

  const query = useQuery({
    queryKey: ["studentSettings"],
    queryFn: getStudentSettings,
    ...queryOptions,
  });

  const mutation = useMutation({
    mutationFn: createOrUpdateStudentSettings,
    onSuccess: () => {
      toast.success("Settings saved ⚙️");
      queryClient.invalidateQueries({ queryKey: ["studentSettings"] });
    },
    onError: () => toast.error("Failed to save settings ❌"),
  });

  return {
    ...query,
    upsert: mutation.mutate,
    isSaving: mutation.isPending,
  };
};

// Admin version: Get settings by user ID
export const useStudentSettingsByUserId = (userId: string) => {
  return useQuery({
    queryKey: ["studentSettings", userId],
    queryFn: () => getStudentSettingsByUserId(userId),
    ...queryOptions,
  });
};
