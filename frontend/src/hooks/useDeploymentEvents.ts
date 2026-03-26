import { useCallback, useEffect, useRef, useState } from 'react';
import type { DeploymentEvent, DeploymentLogEntry, NodeStatus } from '@/types/flow';

interface UseDeploymentEventsOptions {
  deploymentId: string | null;
  nodeIds: string[];
  onNodeStatusChange?: (nodeId: string, status: NodeStatus, message?: string) => void;
  onLevelNodesStart?: (nodeIds: string[]) => void;
  onLevelStart?: (levelIdx: number) => void;
  onLevelComplete?: (levelIdx: number) => void;
  onError?: (message: string) => void;
  onComplete?: () => void;
}

interface UseDeploymentEventsReturn {
  events: DeploymentEvent[];
  logs: DeploymentLogEntry[];
  isConnected: boolean;
  reconnect: () => void;
}

const MAX_RETRIES = 3;
const RETRY_DELAY_BASE = 1000;

function mapBackendStatus(status: string): NodeStatus {
  const normalized = status.toLowerCase();
  if (normalized === 'active' || normalized === 'healthy' || normalized === 'running') {
    return 'active';
  }
  if (normalized === 'deploying' || normalized === 'creating' || normalized === 'installing') {
    return 'deploying';
  }
  if (normalized === 'pending' || normalized === 'queued') {
    return 'pending';
  }
  if (normalized === 'error' || normalized === 'failed') {
    return 'error';
  }
  return 'inactive';
}

function generateLogId(event: DeploymentEvent): string {
  return `log-${event.timestamp}-${event.type}${event.node_id ? `-${event.node_id}` : ''}`;
}

function createLogEntry(event: DeploymentEvent, nodeIds: string[]): DeploymentLogEntry | null {
  const nodeLabel = event.node_id ? nodeIds.includes(event.node_id) ? `Node ${event.node_id}` : '' : '';
  
  switch (event.type) {
    case 'connected':
      return {
        id: generateLogId(event),
        type: 'info',
        message: 'Connected to deployment stream',
        timestamp: event.timestamp,
      };
    case 'node_update':
      if (event.status) {
        const status = mapBackendStatus(event.status);
        if (status === 'active') {
          return {
            id: generateLogId(event),
            type: 'success',
            message: `${nodeLabel} deployed successfully`,
            nodeId: event.node_id,
            timestamp: event.timestamp,
          };
        } else if (status === 'deploying') {
          return {
            id: generateLogId(event),
            type: 'info',
            message: `${nodeLabel} is being deployed`,
            nodeId: event.node_id,
            timestamp: event.timestamp,
          };
        } else if (status === 'error') {
          return {
            id: generateLogId(event),
            type: 'error',
            message: `${nodeLabel} failed: ${event.message || 'Unknown error'}`,
            nodeId: event.node_id,
            timestamp: event.timestamp,
          };
        }
      }
      return null;
    case 'level_start':
      const nodeCount = event.node_ids?.length ?? 0;
      return {
        id: generateLogId(event),
        type: 'info',
        message: `Starting deployment of ${nodeCount} node${nodeCount !== 1 ? 's' : ''}`,
        timestamp: event.timestamp,
      };
    case 'level_complete':
      return {
        id: generateLogId(event),
        type: 'success',
        message: 'Level completed',
        timestamp: event.timestamp,
      };
    case 'error':
      return {
        id: generateLogId(event),
        type: 'error',
        message: event.message || 'Deployment error',
        timestamp: event.timestamp,
      };
    default:
      return null;
  }
}

