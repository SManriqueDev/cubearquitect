import { useState, memo } from 'react';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Server, Database, Trash2, X, Loader2 } from 'lucide-react';
import { usePricingStore } from '@/stores/pricingStore';
import { cn } from '@/lib/utils';
import type { FlowNode, AppNodeData, DatabaseNodeData } from '@/types/flow';

interface ConfigurationPanelProps {
  selectedNode: FlowNode | null;
  onUpdateNode?: (node: FlowNode) => void;
  onDeleteNode?: (id: string) => void;
}

function ConfigurationPanelComponent({
  selectedNode,
  onUpdateNode,
  onDeleteNode,
}: ConfigurationPanelProps) {
  const { pricing, isPending } = usePricingStore();

  if (!selectedNode) {
    return (
      <div className="w-80 bg-white border-l border-gray-200 p-8 flex flex-col items-center justify-center min-h-screen">
        <div className="text-center space-y-3">
          <div className="w-12 h-12 rounded-full bg-gray-100 flex items-center justify-center mx-auto">
            <X className="w-5 h-5 text-gray-400" />
          </div>
          <p className="text-gray-500 text-sm font-medium">No node selected</p>
          <p className="text-gray-400 text-xs">Click on a node to configure</p>
        </div>
      </div>
    );
  }

  if (isPending) {
    return (
      <div className="w-80 bg-white border-l border-gray-200 p-8 flex flex-col items-center justify-center min-h-screen">
        <Loader2 className="w-8 h-8 animate-spin text-gray-400" />
        <p className="mt-3 text-gray-500 text-sm">Loading options...</p>
      </div>
    );
  }

  return (
    <NodeForm
      key={selectedNode.id}
      selectedNode={selectedNode}
      pricing={pricing}
      onUpdateNode={onUpdateNode}
      onDeleteNode={onDeleteNode}
    />
  );
}

interface NodeFormProps {
  selectedNode: FlowNode;
  pricing: { plans: Array<{ plan_name: string; cpu: number; ram: number }>; locations: Array<{ location_name: string; description: string }>; templates: Array<{ template_name: string; os_name: string; version: string }> } | undefined | null;
  onUpdateNode?: (node: FlowNode) => void;
  onDeleteNode?: (id: string) => void;
}

