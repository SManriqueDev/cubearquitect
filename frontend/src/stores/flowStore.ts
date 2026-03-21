import { create } from 'zustand';

interface CanvasUiState {
  selectedNodeId: string | null;
  setSelectedNodeId: (nodeId: string | null) => void;
  clearSelection: () => void;
}

export const useFlowStore = create<CanvasUiState>((set) => ({
  selectedNodeId: null,
  setSelectedNodeId: (nodeId) => {
    set({ selectedNodeId: nodeId });
  },
  clearSelection: () => {
    set({ selectedNodeId: null });
  },
}));

