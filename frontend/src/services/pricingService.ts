import { apiFetch } from './api';
import type { PricingData } from '@/types/flow';

export async function fetchPricing(): Promise<PricingData> {
  const data = await apiFetch<{
    vps: {
      locations: Array<{
        location_name: string;
        description: string;
        clusters: Array<{
          plans: Array<{
            plan_name: string;
            cpu: number;
            ram: number;
            storage: number;
            bandwidth: number;
            price_per_hour: string;
          }>;
        }>;
        templates: Array<{
          template_name: string;
          os_name: string;
          version: string;
        }>;
      }>;
    };
  }>('/api/pricing');

  // Flatten nested structure from API
  const locations = data.vps.locations.map((loc) => ({
    location_name: loc.location_name,
    description: loc.description,
  }));

  const plans: PricingData['plans'] = [];
  const templates: PricingData['templates'] = [];

  data.vps.locations.forEach((loc) => {
    loc.clusters.forEach((cluster) => {
      cluster.plans.forEach((plan) => {
        if (!plans.some((p) => p.plan_name === plan.plan_name)) {
          plans.push(plan);
        }
      });
    });
    loc.templates.forEach((tmpl) => {
      if (!templates.some((t) => t.template_name === tmpl.template_name)) {
        templates.push(tmpl);
      }
    });
  });

  return { locations, plans, templates };
}
