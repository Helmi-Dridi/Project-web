import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  assignStudentToProgram,
  unassignStudentFromProgram,
  getStudentEnrollment,
  getEnrollmentsByProgram,
  getAllStudentEnrollments,
  type EnrollmentResponse,
} from "../services/enrollment.service";
import toast from "react-hot-toast";

// Shared options
const queryOptions = {
  refetchOnWindowFocus: false,
  refetchOnReconnect: false,
  staleTime: 1000 * 60 * 10, // cache for 10 minutes
};

// ✅ Get the current student's enrollment (used in student dashboard)
export const useStudentEnrollment = () => {
  return useQuery<EnrollmentResponse>({
    queryKey: ["student-enrollment"],
    queryFn: getStudentEnrollment,
    ...queryOptions,
  });
};

// ✅ Assign a student to a program (used in admin dashboard)
export const useAssignStudent = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      userId,
      programId,
    }: {
      userId: string;
      programId: string;
    }) => assignStudentToProgram(userId, programId),

    onSuccess: (_data, variables) => {
      toast.success("Student successfully assigned to program!");
      queryClient.invalidateQueries({
        queryKey: ["program-enrollments", variables.programId],
      });
    },

    onError: (error: any) => {
      toast.error(error?.message || "Failed to assign student");
    },
  });
};

// ✅ Unassign a student (used in admin dashboard)
export const useUnassignStudent = (programId: string) => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (userId: string) => unassignStudentFromProgram(userId),

    onSuccess: () => {
      toast.success("Student successfully unassigned.");
      queryClient.invalidateQueries({
        queryKey: ["program-enrollments", programId],
      });
    },

    onError: (error: any) => {
      toast.error(error?.message || "Failed to unassign student");
    },
  });
};


// ✅ NEW: Get all enrollments for a specific program
export const useEnrollmentsByProgram = (programId: string) => {
  return useQuery<EnrollmentResponse[]>({
    queryKey: ["program-enrollments", programId],
    queryFn: () => getEnrollmentsByProgram(programId),
    enabled: !!programId,
    ...queryOptions,
  });
};

// ✅ NEW: Get all student enrollments across all programs
export const useAllStudentEnrollments = () => {
  return useQuery<EnrollmentResponse[]>({
    queryKey: ["all-enrollments"],
    queryFn: getAllStudentEnrollments,
    ...queryOptions,
  });
};
