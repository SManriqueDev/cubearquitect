import { memo } from 'react';
import { Button } from '@/components/ui/button';
import { Server, Database, Rocket } from 'lucide-react';
import { useFlowStore } from '@/stores/flowStore';
import { useDeploy } from '@/hooks/useDeploy';
import { createDeployPayload } from '@/utils/nodeUtils';

interface FlowToolbarProps {
  onAddNode?: (type: 'app' | 'database') => void;
  nodeCount?: { total: number; apps: number; databases: number };
  edgeCount?: number;
}

function FlowToolbarComponent({ onAddNode, nodeCount, edgeCount = 0 }: FlowToolbarProps) {
  const nodes = useFlowStore((state) => state.nodes);
  const edges = useFlowStore((state) => state.edges);

  const { mutate: deploy, isPending: isDeploying } = useDeploy();

  const handleAddApp = () => {
    onAddNode?.('app');
  };

  const handleAddDatabase = () => {
    onAddNode?.('database');
  };

  const handleDeploy = () => {
    const payload = createDeployPayload(nodes, edges);
    deploy(payload);
  };

  const hasNodes = nodes.length > 0;

  return (
    <div className="absolute top-4 left-4 z-10 flex gap-2 bg-white rounded-lg shadow-md p-2">
      <Button
        variant="outline"
        size="sm"
        onClick={handleAddApp}
        className="gap-2"
      >
        <Server className="w-4 h-4" />
        Add App
      </Button>

      <Button
        variant="outline"
        size="sm"
        onClick={handleAddDatabase}
        className="gap-2"
      >
        <Database className="w-4 h-4" />
        Add DB
      </Button>

      <div className="w-px bg-gray-200 mx-1" />

      <Button
        variant="default"
        size="sm"
        onClick={handleDeploy}
        disabled={!hasNodes || isDeploying}
        className="gap-2 bg-green-600 hover:bg-green-700"
      >
        <Rocket className="w-4 h-4" />
        {isDeploying ? 'Deploying...' : 'Deploy'}
      </Button>

      {nodeCount && (
        <div className="flex items-center gap-2 text-xs text-gray-500 ml-2">
          <span>{nodeCount.apps} apps</span>
          <span>|</span>
          <span>{nodeCount.databases} DBs</span>
          <span>|</span>
          <span>{edgeCount} edges</span>
        </div>
      )}
    </div>
  );
}

export const FlowToolbar = memo(FlowToolbarComponent);
