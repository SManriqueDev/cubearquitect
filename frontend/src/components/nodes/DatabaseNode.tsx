import { memo } from 'react';
import { Database } from 'lucide-react';
import { BaseNode } from './BaseNode';
import type { DatabaseNodeData } from '@/types/flow';

interface DatabaseNodeProps {
  data: DatabaseNodeData;
}

function DatabaseNodeComponent({ data }: DatabaseNodeProps) {
  return (
    <BaseNode
      data={data}
      icon={<Database className="w-4 h-4" />}
      iconColor="text-purple-600"
    />
  );
}

const DatabaseNode = memo(DatabaseNodeComponent);

export default DatabaseNode;
