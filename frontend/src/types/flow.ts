// Flow data types

export type NodeStatus = 'active' | 'inactive' | 'error';
export type NodeType = 'app' | 'database';

// Base node interface with index signature for React Flow compatibility
export type BaseNodeData = {
  id: string;
  type: NodeType;
  label: string;
  planName: string;
  locationName: string;
  status: NodeStatus;
  projectId: number;
  isSelected?: boolean;
  [key: string]: unknown;
};

// AppNode (VPS) - según CubePath VPSCreateRequest
export type AppNodeData = BaseNodeData & {
  type: 'app';
  name: string;
  templateName: string;
  ip?: string;
  sshKeyNames?: string[];
  ipv4: boolean;
  enableBackups: boolean;
};

// DatabaseNode (Managed DB via Cloud-init)
export type DatabaseNodeData = BaseNodeData & {
  type: 'database';
  name: string;
  ipv4: boolean;
  enableBackups: boolean;
};

// Union type for all node types
export type FlowNode = AppNodeData | DatabaseNodeData;

// Edge between nodes
export interface FlowEdge {
  id: string;
  source: string;
  target: string;
  label?: string;
  dependency?: 'network' | 'execution' | 'storage';
}

// Canvas data from API
export interface CanvasData {
  nodes: FlowNode[];
  edges: FlowEdge[];
}

// Deploy payload
export interface DeployNode {
  id: string;
  type: 'app' | 'database';
  name: string;
  plan_name: string;
  template_name?: string;
  location_name: string;
  label?: string;
  ssh_key_names?: string[];
  ipv4?: boolean;
  enable_backups?: boolean;
  custom_cloudinit?: string;
}

export interface DeployPayload {
  nodes: DeployNode[];
  edges: { source: string; target: string }[];
}

// Pricing types from /pricing
export interface VPSPlan {
  plan_name: string;
  cpu: number;
  ram: number;
  storage: number;
  bandwidth: number;
  price_per_hour: string;
}

export interface VPSTemplate {
  template_name: string;
  os_name: string;
  version: string;
}

export interface VPSLocation {
  location_name: string;
  description: string;
}

export interface PricingData {
  locations: VPSLocation[];
  plans: VPSPlan[];
  templates: VPSTemplate[];
}
