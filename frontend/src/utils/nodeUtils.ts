import type { FlowNode, DeployPayload, DeployNode, NodeStatus } from '@/types/flow';

export function toDeployNode(node: FlowNode): DeployNode {
  if (node.type === 'app') {
    return {
      id: node.id,
      type: 'app',
      name: node.name,
      plan_name: node.planName,
      template_name: node.templateName,
      location_name: node.locationName,
      label: node.label,
      ssh_key_names: node.sshKeyNames,
      ipv4: node.ipv4,
      enable_backups: node.enableBackups,
    };
  }

  return {
    id: node.id,
    type: 'database',
    name: node.name,
    plan_name: node.planName,
    location_name: node.locationName,
    label: node.label,
    ipv4: node.ipv4,
    enable_backups: node.enableBackups,
  };
}

export function createDeployPayload(
  nodes: FlowNode[],
  edges: { source: string; target: string }[]
): DeployPayload {
  return {
    nodes: nodes.map(toDeployNode),
    edges,
  };
}

export function getStatusColor(status: NodeStatus): string {
  switch (status) {
    case 'active':
      return 'text-green-500';
    case 'error':
      return 'text-red-500';
    case 'pending':
      return 'text-yellow-500';
    case 'deploying':
      return 'text-blue-500';
    default:
      return 'text-gray-400';
  }
}

export function getStatusLabel(status: NodeStatus): string {
  switch (status) {
    case 'active':
      return 'Active';
    case 'error':
      return 'Error';
    case 'pending':
      return 'Pending';
    case 'deploying':
      return 'Deploying';
    default:
      return status;
  }
}
