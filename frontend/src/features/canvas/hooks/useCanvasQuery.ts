import { useQuery } from '@tanstack/react-query';
import { fetchCanvasData } from '@/features/canvas/api/fetchCanvasData';
import { canvasKeys } from '@/features/canvas/api/queryKeys';
import type { CanvasData } from '@/types/canvas';

interface UseCanvasQueryResult {
  data: CanvasData | undefined;
  isPending: boolean;
  error: Error | null;
}

export function useCanvasQuery(): UseCanvasQueryResult {
  const { data, isPending, error } = useQuery({
    queryKey: canvasKeys.all,
    queryFn: fetchCanvasData,
    refetchOnWindowFocus: true,
  });

  return {
    data,
    isPending,
    error: error as Error | null,
  };
}

