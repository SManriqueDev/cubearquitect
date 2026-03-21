import { create } from 'zustand';
import { fetchPricing } from '@/services/pricingService';

export interface PricingData {
  plans: Array<{ plan_name: string; cpu: number; ram: number }>;
  locations: Array<{ location_name: string; description: string }>;
  templates: Array<{ template_name: string; os_name: string; version: string }>;
}

interface PricingState {
  pricing: PricingData | null;
  isPending: boolean;
  error: Error | null;
  fetch: () => Promise<void>;
}

export const usePricingStore = create<PricingState>((set, get) => ({
  pricing: null,
  isPending: false,
  error: null,

  fetch: async () => {
    const state = get();
    if (state.isPending || state.pricing) return;
    
    set({ isPending: true, error: null });

    try {
      const data = await fetchPricing();
      set({ pricing: data, isPending: false });
    } catch (err) {
      set({ error: err instanceof Error ? err : new Error('Failed to fetch pricing'), isPending: false });
    }
  },
}));
