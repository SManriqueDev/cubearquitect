import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import {
  getToken,
  setToken,
  removeToken,
  getProjectId,
  setProjectId,
  removeProjectId,
} from '@/services/api';

export interface Project {
  id: number;
  name: string;
  description: string;
}

export interface SSHKey {
  id: number;
  name: string;
  key_type: string;
  fingerprint: string;
}

interface AccountState {
  isConfigured: boolean;
  token: string | null;
  projectId: number | null;
  projectName: string | null;
  sshKeys: SSHKey[];
  selectedSSHKeys: string[];
  
  // Actions
  configure: (token: string) => void;
  configureWithProject: (token: string, projectId: number, projectName: string) => void;
  clear: () => void;
  setProject: (projectId: number, projectName: string) => void;
  setSSHKeys: (keys: SSHKey[]) => void;
  toggleSSHKey: (keyName: string) => void;
  setSelectedSSHKeys: (keys: string[]) => void;
  
  // Initialization
  initialize: () => void;
}

export const useAccountStore = create<AccountState>()(
  persist(
    (set, get) => ({
      isConfigured: false,
      token: null,
      projectId: null,
      projectName: null,
      sshKeys: [],
      selectedSSHKeys: [],

      configure: (token: string) => {
        setToken(token);
        set({
          token,
          isConfigured: true,
        });
      },

      configureWithProject: (token: string, projectId: number, projectName: string) => {
        setToken(token);
        setProjectId(projectId);
        set({
          token,
          isConfigured: true,
          projectId,
          projectName,
        });
      },

      clear: () => {
        removeToken();
        removeProjectId();
        set({
          isConfigured: false,
          token: null,
          projectId: null,
          projectName: null,
          sshKeys: [],
          selectedSSHKeys: [],
        });
      },

      setProject: (projectId: number, projectName: string) => {
        setProjectId(projectId);
        set({ projectId, projectName });
      },

      setSSHKeys: (keys: SSHKey[]) => {
        set({ sshKeys: keys });
      },

      toggleSSHKey: (keyName: string) => {
        const { selectedSSHKeys } = get();
        if (selectedSSHKeys.includes(keyName)) {
          set({
            selectedSSHKeys: selectedSSHKeys.filter(k => k !== keyName),
          });
        } else {
          set({
            selectedSSHKeys: [...selectedSSHKeys, keyName],
          });
        }
      },

      setSelectedSSHKeys: (keys: string[]) => {
        set({ selectedSSHKeys: keys });
      },

      initialize: () => {
        const token = getToken();
        const projectId = getProjectId();
        set({
          token,
          projectId,
          isConfigured: token !== null,
        });
      },
    }),
    {
      name: 'account-storage',
      partialize: (state) => ({
        isConfigured: state.isConfigured,
        token: state.token,
        projectId: state.projectId,
        projectName: state.projectName,
        sshKeys: state.sshKeys,
        selectedSSHKeys: state.selectedSSHKeys,
      }),
    }
  )
);
