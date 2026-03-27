import { useEffect, useRef } from 'react';
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetDescription } from '@/components/ui/sheet';
import type { DeploymentLogEntry } from '@/types/flow';
import { CheckCircle, XCircle, AlertCircle, Info, Terminal } from 'lucide-react';

interface DeploymentLogsPanelProps {
  isOpen: boolean;
  onClose: () => void;
  logs: DeploymentLogEntry[];
  isConnected: boolean;
  deploymentId: string | null;
}

function getLogIcon(type: DeploymentLogEntry['type']) {
  switch (type) {
    case 'success':
      return <CheckCircle className="h-4 w-4 text-green-500" />;
    case 'error':
      return <XCircle className="h-4 w-4 text-red-500" />;
    case 'warning':
      return <AlertCircle className="h-4 w-4 text-yellow-500" />;
    default:
      return <Info className="h-4 w-4 text-blue-500" />;
  }
}

function formatTimestamp(timestamp: number): string {
  if (!timestamp || timestamp <= 0) {
    return new Date().toLocaleTimeString('en-US', {
      hour12: false,
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    });
  }
  const date = new Date(timestamp);
  return date.toLocaleTimeString('en-US', {
    hour12: false,
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  });
}

export function DeploymentLogsPanel({
  isOpen,
  onClose,
  logs,
  isConnected,
  deploymentId,
}: DeploymentLogsPanelProps) {
  const logsEndRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    logsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [logs]);

  return (
    <Sheet open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <SheetContent side="right" className="w-[400px] sm:max-w-[400px] flex flex-col">
        <SheetHeader className="pb-4 border-b">
          <div className="flex items-center gap-2">
            <Terminal className="h-5 w-5" />
            <SheetTitle>Deployment Logs</SheetTitle>
          </div>
          <SheetDescription className="flex items-center justify-between">
            <span className="truncate">
              {deploymentId ? `ID: ${deploymentId.slice(0, 8)}...` : 'No active deployment'}
            </span>
            <span
              className={`flex items-center gap-1.5 text-xs ${
                isConnected ? 'text-green-600' : 'text-gray-500'
              }`}
            >
              <span
                className={`h-2 w-2 rounded-full ${
                  isConnected ? 'bg-green-500 animate-pulse' : 'bg-gray-400'
                }`}
              />
              {isConnected ? 'Connected' : 'Disconnected'}
            </span>
          </SheetDescription>
        </SheetHeader>

        <div className="flex-1 overflow-y-auto py-4 px-2 space-y-2">
          {logs.length === 0 ? (
            <div className="text-center text-muted-foreground py-8">
              <Terminal className="h-8 w-8 mx-auto mb-2 opacity-50" />
              <p className="text-sm">Waiting for events...</p>
            </div>
          ) : (
            logs.map((log) => (
              <div
                key={log.id}
                className={`flex items-start gap-3 p-3 rounded-lg border ${
                  log.type === 'error'
                    ? 'bg-red-50 border-red-200'
                    : log.type === 'success'
                    ? 'bg-green-50 border-green-200'
                    : log.type === 'warning'
                    ? 'bg-yellow-50 border-yellow-200'
                    : 'bg-muted/50 border-border'
                }`}
              >
                <div className="flex-shrink-0 mt-0.5">{getLogIcon(log.type)}</div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm break-words">{log.message}</p>
                  <p className="text-xs text-muted-foreground mt-1">
                    {formatTimestamp(log.timestamp)}
                  </p>
                </div>
              </div>
            ))
          )}
          <div ref={logsEndRef} />
        </div>
      </SheetContent>
    </Sheet>
  );
}
