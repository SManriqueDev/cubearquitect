import { useMemo } from 'react';
import { useFlowStore } from '@/stores/flowStore';

export function useNodeCount(): { total: number; apps: number; databases: number } {
  const nodes = useFlowStore((state) => state.nodes);

  return useMemo(() => ({
    total: nodes.length,
    apps: nodes.filter((n) => n.type === 'app').length,
    databases: nodes.filter((n) => n.type === 'database').length,
  }), [nodes]);
}

export function useEdgeCount(): number {
  return useFlowStore((state) => state.edges.length);
}
