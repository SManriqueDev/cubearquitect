import { useEffect, useState } from 'react';
import { fetchCanvasData } from '@/services/canvasService';
import type { CanvasData, FetchState } from '@/types/canvas';

export function useCanvasData(): FetchState {
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [data, setData] = useState<CanvasData | null>(null);

  useEffect(() => {
    let isMounted = true;

    const loadData = async () => {
      try {
        setLoading(true);
        setError(null);
        const canvasData = await fetchCanvasData();
        if (isMounted) {
          setData(canvasData);
        }
      } catch (err) {
        if (isMounted) {
          const errorMessage = err instanceof Error ? err.message : 'Failed to load canvas data';
          setError(errorMessage);
        }
      } finally {
        if (isMounted) {
          setLoading(false);
        }
      }
    };

    loadData();

    return () => {
      isMounted = false;
    };
  }, []);

  return { loading, error, data };
}