export function useDeploymentEvents({
  deploymentId,
  nodeIds,
  onNodeStatusChange,
  onLevelNodesStart,
  onLevelStart,
  onLevelComplete,
  onError,
  onComplete,
}: UseDeploymentEventsOptions): UseDeploymentEventsReturn {
  const [events, setEvents] = useState<DeploymentEvent[]>([]);
  const [logs, setLogs] = useState<DeploymentLogEntry[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  
  const wsRef = useRef<WebSocket | null>(null);
  const retryCountRef = useRef(0);
  const retryTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const isCompletedRef = useRef(false);
  const isConnectingRef = useRef(false);

  const onNodeStatusChangeRef = useRef(onNodeStatusChange);
  const onLevelNodesStartRef = useRef(onLevelNodesStart);
  const onLevelStartRef = useRef(onLevelStart);
  const onLevelCompleteRef = useRef(onLevelComplete);
  const onErrorRef = useRef(onError);
  const onCompleteRef = useRef(onComplete);
  const nodeIdsRef = useRef(nodeIds);

  useEffect(() => {
    onNodeStatusChangeRef.current = onNodeStatusChange;
    onLevelNodesStartRef.current = onLevelNodesStart;
    onLevelStartRef.current = onLevelStart;
    onLevelCompleteRef.current = onLevelComplete;
    onErrorRef.current = onError;
    onCompleteRef.current = onComplete;
  }, [onNodeStatusChange, onLevelNodesStart, onLevelStart, onLevelComplete, onError, onComplete]);

  useEffect(() => {
    nodeIdsRef.current = nodeIds;
  }, [nodeIds]);

  const connect = useCallback(() => {
    if (!deploymentId || isCompletedRef.current || isConnectingRef.current) return;

    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/deployments/${deploymentId}/events`;

    isConnectingRef.current = true;

    try {
      const ws = new WebSocket(wsUrl);
      wsRef.current = ws;

      ws.onopen = () => {
        setIsConnected(true);
        retryCountRef.current = 0;
        isConnectingRef.current = false;
      };

      ws.onmessage = (event) => {
        try {
          const data: DeploymentEvent = JSON.parse(event.data);
          setEvents((prev) => [...prev, data]);

          const logEntry = createLogEntry(data, nodeIdsRef.current);
          if (logEntry) {
            setLogs((prev) => [...prev, logEntry]);
          }

          switch (data.type) {
            case 'node_update':
              if (data.node_id && data.status) {
                const status = mapBackendStatus(data.status);
                onNodeStatusChangeRef.current?.(data.node_id, status, data.message);
              }
              break;
            case 'level_start':
              if (data.node_ids && data.node_ids.length > 0) {
                onLevelNodesStartRef.current?.(data.node_ids);
              }
              if (data.level_idx !== undefined) {
                onLevelStartRef.current?.(data.level_idx);
              }
              break;
            case 'level_complete':
              if (data.level_idx !== undefined) {
                onLevelCompleteRef.current?.(data.level_idx);
              }
              break;
            case 'error':
              onErrorRef.current?.(data.message || 'Deployment failed');
              break;
          }
        } catch (err) {
          console.error('Failed to parse WebSocket message:', err);
        }
      };

      ws.onclose = () => {
        setIsConnected(false);
        wsRef.current = null;
        isConnectingRef.current = false;

        if (!isCompletedRef.current && retryCountRef.current < MAX_RETRIES) {
          const delay = RETRY_DELAY_BASE * Math.pow(2, retryCountRef.current);
          retryCountRef.current++;
          retryTimeoutRef.current = setTimeout(connect, delay);
        }
      };

      ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        isConnectingRef.current = false;
      };
    } catch (err) {
      console.error('Failed to create WebSocket:', err);
      isConnectingRef.current = false;
    }
  }, [deploymentId]);

  const reconnect = useCallback(() => {
    if (retryTimeoutRef.current) {
      clearTimeout(retryTimeoutRef.current);
    }
    retryCountRef.current = 0;
    isCompletedRef.current = false;
    isConnectingRef.current = false;
    
    if (wsRef.current) {
      wsRef.current.close();
    }
    
    connect();
  }, [connect]);

  useEffect(() => {
    if (deploymentId && nodeIds.length > 0) {
      isCompletedRef.current = false;
      setEvents([]);
      setLogs([]);
      connect();
    }

    return () => {
      if (retryTimeoutRef.current) {
        clearTimeout(retryTimeoutRef.current);
      }
      // Prevent reconnection logic in onclose from running after unmount/intentional close
      isCompletedRef.current = true;
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [deploymentId, nodeIds.length, connect]);

  useEffect(() => {
    if (nodeIds.length > 0 && events.length > 0) {
      const lastEvent = events[events.length - 1];
      
      if (lastEvent.type === 'error') {
        isCompletedRef.current = true;
        onCompleteRef.current?.();
      }
      
      const completedNodeUpdates = events.filter(
        (e) => e.type === 'node_update' && e.status && 
        ['active', 'error'].includes(mapBackendStatus(e.status).toLowerCase())
      );
      
      if (completedNodeUpdates.length === nodeIds.length) {
        isCompletedRef.current = true;
        setTimeout(() => onCompleteRef.current?.(), 1000);
      }
    }
  }, [events, nodeIds.length]);

  return {
    events,
    logs,
    isConnected,
    reconnect,
  };
}
