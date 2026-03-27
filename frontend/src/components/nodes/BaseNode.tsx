import { memo, type ReactNode, type FC } from 'react';
import { Handle, Position } from '@xyflow/react';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Loader2, Clock, AlertCircle } from 'lucide-react';
import type { FlowNode, NodeStatus } from '@/types/flow';
import { getStatusColor } from '@/utils/nodeUtils';

interface BaseNodeProps {
  data: FlowNode;
  icon: ReactNode;
  iconColor: string;
  children?: ReactNode;
}

function StatusIndicator({ status }: { status: NodeStatus }) {
  switch (status) {
    case 'deploying':
      return (
        <div className="flex items-center gap-1">
          <Loader2 className="h-3 w-3 animate-spin" />
          <span className="text-xs">Deploying...</span>
        </div>
      );
    case 'pending':
      return (
        <div className="flex items-center gap-1">
          <Clock className="h-3 w-3" />
          <span className="text-xs">Pending...</span>
        </div>
      );
    case 'error':
      return (
        <div className="flex items-center gap-1">
          <AlertCircle className="h-3 w-3" />
          <span className="text-xs">Error</span>
        </div>
      );
    default:
      return null;
  }
}

const BaseNodeComponent: FC<BaseNodeProps> = ({
  data,
  icon,
  iconColor,
  children,
}) => {
  const statusColor = getStatusColor(data.status);

  return (
    <div className={`transition-all duration-300 ${data.isSelected ? 'ring-2 ring-blue-500 ring-offset-2' : ''}`}>
      <Card className="p-4 w-56 transition-all duration-300 relative">
        {/* Indicador circular animado - esquina superior derecha */}
        {(data.status === 'pending' || data.status === 'deploying' || data.status === 'active' || data.status === 'error') && (
          <div className="absolute top-2 right-2">
            <span className="relative flex h-3 w-3">
              <span className={`absolute inline-flex h-full w-full rounded-full opacity-75 animate-ping ${
                data.status === 'pending' ? 'bg-yellow-400' :
                data.status === 'deploying' ? 'bg-blue-400' :
                data.status === 'active' ? 'bg-green-400' :
                'bg-red-400'
              }`} />
              <span className={`relative inline-flex rounded-full h-3 w-3 ${
                data.status === 'pending' ? 'bg-yellow-500' :
                data.status === 'deploying' ? 'bg-blue-500' :
                data.status === 'active' ? 'bg-green-500' :
                'bg-red-500'
              }`} />
            </span>
          </div>
        )}

        <div className="space-y-3">
          {/* Header */}
          <div className="flex items-center justify-between gap-2">
            <div className="flex items-center gap-2 flex-1 min-w-0">
              <span className={iconColor}>{icon}</span>
              <h3 className="font-bold text-sm truncate">{data.name}</h3>
            </div>
            <div className={`text-xs ${statusColor}`}>
              <StatusIndicator status={data.status} />
            </div>
          </div>

          {/* Error message */}
          {data.status === 'error' && data.errorMessage && (
            <div className="text-xs text-red-600 bg-red-100 rounded p-2 border border-red-200 animate-pulse">
              {data.errorMessage}
            </div>
          )}

          {/* Content */}
          {children}

          {/* Footer */}
          <div className="flex gap-2 flex-wrap">
            <Badge variant="outline" className="text-xs">
              {data.planName}
            </Badge>
            <Badge variant="outline" className="text-xs">
              {data.locationName}
            </Badge>
          </div>
        </div>
      </Card>

      <Handle type="target" position={Position.Left} className="!bg-gray-400" />
      <Handle type="source" position={Position.Right} className="!bg-gray-400" />
    </div>
  );
};

export const BaseNode = memo(BaseNodeComponent);
