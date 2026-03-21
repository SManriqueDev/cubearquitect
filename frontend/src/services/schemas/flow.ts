import { z } from 'zod';

export const nodeTypeSchema = z.enum(['app', 'database']);
export const nodeStatusSchema = z.enum(['active', 'inactive', 'error']);

export const baseNodeSchema = z.object({
  id: z.string().min(1),
  type: nodeTypeSchema,
  label: z.string().min(1),
  planName: z.string().min(1),
  locationName: z.string().min(1),
  status: nodeStatusSchema,
  projectId: z.number().int().positive(),
  isSelected: z.boolean().optional(),
});

export const appNodeSchema = baseNodeSchema.extend({
  type: z.literal('app'),
  name: z.string().min(1).max(63),
  templateName: z.string().min(1),
  ip: z.string().optional(),
  sshKeyNames: z.array(z.string()).optional(),
  ipv4: z.boolean(),
  enableBackups: z.boolean(),
});

export const databaseNodeSchema = baseNodeSchema.extend({
  type: z.literal('database'),
  name: z.string().min(1).max(63),
  cloudInitConfig: z.string().optional(),
});

export const flowNodeSchema = z.union([appNodeSchema, databaseNodeSchema]);

export const flowEdgeSchema = z.object({
  id: z.string().min(1),
  source: z.string().min(1),
  target: z.string().min(1),
  label: z.string().optional(),
  dependency: z.enum(['network', 'execution', 'storage']).optional(),
});

export const canvasDataSchema = z.object({
  nodes: z.array(flowNodeSchema),
  edges: z.array(flowEdgeSchema),
});

export const deployNodeSchema = z.object({
  type: nodeTypeSchema,
  name: z.string().min(1).max(63),
  plan_name: z.string().min(1),
  template_name: z.string().optional(),
  location_name: z.string().min(1),
  label: z.string().optional(),
  ssh_key_names: z.array(z.string()).optional(),
  ipv4: z.boolean().optional(),
  enable_backups: z.boolean().optional(),
  custom_cloudinit: z.string().optional(),
});

export const deployPayloadSchema = z.object({
  nodes: z.array(deployNodeSchema),
  edges: z.array(z.object({
    source: z.string().min(1),
    target: z.string().min(1),
  })),
});
