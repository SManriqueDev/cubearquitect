import { useQuery } from '@tanstack/react-query';
import { fetchCanvasData } from '@/services/fetchProjects';
import { canvasKeys } from '@/services/queryKeys';
import type { CanvasData } from '@/types/flow';

interface UseProjectsResult {
  data: CanvasData | undefined;
  isPending: boolean;
  error: Error | null;
}

export function useProjects(): UseProjectsResult {
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

