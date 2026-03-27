import { memo } from 'react';
import { Rocket, Server, Database, Terminal } from 'lucide-react';
import { useFlowStore } from '@/stores/flowStore';
import {
  ToolbarButton,
  ToolbarAddButton,
  ToolbarSeparator,
  ToolbarStats,
} from './components';

interface FlowToolbarProps {
  onAddNode?: (type: 'app' | 'database') => void;
  onDeploy?: () => void;
  isDeploying?: boolean;
}

export const FlowToolbar = memo(function FlowToolbar({
  onAddNode,
  onDeploy,
  isDeploying = false,
}: FlowToolbarProps) {
  const nodes = useFlowStore((state) => state.nodes);
  const edges = useFlowStore((state) => state.edges);
  const showLogs = useFlowStore((state) => state.showLogs);
  const setShowLogs = useFlowStore((state) => state.setShowLogs);

  const appCount = nodes.filter((n) => n.type === 'app').length;
  const dbCount = nodes.filter((n) => n.type === 'database').length;
  const hasNodes = nodes.length > 0;

  return (
    <div
      className="absolute top-4 left-4 z-10 flex items-center gap-2 bg-card rounded-lg border border-border/60 shadow-sm p-2"
      role="toolbar"
      aria-label="Flow editor toolbar"
    >
      <ToolbarAddButton onAddNode={onAddNode} />

      <ToolbarSeparator />

      <ToolbarButton
        icon={Rocket}
        label={isDeploying ? 'Deploying...' : 'Deploy'}
        variant="deploy"
        onClick={onDeploy}
        disabled={!hasNodes || isDeploying}
        loading={isDeploying}
      />

      <ToolbarSeparator />

      <ToolbarButton
        icon={Terminal}
        label={showLogs ? 'Hide Logs' : 'Show Logs'}
        onClick={() => setShowLogs(!showLogs)}
        variant={showLogs ? 'default' : 'ghost'}
      />

      <ToolbarSeparator />

      <ToolbarStats
        nodes={[
          { type: 'app', label: 'Apps', icon: Server, count: appCount },
          { type: 'database', label: 'DBs', icon: Database, count: dbCount },
        ]}
        edgeCount={edges.length}
      />
    </div>
  );
});