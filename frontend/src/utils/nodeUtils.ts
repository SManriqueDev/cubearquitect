import type { FlowNode, DeployPayload, DeployNode } from '@/types/flow';

export function toDeployNode(node: FlowNode): DeployNode {
  if (node.type === 'app') {
    return {
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
    type: 'database',
    name: node.name,
    plan_name: node.planName,
    location_name: node.locationName,
    label: node.label,
    custom_cloudinit: node.cloudInitConfig,
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

export function getStatusColor(status: string): string {
  switch (status) {
    case 'active':
      return 'text-green-500';
    case 'inactive':
      return 'text-gray-400';
    case 'error':
      return 'text-red-500';
    default:
      return 'text-gray-400';
  }
}
