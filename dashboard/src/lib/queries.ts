import {
  useQuery,
  useMutation,
  useQueryClient,
} from "@tanstack/react-query";
import { fetchUrls, shortenUrl, fetchStats, deleteUrl } from "./api";

export const queryKeys = {
  urls: ["urls"] as const,
  stats: (code: string) => ["stats", code] as const,
};

export function useUrls() {
  return useQuery({
    queryKey: queryKeys.urls,
    queryFn: fetchUrls,
  });
}

export function useStats(code: string) {
  return useQuery({
    queryKey: queryKeys.stats(code),
    queryFn: () => fetchStats(code),
  });
}

export function useShortenUrl() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (url: string) => shortenUrl(url),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.urls });
    },
  });
}

export function useDeleteUrl() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (code: string) => deleteUrl(code),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: queryKeys.urls });
    },
  });
}
