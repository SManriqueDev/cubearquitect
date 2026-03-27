import { create } from 'zustand';
import type {
  FlowNode,
  FlowEdge,
  AppNodeData,
  DatabaseNodeData,
  NodeStatus,
} from '@/types/flow';

interface FlowState {
  selectedNodeId: string | null;
  setSelectedNodeId: (nodeId: string | null) => void;
  clearSelection: () => void;

  nodes: FlowNode[];
  edges: FlowEdge[];
  setNodes: (nodes: FlowNode[]) => void;
  setEdges: (edges: FlowEdge[]) => void;

  addNode: (type: 'app' | 'database') => FlowNode;
  updateNode: (id: string, data: Partial<FlowNode>) => void;
  removeNode: (id: string) => void;

  addEdge: (source: string, target: string) => void;
  removeEdge: (id: string) => void;

  loadFromApi: (nodes: FlowNode[], edges: FlowEdge[]) => void;
  reset: () => void;
  getNode: (id: string) => FlowNode | undefined;

  // Deployment context
  deploymentId: string | null;
  deployingNodeIds: string[];
  pendingNodeIds: string[];
  isDeploying: boolean;
  showLogs: boolean;
  setDeploymentContext: (deploymentId: string, nodeIds: string[]) => void;
  updateNodeStatus: (
    nodeId: string,
    status: NodeStatus,
    errorMessage?: string,
  ) => void;
  setShowLogs: (show: boolean) => void;
  setPendingNodeIds: (nodeIds: string[]) => void;
  clearDeployment: () => void;
}

const generateId = () =>
  `${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;

const defaultAppNode = (projectId = 1): AppNodeData => ({
  id: generateId(),
  type: 'app',
  name: 'new-vps',
  label: 'app Node.js new-vps',
  planName: 'gp.nano',
  templateName: 'ubuntu-24',
  locationName: 'us-mia-1',
  status: 'pending',
  ipv4: true,
  enableBackups: false,
  projectId,
  isSelected: false,
});

const defaultDatabaseNode = (projectId = 1): DatabaseNodeData => ({
  id: generateId(),
  type: 'database',
  name: 'new-database',
  label: 'database PostgreSQL new-database',
  planName: 'gp.nano',
  locationName: 'us-mia-1',
  status: 'pending',
  ipv4: true,
  enableBackups: false,
  projectId,
  isSelected: false,
});

const initialState = {
  selectedNodeId: null as string | null,
  nodes: [] as FlowNode[],
  edges: [] as FlowEdge[],
  deploymentId: null as string | null,
  deployingNodeIds: [] as string[],
  pendingNodeIds: [] as string[],
  isDeploying: false,
  showLogs: false,
};

export const useFlowStore = create<FlowState>((set, get) => ({
  ...initialState,

  setSelectedNodeId: (nodeId) => set({ selectedNodeId: nodeId }),
  clearSelection: () => set({ selectedNodeId: null }),

  setNodes: (nodes) => set({ nodes }),
  setEdges: (edges) => set({ edges }),

  addNode: (type) => {
    const newNode = type === 'app' ? defaultAppNode() : defaultDatabaseNode();

    set((state) => ({
      nodes: [...state.nodes, newNode],
    }));

    return newNode;
  },

  updateNode: (id, data) => {
    set((state) => ({
      nodes: state.nodes.map((node) =>
        node.id === id ? ({ ...node, ...data } as FlowNode) : node,
      ),
    }));
  },

  removeNode: (id) => {
    set((state) => ({
      nodes: state.nodes.filter((node) => node.id !== id),
      edges: state.edges.filter(
        (edge) => edge.source !== id && edge.target !== id,
      ),
      selectedNodeId: state.selectedNodeId === id ? null : state.selectedNodeId,
    }));
  },

  addEdge: (source, target) => {
    const newEdge: FlowEdge = {
      id: `edge-${source}-${target}`,
      source,
      target,
    };

    set((state) => {
      const exists = state.edges.some(
        (e) => e.source === source && e.target === target,
      );
      if (exists) return state;

      return { edges: [...state.edges, newEdge] };
    });
  },

  removeEdge: (id) => {
    set((state) => ({
      edges: state.edges.filter((edge) => edge.id !== id),
    }));
  },

  loadFromApi: (nodes, edges) => set({ nodes, edges }),

  reset: () => set(initialState),

  getNode: (id) => {
    return get().nodes.find((node) => node.id === id);
  },

  deploymentId: null,
  deployingNodeIds: [],
  isDeploying: false,

  setDeploymentContext: (deploymentId, nodeIds) => {
    const idsToDeploy = new Set(nodeIds);
    set((state) => ({
      deploymentId,
      deployingNodeIds: nodeIds,
      pendingNodeIds: nodeIds,
      isDeploying: true,
      nodes: state.nodes.map((node) =>
        idsToDeploy.has(node.id)
          ? { ...node, status: 'pending' as NodeStatus }
          : node,
      ),
    }));
  },

  updateNodeStatus: (nodeId, status, errorMessage) => {
    set((state) => ({
      nodes: state.nodes.map((node) =>
        node.id === nodeId
          ? ({ ...node, status, errorMessage } as FlowNode)
          : node,
      ),
    }));
  },

  clearDeployment: () => {
    set({
      deploymentId: null,
      deployingNodeIds: [],
      pendingNodeIds: [],
      isDeploying: false,
      showLogs: false,
    });
  },

  setShowLogs: (show) => set({ showLogs: show }),

  setPendingNodeIds: (nodeIds) => set({ pendingNodeIds: nodeIds }),
}));
