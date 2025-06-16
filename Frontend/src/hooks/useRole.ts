// src/hooks/useRole.ts
import { useQuery } from "@tanstack/react-query";
import { fetchRoleall } from "../services/role.service";

// ✅ Hook to fetch all roles for current company
export const useRoles = () => {
  return useQuery({
    queryKey: ["roles"],
    queryFn: fetchRoleall,
    staleTime: 1000 * 60 * 60, // ✅ Cache roles for 1 hour
    refetchOnWindowFocus: false, // ✅ Prevent re-fetch on tab switch
  });
};
