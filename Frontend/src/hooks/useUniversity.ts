import { useQuery } from "@tanstack/react-query";
import {
  fetchAllUniversities,
  fetchUniversityPrograms,
  type UniversityProgramInfo,
  type UniversityPagination,
} from "../services/university.service";

const queryOptions = {
  refetchOnWindowFocus: true,
  refetchOnReconnect: true,
  staleTime: 0,
};

export const useUniversities = (page: number, limit: number) => {
  return useQuery<UniversityPagination>({
    queryKey: ["universities", page, limit],
    queryFn: () => fetchAllUniversities(page, limit),
    ...queryOptions,
  });
};

export const useUniversityPrograms = (universityId: string) => {
  return useQuery<UniversityProgramInfo[]>({
    queryKey: ["university-programs", universityId],
    queryFn: () => fetchUniversityPrograms(universityId),
    enabled: !!universityId,
    ...queryOptions,
  });
};
