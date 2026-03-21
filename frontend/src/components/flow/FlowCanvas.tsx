import { memo, useCallback, useEffect, useRef } from 'react';
import {
  ReactFlow,
  Background,
  Controls,
  MiniMap,
  useNodesState,
  useEdgesState,
  addEdge,
  type Connection,
  type Node,
  type Edge,
  type NodeTypes,
  type NodeChange,
  type EdgeChange,
  type OnNodesChange,
  type OnEdgesChange,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';

import { useProjects } from '@/hooks/useProjects';
import { useFlowStore } from '@/stores/flowStore';
import { useNodeCount, useEdgeCount } from '@/hooks/useFlowStore';
import AppNode from '@/components/nodes/AppNode';
import DatabaseNode from '@/components/nodes/DatabaseNode';
import { FlowToolbar } from './FlowToolbar';
import { ConfigurationPanel } from './ConfigurationPanel';
import type { FlowNode, NodeType } from '@/types/flow';

const nodeTypes: NodeTypes = {
  app: AppNode,
  database: DatabaseNode,
};

const canConnect = (sourceType: NodeType, targetType: NodeType): boolean => {
  return !(sourceType === 'database' && targetType === 'database');
};

function FlowCanvasComponent() {
  const { data, isPending, error } = useProjects();
  const isInitializedRef = useRef(false);

  const {
    nodes: storeNodes,
    setSelectedNodeId,
    addNode,
    updateNode,
    removeNode,
    addEdge: addStoreEdge,
    removeEdge: removeStoreEdge,
    loadFromApi,
  } = useFlowStore();

  const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);

  const selectedNodeId = useFlowStore((state) => state.selectedNodeId);
  const nodeCount = useNodeCount();
  const edgeCount = useEdgeCount();

  const selectedNode: FlowNode | null =
    storeNodes.find((n) => n.id === selectedNodeId) ?? null;

  // Sync selection state to React Flow nodes
  useEffect(() => {
    setNodes((nds) =>
      nds.map((n) => ({
        ...n,
        data: {
          ...n.data,
          isSelected: n.id === selectedNodeId,
        },
      }))
    );
  }, [selectedNodeId, setNodes]);

  useEffect(() => {
    if (data && !isInitializedRef.current) {
      isInitializedRef.current = true;
      loadFromApi(data.nodes, data.edges);

      const initialNodes: Node[] = data.nodes.map((node, idx) => ({
        id: node.id,
        type: node.type,
        position: { x: (idx % 3) * 350, y: Math.floor(idx / 3) * 300 },
        data: { ...node, isSelected: false },
      }));

      const initialEdges: Edge[] = data.edges.map((edge) => ({
        id: edge.id,
        source: edge.source,
        target: edge.target,
        animated: true,
      }));

      setNodes(initialNodes);
      setEdges(initialEdges);
    }
  }, [data, loadFromApi, setNodes, setEdges]);

  const handleConnect = useCallback(
    (connection: Connection) => {
      if (!connection.source || !connection.target) return;

      const sourceNode = nodes.find((n) => n.id === connection.source);
      const targetNode = nodes.find((n) => n.id === connection.target);

      if (!sourceNode || !targetNode) return;

      const sourceType = sourceNode.type as NodeType;
      const targetType = targetNode.type as NodeType;

      if (!canConnect(sourceType, targetType)) return;

      const newEdge = { ...connection, animated: true };
      setEdges((eds) => addEdge(newEdge, eds));
      addStoreEdge(connection.source, connection.target);
    },
    [nodes, setEdges, addStoreEdge]
  );

  const handleNodesChange: OnNodesChange = useCallback(
    (changes: NodeChange[]) => {
      onNodesChange(changes);

      const removeChanges = changes.filter((c) => c.type === 'remove');
      if (removeChanges.length > 0) {
        removeChanges.forEach((change) => {
          if (change.type === 'remove') {
            removeNode(change.id);
          }
        });
      }
    },
    [onNodesChange, removeNode]
  );

  const handleEdgesChange: OnEdgesChange = useCallback(
    (changes: EdgeChange[]) => {
      onEdgesChange(changes);

      changes.forEach((change) => {
        if (change.type === 'remove') {
          removeStoreEdge(change.id);
        }
      });
    },
    [onEdgesChange, removeStoreEdge]
  );

  const handleNodeClick = useCallback(
    (_: React.MouseEvent, node: Node) => {
      setSelectedNodeId(node.id);
    },
    [setSelectedNodeId]
  );

  const handlePaneClick = useCallback(() => {
    setSelectedNodeId(null);
  }, [setSelectedNodeId]);

  const handleAddNode = useCallback(
    (type: 'app' | 'database') => {
      const position = { x: 200 + Math.random() * 200, y: 200 + Math.random() * 200 };
      const newNode = addNode(type, position);

      const flowNode: Node = {
        id: newNode.id,
        type: newNode.type,
        position,
        data: { ...newNode, isSelected: true },
      };

      setNodes((nds) => [...nds, flowNode]);
      setSelectedNodeId(newNode.id);
    },
    [addNode, setNodes, setSelectedNodeId]
  );

  const handleUpdateNode = useCallback(
    (updatedNode: FlowNode) => {
      updateNode(updatedNode.id, updatedNode);
      setNodes((nds) =>
        nds.map((n) =>
          n.id === updatedNode.id ? { ...n, data: { ...updatedNode, isSelected: n.data.isSelected } } : n
        )
      );
    },
    [updateNode, setNodes]
  );

  const handleDeleteNode = useCallback(
    (id: string) => {
      removeNode(id);
      setNodes((nds) => nds.filter((n) => n.id !== id));
      setEdges((eds) => eds.filter((e) => e.source !== id && e.target !== id));
    },
    [removeNode, setNodes, setEdges]
  );

  if (isPending) {
    return (
      <div className="w-full h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500 mx-auto mb-4" />
          <p className="text-gray-600 text-lg">Loading infrastructure...</p>
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

  return (
    <div className="flex w-full h-screen bg-gray-50">
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={handleNodesChange}
        onEdgesChange={handleEdgesChange}
        onConnect={handleConnect}
        onNodeClick={handleNodeClick}
        onPaneClick={handlePaneClick}
        nodeTypes={nodeTypes}
        fitView
      >
        <Background />
        <Controls />
        <MiniMap nodeStrokeWidth={3} zoomable pannable />
      </ReactFlow>

      <FlowToolbar onAddNode={handleAddNode} nodeCount={nodeCount} edgeCount={edgeCount} />

      <ConfigurationPanel
        selectedNode={selectedNode}
        onUpdateNode={handleUpdateNode}
        onDeleteNode={handleDeleteNode}
      />
    </div>
  );
}

export const FlowCanvas = memo(FlowCanvasComponent);
