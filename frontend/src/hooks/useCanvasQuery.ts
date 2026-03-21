import { useQuery } from '@tanstack/react-query';
import { fetchCanvasData } from '@/services/canvasService';
import type { CanvasData } from '@/types/canvas';

interface UseCanvasQueryResult {
  data: CanvasData | undefined;
  isPending: boolean;
  error: Error | null;
}

export function useCanvasQuery(): UseCanvasQueryResult {
  const { data, isPending, error } = useQuery({
    queryKey: ['canvas'],
    queryFn: fetchCanvasData,
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000, // 10 minutes (formerly cacheTime)
    retry: 1,
    refetchOnWindowFocus: true,
  });

  return {
    data,
    isPending,
    error: error as Error | null,
  };
}
