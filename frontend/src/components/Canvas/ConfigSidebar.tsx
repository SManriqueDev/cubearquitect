import { useState, useEffect } from 'react';
import type { CanvasNode } from '@/types/canvas';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Server, Database, X, CheckCircle, AlertCircle, Clock } from 'lucide-react';

interface ConfigSidebarProps {
  selectedNode: CanvasNode | null;
  onUpdateNode?: (node: CanvasNode) => void;
}

const statusConfig = {
  active: {
    icon: CheckCircle,
    color: 'text-green-600',
    bgColor: 'bg-green-50',
    borderColor: 'border-green-200',
    label: 'Active',
  },
  inactive: {
    icon: Clock,
    color: 'text-gray-500',
    bgColor: 'bg-gray-50',
    borderColor: 'border-gray-200',
    label: 'Inactive',
  },
  error: {
    icon: AlertCircle,
    color: 'text-red-600',
    bgColor: 'bg-red-50',
    borderColor: 'border-red-200',
    label: 'Error',
  },
};

export function ConfigSidebar({ selectedNode }: ConfigSidebarProps) {
  const [editedNode, setEditedNode] = useState<CanvasNode | null>(selectedNode);

  useEffect(() => {
    setEditedNode(selectedNode);
  }, [selectedNode]);

  if (!selectedNode) {
    return (
      <div className="w-80 bg-white border-l border-gray-200 p-8 flex flex-col items-center justify-center min-h-screen">
        <div className="text-center space-y-3">
          <div className="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center mx-auto">
            <X className="w-5 h-5 text-gray-400" />
          </div>
          <p className="text-gray-500 text-sm font-medium">No node selected</p>
          <p className="text-gray-400 text-xs">Click on a node in the canvas to view details</p>
        </div>
      </div>
    );
  }

  return (
    <div className="w-80 bg-white border-l border-gray-200 overflow-y-auto max-h-screen flex flex-col">
      {/* Header */}
      <div className="sticky top-0 bg-white border-b border-gray-100 p-6 space-y-3">
        <h2 className="font-bold text-lg text-gray-900">{selectedNode.label}</h2>
        <Badge className="flex items-center gap-2 w-fit bg-blue-50 text-blue-700 border-blue-200 hover:bg-blue-100">
          {selectedNode.type === 'app' ? (
            <>
              <Server className="w-3 h-3" />
              App Server
            </>
          ) : (
            <>
              <Database className="w-3 h-3" />
              Database
            </>
          )}
        </Badge>
      </div>

      {/* Content */}
      <div className="flex-1 p-6 space-y-5">
        <div className="space-y-4">
          {/* IP Address */}
          <div className="space-y-2">
            <Label htmlFor="ip" className="text-xs font-semibold text-gray-700 uppercase tracking-wide">
              IP Address
            </Label>
            <Input
              id="ip"
              value={editedNode?.ip || ''}
              className="text-sm bg-gray-50 border-gray-200 font-mono cursor-default"
              disabled
            />
          </div>

          {/* Plan */}
          <div className="space-y-2">
            <Label htmlFor="plan" className="text-xs font-semibold text-gray-700 uppercase tracking-wide">
              Plan
            </Label>
            <Input
              id="plan"
              value={editedNode?.plan || ''}
              className="text-sm bg-gray-50 border-gray-200 cursor-default"
              disabled
            />
          </div>

          {/* Status */}
          <div className="space-y-2">
            <Label className="text-xs font-semibold text-gray-700 uppercase tracking-wide">
              Status
            </Label>
            {(() => {
              const status = editedNode?.status || 'inactive';
              const config = statusConfig[status as keyof typeof statusConfig] || statusConfig.inactive;
              const StatusIcon = config.icon;
              return (
                <div className={`p-3 rounded-md border ${config.bgColor} ${config.borderColor} flex items-center gap-2`}>
                  <StatusIcon className={`w-5 h-5 flex-shrink-0 ${config.color}`} />
                  <span className={`text-sm font-semibold ${config.color}`}>{config.label}</span>
                </div>
              );
            })()}
          </div>

          {/* Region */}
          <div className="space-y-2">
            <Label htmlFor="region" className="text-xs font-semibold text-gray-700 uppercase tracking-wide">
              Region
            </Label>
            <Input
              id="region"
              value={editedNode?.region || ''}
              className="text-sm bg-gray-50 border-gray-200 cursor-default"
              disabled
            />
          </div>
        </div>
      </div>
    </div>
  );
}
