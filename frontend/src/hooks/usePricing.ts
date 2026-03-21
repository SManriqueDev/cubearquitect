import { useQuery } from '@tanstack/react-query';
import { fetchPricing } from '@/services/pricingService';

export function usePricing() {
  return useQuery({
    queryKey: ['pricing'],
    queryFn: fetchPricing,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}
