import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { apiFetch } from '@/services/api';
import { useAccountStore } from '@/stores/accountStore';
import { canvasKeys } from '@/services/queryKeys';
import { deployPayloadSchema } from '@/services/schemas/flow';

type DeployResponse = {
  success: boolean;
  message: string;
  deployment_id?: string;
};

interface UseDeployOptions {
  onDeployStarted?: (deploymentId: string, nodeIds: string[]) => void;
  onSuccess?: (data: DeployResponse) => void;
  onError?: (error: Error) => void;
}

export function useDeploy(options: UseDeployOptions = {}) {
  const queryClient = useQueryClient();
  const { projectId, selectedSSHKeys } = useAccountStore();

  return useMutation({
    mutationFn: async (payload: unknown) => {
      const validated = deployPayloadSchema.parse(payload);
      
      const deployPayload = {
        ...validated,
        project_id: projectId,
        ssh_key_names: selectedSSHKeys,
      };

      return apiFetch<DeployResponse>('/api/deploy', {
        method: 'POST',
        body: JSON.stringify(deployPayload),
      });
    },
    onSuccess: (data, variables) => {
      toast.success('Deployment started!', {
        description: data.message,
      });
      
      if (data.deployment_id) {
        const parsed = deployPayloadSchema.safeParse(variables);
        const nodeIds = parsed.success ? parsed.data.nodes.map((n) => n.id) : [];
        options.onDeployStarted?.(data.deployment_id, nodeIds);
      }
      
      queryClient.invalidateQueries({ queryKey: canvasKeys.all });
      options.onSuccess?.(data);
    },
    onError: (error) => {
      const message = error instanceof Error ? error.message : 'Unknown error occurred';
      toast.error('Deployment failed', {
        description: message,
      });
      options.onError?.(error);
    },
  });
}

type UseDeployReturn = ReturnType<typeof useDeploy>;
export type { UseDeployReturn };
