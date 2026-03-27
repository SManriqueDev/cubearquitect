import type { CanvasData, FlowNode, FlowEdge, NodeStatus } from '@/types/flow';

interface FloatingIP {
  address: string;
  type: 'IPv4' | 'IPv6';
  is_primary: boolean;
}

interface VPSPlan {
  plan_name: string;
  cpu: number;
  ram: number;
  storage: number;
}

interface VPSItem {
  id: number;
  name: string;
  status: string;
  user: string;
  label: string;
  plan: VPSPlan;
  floating_ips?: {
    list: FloatingIP[];
  };
  hostname?: string;
  location?: {
    location_name: string;
    description: string;
  };
  node_type?: string;
}

interface ProjectResponse {
  project: {
    id: number;
    name: string;
    description: string;
  };
  networks: unknown[];
  baremetals: unknown[];
  vps: VPSItem[];
}

export async function fetchCanvasData(): Promise<CanvasData> {
  try {
    const response = await fetch('/api/projects', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch projects: ${response.statusText}`);
    }

    const projects: ProjectResponse[] = await response.json();
    const nodes: FlowNode[] = [];
    const edges: FlowEdge[] = [];

    projects.forEach((project) => {
      const projectId = project.project.id;

      project.vps.forEach((vps) => {
        const nodeId = `vps-${projectId}-${vps.id}`;
        const primaryIP =
          vps.floating_ips?.list.find((ip) => ip.is_primary && ip.type === 'IPv4')
            ?.address ||
          vps.floating_ips?.list[0]?.address ||
          '';

        const nodeType = vps.node_type === 'database' ? 'database' : 'app';

        nodes.push({
          id: nodeId,
          type: nodeType,
          name: vps.name || vps.label,
          label: vps.label || vps.name,
          planName: vps.plan?.plan_name || 'default',
          templateName: 'ubuntu-24', // Default template
          locationName: vps.location?.location_name || 'us-mia-1',
          ip: primaryIP,
          status: mapVPSStatus(vps.status),
          ipv4: true,
          enableBackups: false,
          projectId,
        });
      });
    });

    return { nodes, edges };
  } catch (error) {
    const message = error instanceof Error ? error.message : 'Unknown error';
    throw new Error(`Canvas data fetch error: ${message}`);
  }
}

function mapVPSStatus(status: string): NodeStatus {
  const normalizedStatus = status.toLowerCase();
  if (normalizedStatus.includes('active') || normalizedStatus.includes('running') || normalizedStatus.includes('healthy')) {
    return 'active';
  }
  if (normalizedStatus.includes('error') || normalizedStatus.includes('failed')) {
    return 'error';
  }
  return 'pending';
}
