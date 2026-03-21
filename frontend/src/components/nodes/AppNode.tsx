import { memo } from 'react';
import { Server } from 'lucide-react';
import { BaseNode } from './BaseNode';
import type { AppNodeData } from '@/types/flow';

interface AppNodeProps {
  data: AppNodeData;
}

function AppNodeComponent({ data }: AppNodeProps) {
  return (
    <BaseNode
      data={data}
      icon={<Server className="w-4 h-4" />}
      iconColor="text-blue-500"
    >
      {data.ip && (
        <p className="text-xs text-gray-600 font-mono">{data.ip}</p>
      )}
      <p className="text-xs text-gray-500 truncate">{data.templateName}</p>
    </BaseNode>
  );
}

const AppNode = memo(AppNodeComponent);

export default AppNode;
