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
import { usePricingStore } from '@/stores/pricingStore';
import { useDeploy } from '@/hooks/useDeploy';
import { useDeploymentEvents } from '@/hooks/useDeploymentEvents';
import AppNode from '@/components/nodes/AppNode';
import DatabaseNode from '@/components/nodes/DatabaseNode';
import { FlowToolbar } from './FlowToolbar';
import { ConfigurationPanel } from './ConfigurationPanel';
import { DeploymentLogsPanel } from './DeploymentLogsPanel';
import { createDeployPayload } from '@/utils/nodeUtils';
import type { FlowNode, NodeType } from '@/types/flow';

const nodeTypes: NodeTypes = {
  app: AppNode,
  database: DatabaseNode,
};

const canConnect = (sourceType: NodeType, targetType: NodeType): boolean => {
  // Only allow App → Database connections
  // App represents the application that depends on Database
  // Database provides configuration (connection string) to App
  return sourceType === 'app' && targetType === 'database';
};

function FlowCanvasComponent() {
  const { data, isPending, error } = useProjects();

  const {
    nodes: storeNodes,
    setSelectedNodeId,
    addNode,
    updateNode,
    removeNode,
    addEdge: addStoreEdge,
    removeEdge: removeStoreEdge,
    loadFromApi,
    deploymentId,
    isDeploying,
    pendingNodeIds,
    showLogs,
    setShowLogs,
    setDeploymentContext,
    updateNodeStatus,
  } = useFlowStore();

  const [flowNodes, setFlowNodes, onNodesChange] = useNodesState<Node>([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);

  const selectedNodeId = useFlowStore((state) => state.selectedNodeId);
  const selectedNodeIdRef = useRef(selectedNodeId);

  const selectedNode: FlowNode | null =
    storeNodes.find((n) => n.id === selectedNodeId) ?? null;

  const { pricing, fetch: fetchPricing } = usePricingStore();

  // Sync Zustand store to React Flow
  useEffect(() => {
    setFlowNodes((nds) =>
      nds.map((n) => {
        const sn = storeNodes.find((s) => s.id === n.id);
        if (!sn) return n;
        return {
          ...n,
          data: {
            ...n.data,
            ...sn,
            isSelected: sn.id === selectedNodeId,
          },
        };
      })
    );
  }, [storeNodes, selectedNodeId, setFlowNodes]);

  useEffect(() => {
    if (!data) return;

    loadFromApi(data.nodes, data.edges);

    setFlowNodes((currentNodes) => {
      const positionMap = new Map(currentNodes.map((n) => [n.id, n.position]));
      return data.nodes.map((node, idx) => ({
        id: node.id,
        type: node.type,
        position: positionMap.get(node.id) ?? { x: (idx % 3) * 350, y: Math.floor(idx / 3) * 300 },
        data: { ...node, isSelected: node.id === selectedNodeIdRef.current },
      }));
    });

    setEdges(
      data.edges.map((edge) => ({
        id: edge.id,
        source: edge.source,
        target: edge.target,
        animated: true,
      }))
    );
  }, [data, loadFromApi, setFlowNodes, setEdges]);

  useEffect(() => {
    if (!pricing) {
      fetchPricing();
    }
  }, [pricing, fetchPricing]);

  // Handle deployment
  const { mutate: deploy } = useDeploy({
    onDeployStarted: (deploymentId, nodeIds) => {
      const ids = nodeIds.length > 0 ? nodeIds : storeNodes.map(n => n.id);
      setDeploymentContext(deploymentId, ids);
      setShowLogs(true);
    },
  });

  // WebSocket events for real-time updates
  const { logs, isConnected } = useDeploymentEvents({
    deploymentId,
    nodeIds: pendingNodeIds,
    onLevelNodesStart: (nodeIds) => {
      nodeIds.forEach((nodeId) => {
        updateNodeStatus(nodeId, 'deploying');
      });
    },
    onNodeStatusChange: (nodeId, status, message) => {
      updateNodeStatus(nodeId, status, message);
    },
    onComplete: () => {
      // Don't auto-close logs - user can review them and close manually
    },
  });

  const handleConnect = useCallback(
    (connection: Connection) => {
      if (!connection.source || !connection.target) return;

      const sourceNode = flowNodes.find((n) => n.id === connection.source);
      const targetNode = flowNodes.find((n) => n.id === connection.target);

      if (!sourceNode || !targetNode) return;

      const sourceType = sourceNode.type as NodeType;
      const targetType = targetNode.type as NodeType;

      if (!canConnect(sourceType, targetType)) return;

      const newEdge = { ...connection, animated: true };
      setEdges((eds) => addEdge(newEdge, eds));
      addStoreEdge(connection.source, connection.target);
    },
    [flowNodes, setEdges, addStoreEdge]
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

  const handleDeploy = useCallback(() => {
    const payload = createDeployPayload(
      storeNodes,
      edges.map((edge) => ({ source: edge.source, target: edge.target }))
    );
    deploy(payload);
  }, [storeNodes, edges, deploy]);

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

      setFlowNodes((nds) => [...nds, flowNode]);
      setSelectedNodeId(newNode.id);
    },
    [addNode, setFlowNodes, setSelectedNodeId]
  );

  const handleUpdateNode = useCallback(
    (updatedNode: FlowNode) => {
      updateNode(updatedNode.id, updatedNode);
      setFlowNodes((nds) =>
        nds.map((n) =>
          n.id === updatedNode.id ? { ...n, data: { ...updatedNode, isSelected: n.data.isSelected } } : n
        )
      );
    },
    [updateNode, setFlowNodes]
  );

  const handleDeleteNode = useCallback(
    (id: string) => {
      removeNode(id);
      setFlowNodes((nds) => nds.filter((n) => n.id !== id));
      setEdges((eds) => eds.filter((e) => e.source !== id && e.target !== id));
    },
    [removeNode, setFlowNodes, setEdges]
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
        nodes={flowNodes}
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

      <FlowToolbar 
        onAddNode={handleAddNode} 
        onDeploy={handleDeploy}
        isDeploying={isDeploying}
      />

      <ConfigurationPanel
        selectedNode={selectedNode}
        onUpdateNode={handleUpdateNode}
        onDeleteNode={handleDeleteNode}
      />

      <DeploymentLogsPanel
        isOpen={showLogs}
        onClose={() => setShowLogs(false)}
        logs={logs}
        isConnected={isConnected}
        deploymentId={deploymentId}
      />
    </div>
  );
}

export const FlowCanvas = memo(FlowCanvasComponent);
