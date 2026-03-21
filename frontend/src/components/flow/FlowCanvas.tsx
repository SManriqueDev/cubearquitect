import { memo, useCallback, useEffect } from 'react';
import {
  ReactFlow,
  Background,
  Controls,
  useNodesState,
  useEdgesState,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import type { Node, Edge, NodeTypes } from '@xyflow/react';

import { useProjects } from '@/hooks/useProjects';
import { useFlowStore } from '@/stores/flowStore';
import AppNode from '@/components/nodes/AppNode';
import DatabaseNode from '@/components/nodes/DatabaseNode';
import { ConfigurationPanel } from './ConfigurationPanel';
import { normalizeEdges } from '@/lib/edgeNormalizer';
import type { CanvasNode } from '@/types/flow';

const nodeTypes: NodeTypes = {
  app: AppNode,
  database: DatabaseNode,
};

function FlowCanvasComponent() {
  const { data, isPending, error } = useProjects();
  const selectedNodeId = useFlowStore((state) => state.selectedNodeId);
  const setSelectedNodeId = useFlowStore((state) => state.setSelectedNodeId);

  const emptyNodes: Node[] = [];
  const emptyEdges: Edge[] = [];
  
  const [nodes, setNodes, onNodesChange] = useNodesState(emptyNodes);
  const [, , onEdgesChange] = useEdgesState(emptyEdges);

  // Optimize re-renders: only when data.nodes changes, not selectedNode
  useEffect(() => {
    if (!data?.nodes || data.nodes.length === 0) {
      setNodes([]);
      return;
    }

    const newNodes: Node[] = data.nodes.map((node, idx) => ({
      id: node.id,
      data: { ...node, isSelected: selectedNodeId === node.id },
      position: {
        x: (idx % 3) * 400,
        y: Math.floor(idx / 3) * 400,
      },
      type: node.type,
    }));

    setNodes(newNodes);
  }, [data?.nodes, setNodes, selectedNodeId]);

  const selectedNode: CanvasNode | null =
    data?.nodes?.find((item) => item.id === selectedNodeId) ?? null;

  const handleNodeClick = useCallback(
    (_: React.MouseEvent, node: Node) => {
      setSelectedNodeId(node.id);
      setNodes((nds) =>
        nds.map((n) => {
          const currentData = n.data as CanvasNode;
          return {
            ...n,
            data: { ...currentData, isSelected: n.id === node.id },
          };
        })
      );
    },
    [setNodes, setSelectedNodeId]
  );

  const handleUpdateNode = useCallback(
    (updatedNode: CanvasNode) => {
      setNodes((nds) =>
        nds.map((n) =>
          n.id === updatedNode.id
            ? { ...n, data: updatedNode }
            : n
        )
      );
    },
    [setNodes]
  );

  if (isPending) {
    return (
      <div className="w-full h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4" />
          <p className="text-gray-600 text-lg">Loading canvas...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="w-full h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center max-w-md">
          <p className="text-red-500 font-bold text-lg mb-2">Error</p>
          <p className="text-gray-600 text-sm">{error.message}</p>
        </div>
      </div>
    );
  }

  if (!data || data.nodes.length === 0) {
    return (
      <div className="w-full h-screen flex items-center justify-center bg-gray-50">
        <p className="text-gray-600">No nodes found</p>
      </div>
    );
  }

  const edges = data.edges ? normalizeEdges(data.edges) : [];

  return (
    <div className="flex w-full h-screen bg-white">
      <div className="flex-1">
        <ReactFlow
          nodes={nodes}
          edges={edges}
          onNodesChange={onNodesChange}
          onEdgesChange={onEdgesChange}
          onNodeClick={handleNodeClick}
          nodeTypes={nodeTypes}
          fitView
        >
          <Background />
          <Controls />
        </ReactFlow>
      </div>
      <ConfigurationPanel selectedNode={selectedNode} onUpdateNode={handleUpdateNode} />
    </div>
  );
}

export const FlowCanvas = memo(FlowCanvasComponent);
