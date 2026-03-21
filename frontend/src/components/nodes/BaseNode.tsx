import { memo, type ReactNode, type FC } from 'react';
import { Handle, Position } from '@xyflow/react';
import { Card } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import type { FlowNode } from '@/types/flow';
import { getStatusColor } from '@/utils/nodeUtils';

interface BaseNodeProps {
  data: FlowNode;
  icon: ReactNode;
  iconColor: string;
  children?: ReactNode;
}

const BaseNodeComponent: FC<BaseNodeProps> = ({
  data,
  icon,
  iconColor,
  children,
}) => {
  const statusColor = getStatusColor(data.status);

  return (
    <div className={`transition-all ${data.isSelected ? 'ring-2 ring-blue-500 ring-offset-2' : ''}`}>
      <Card className="p-4 w-56">
        <div className="space-y-3">
          {/* Header */}
          <div className="flex items-center justify-between gap-2">
            <div className="flex items-center gap-2 flex-1 min-w-0">
              <span className={iconColor}>{icon}</span>
              <h3 className="font-bold text-sm truncate">{data.label}</h3>
            </div>
            <Badge
              variant="secondary"
              className={`text-xs ${statusColor}`}
            >
              {data.status}
            </Badge>
          </div>

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
