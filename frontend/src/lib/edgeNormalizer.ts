import type { Edge as ReactFlowEdge } from '@xyflow/react';

// Convert canvas edges to React Flow edge format
export function normalizeEdges(edges: Array<{
  id: string;
  source: string;
  target: string;
  label?: string;
  dependency: string;
}>) {
  return edges.map((edge) => ({
    id: edge.id,
    source: edge.source,
    target: edge.target,
    label: edge.label,
    animated: edge.dependency === 'execution',
    style: {
      stroke: getDependencyColor(edge.dependency),
      strokeWidth: 2,
    },
  } as ReactFlowEdge));
}

function getDependencyColor(dependency: string): string {
  const colors: Record<string, string> = {
    execution: '#3b82f6', // blue for execution order
    network: '#10b981',   // green for network dependencies
    storage: '#f59e0b',   // amber for storage dependencies
  };
  return colors[dependency] || '#6b7280';
}
