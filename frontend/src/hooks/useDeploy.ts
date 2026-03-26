import { useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import { apiFetch } from '@/services/api';
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

  return useMutation({
    mutationFn: async (payload: unknown) => {
      const validated = deployPayloadSchema.parse(payload);
      return apiFetch<DeployResponse>('/api/deploy', {
        method: 'POST',
        body: JSON.stringify(validated),
      });
    },
    onSuccess: (data) => {
      toast.success('Deployment started!', {
        description: data.message,
      });
      
      if (data.deployment_id) {
        options.onDeployStarted?.(data.deployment_id, []);
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