function NodeForm({ selectedNode, pricing, onUpdateNode, onDeleteNode }: NodeFormProps) {
  const [localEdits, setLocalEdits] = useState<Record<string, unknown>>({});
  const formData = { ...selectedNode, ...localEdits };
  const isApp = selectedNode.type === 'app';

  const handleChange = (field: string, value: unknown) => {
    setLocalEdits((prev) => ({ ...prev, [field]: value }));
  };

  const handleSave = () => {
    if (onUpdateNode) {
      onUpdateNode({ ...selectedNode, ...localEdits } as FlowNode);
    }
  };

  return (
    <div className="w-80 bg-white border-l border-gray-200 overflow-y-auto max-h-screen flex flex-col">
      <div className="sticky top-0 bg-white border-b border-gray-100 p-4 space-y-3">
        <div className="flex items-center justify-between">
          <h2 className="font-bold text-lg text-gray-900 truncate">
            {selectedNode.label}
          </h2>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => onDeleteNode?.(selectedNode.id)}
            className="text-red-500 hover:text-red-600 hover:bg-red-50"
            aria-label={`Delete node ${selectedNode.label}`}
          >
            <Trash2 className="w-4 h-4" />
          </Button>
        </div>
        <Badge
          variant="outline"
          className={`gap-2 w-fit ${
            isApp
              ? 'bg-blue-50 text-blue-700 border-blue-200'
              : 'bg-purple-50 text-purple-700 border-purple-200'
          }`}
        >
          {isApp ? (
            <Server className="w-3 h-3" />
          ) : (
            <Database className="w-3 h-3" />
          )}
          {isApp ? 'App Server' : 'Database'}
        </Badge>
      </div>

      <div className="flex-1 p-4 space-y-4">
        <div className="space-y-2">
          <Label htmlFor="name" className="text-xs font-semibold text-gray-700 uppercase">
            Name
          </Label>
          <Input
            id="name"
            value={(formData as AppNodeData).name || (formData as DatabaseNodeData).name || ''}
            onChange={(e) => handleChange('name', e.target.value)}
            className="text-sm"
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="label" className="text-xs font-semibold text-gray-700 uppercase">
            Label
          </Label>
          <Input
            id="label"
            value={formData.label || ''}
            onChange={(e) => handleChange('label', e.target.value)}
            className="text-sm"
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="plan" className="text-xs font-semibold text-gray-700 uppercase">
            Plan
          </Label>
          <select
            id="plan"
            value={formData.planName || ''}
            onChange={(e) => handleChange('planName', e.target.value)}
            className="w-full h-10 px-3 text-sm border border-gray-200 rounded-md bg-white"
          >
            <option value="">Select plan</option>
            {formData.planName && !pricing?.plans.some(p => p.plan_name === formData.planName) && (
              <option value={formData.planName}>{formData.planName} (current)</option>
            )}
            {pricing?.plans.map((plan) => (
              <option key={plan.plan_name} value={plan.plan_name}>
                {plan.plan_name} ({plan.cpu} CPU, {plan.ram} RAM)
              </option>
            ))}
          </select>
        </div>

        <div className="space-y-2">
          <Label htmlFor="region" className="text-xs font-semibold text-gray-700 uppercase">
            Region
          </Label>
          <select
            id="region"
            value={formData.locationName || ''}
            onChange={(e) => handleChange('locationName', e.target.value)}
            className="w-full h-10 px-3 text-sm border border-gray-200 rounded-md bg-white"
          >
            <option value="">Select region</option>
            {formData.locationName && !pricing?.locations.some(l => l.location_name === formData.locationName) && (
              <option value={formData.locationName}>{formData.locationName} (current)</option>
            )}
            {pricing?.locations.map((loc) => (
              <option key={loc.location_name} value={loc.location_name}>
                {loc.location_name} - {loc.description}
              </option>
            ))}
          </select>
        </div>

        {isApp && (
          <div className="space-y-2">
            <Label htmlFor="template" className="text-xs font-semibold text-gray-700 uppercase">
              OS / Template
            </Label>
            <select
              id="template"
              value={(formData as AppNodeData).templateName || ''}
              onChange={(e) => handleChange('templateName', e.target.value)}
              className="w-full h-10 px-3 text-sm border border-gray-200 rounded-md bg-white"
            >
              <option value="">Select template</option>
              {(formData as AppNodeData).templateName && !pricing?.templates.some(t => t.template_name === (formData as AppNodeData).templateName) && (
                <option value={(formData as AppNodeData).templateName}>{(formData as AppNodeData).templateName} (current)</option>
              )}
              {pricing?.templates.map((tmpl) => (
                <option key={tmpl.template_name} value={tmpl.template_name}>
                  {tmpl.os_name} {tmpl.version}
                </option>
              ))}
            </select>
          </div>
        )}

        <div className="flex items-center gap-2">
          <input
            type="checkbox"
            id="ipv4"
            checked={(formData as AppNodeData).ipv4 ?? (formData as DatabaseNodeData).ipv4 ?? false}
            onChange={(e) => handleChange('ipv4', e.target.checked)}
            className="w-4 h-4 rounded border-gray-300"
          />
          <Label htmlFor="ipv4" className="text-sm font-normal text-gray-700">
            Enable IPv4
          </Label>
        </div>

        <div className="flex items-center gap-2">
          <input
            type="checkbox"
            id="backups"
            checked={(formData as AppNodeData).enableBackups ?? (formData as DatabaseNodeData).enableBackups ?? false}
            onChange={(e) => handleChange('enableBackups', e.target.checked)}
            className="w-4 h-4 rounded border-gray-300"
          />
          <Label htmlFor="backups" className="text-sm font-normal text-gray-700">
            Enable Backups
          </Label>
        </div>



        <div className="space-y-2">
          <Label className="text-xs font-semibold text-gray-700 uppercase">
            Status
          </Label>
          <div className={cn(
            "inline-flex items-center gap-2 px-3 py-2 rounded-md border",
            selectedNode.status === "active" && "bg-green-50 border-green-200",
            selectedNode.status === "inactive" && "bg-gray-50 border-gray-200",
            selectedNode.status === "error" && "bg-red-50 border-red-200"
          )}>
            <span className="relative flex h-2 w-2">
              <span className={cn(
                "absolute inline-flex h-full w-full rounded-full opacity-75 animate-pulse",
                selectedNode.status === "active" && "bg-green-500",
                selectedNode.status === "inactive" && "bg-gray-400",
                selectedNode.status === "error" && "bg-red-500"
              )} />
              <span className={cn(
                "relative inline-flex rounded-full h-2 w-2",
                selectedNode.status === "active" && "bg-green-500",
                selectedNode.status === "inactive" && "bg-gray-400",
                selectedNode.status === "error" && "bg-red-500"
              )} />
            </span>
            <span className={cn(
              "text-sm font-medium capitalize",
              selectedNode.status === "active" && "text-green-700",
              selectedNode.status === "inactive" && "text-gray-600",
              selectedNode.status === "error" && "text-red-700"
            )}>
              {selectedNode.status}
            </span>
          </div>
        </div>
      </div>

      <div className="sticky bottom-0 bg-white border-t border-gray-100 p-4">
        <Button
          onClick={handleSave}
          className="w-full"
        >
          Save Changes
        </Button>
      </div>
    </div>
  );
}

export const ConfigurationPanel = memo(ConfigurationPanelComponent);
