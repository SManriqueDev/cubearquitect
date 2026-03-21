// Canvas data types for Phase 2

export type NodeStatus = 'active' | 'inactive' | 'error';

export interface CanvasNode extends Record<string, unknown> {
  id: string;
  type: 'app' | 'database';
  label: string;
  ip?: string;
  plan: string;
  status: NodeStatus;
  dockerImage?: string;
  region?: string;
  size?: string;
  projectId: number;
  isSelected?: boolean;
}

export interface CanvasEdge {
  id: string;
  source: string;
  target: string;
  label?: string;
  dependency: 'network' | 'execution' | 'storage';
}

export interface CanvasData {
  nodes: CanvasNode[];
  edges: CanvasEdge[];
}

export interface FetchState {
  loading: boolean;
  error: string | null;
  data: CanvasData | null;
}
