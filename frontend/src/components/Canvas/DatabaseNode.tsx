import { memo } from 'react';
import { Handle, Position } from '@xyflow/react';
import type { NodeProps } from '@xyflow/react';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Database, CheckCircle, AlertCircle, Clock } from 'lucide-react';
import type { CanvasNode } from '@/types/canvas';

const statusIcons: Record<string, React.ComponentType<{ className?: string }>> = {
  active: CheckCircle,
  inactive: Clock,
  error: AlertCircle,
};

const statusColors: Record<string, string> = {
  active: 'text-green-500',
  inactive: 'text-gray-400',
  error: 'text-red-500',
};

function DatabaseNodeComponent({
  data,
}: NodeProps) {
  const nodeData = data as CanvasNode;
  const StatusIcon = statusIcons[nodeData.status] || AlertCircle;
  const statusColor = statusColors[nodeData.status] || 'text-gray-400';

  return (
    <div className={`transition-all ${nodeData.isSelected ? 'ring-2 ring-blue-500' : ''}`}>
      <Card className="p-4 w-48 border-l-4 border-l-purple-600">
        <div className="space-y-2">
          <div className="flex items-center justify-between gap-2">
            <div className="flex items-center gap-2 flex-1 min-w-0">
              <Database className="w-4 h-4 flex-shrink-0 text-purple-600" />
              <h3 className="font-bold text-sm truncate">{nodeData.label}</h3>
            </div>
            <StatusIcon className={`w-4 h-4 flex-shrink-0 ${statusColor}`} />
          </div>

          {nodeData.ip && (
            <p className="text-xs text-gray-600 font-mono">{nodeData.ip}</p>
          )}

          <div className="flex gap-2 flex-wrap">
            <Badge variant="outline" className="text-xs">
              {nodeData.plan}
            </Badge>
            <Badge variant="secondary" className="text-xs">
              {nodeData.status}
            </Badge>
          </div>
        </div>
      </Card>

      <Handle type="target" position={Position.Top} />
      <Handle type="source" position={Position.Bottom} />
    </div>
  );
}

const DatabaseNode = memo(DatabaseNodeComponent);

export default DatabaseNode;
