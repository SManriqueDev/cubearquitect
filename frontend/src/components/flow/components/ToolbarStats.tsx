import { memo } from 'react';
import { cn } from '@/lib/utils';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { Layers, GitBranch } from 'lucide-react';

interface NodeStatItem {
  type: string;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  count: number;
}

interface ToolbarStatsProps {
  nodes: NodeStatItem[];
  edgeCount: number;
  className?: string;
}

export const ToolbarStats = memo(function ToolbarStats({
  nodes,
  edgeCount,
  className,
}: ToolbarStatsProps) {
  const totalNodes = nodes.reduce((sum, n) => sum + n.count, 0);
  
  const hasContent = totalNodes > 0 || edgeCount > 0;

  const tooltipContent = (
    <div className="flex flex-col gap-1">
      <div className="font-medium text-xs mb-1">Canvas Stats</div>
      {nodes.map((node) => (
        <div key={node.type} className="flex items-center gap-2">
          <node.icon className="size-3" />
          <span>{node.label}: {node.count}</span>
        </div>
      ))}
      <div className="flex items-center gap-2 mt-1 pt-1 border-t border-white/20">
        <GitBranch className="size-3" />
        <span>Edges: {edgeCount}</span>
      </div>
    </div>
  );

  if (!hasContent) {
    return null;
  }

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <button
          type="button"
          className={cn(
            'inline-flex items-center gap-1.5 rounded-md',
            'h-8 px-2.5 text-xs font-medium',
            'bg-muted/50 text-muted-foreground',
            'hover:bg-muted hover:text-foreground',
            'transition-colors cursor-default',
            className
          )}
          aria-label="View canvas statistics"
        >
          <Layers className="size-3.5" />
          <span>{totalNodes}</span>
          {edgeCount > 0 && (
            <>
              <span className="text-muted-foreground/50">·</span>
              <GitBranch className="size-3.5" />
              <span>{edgeCount}</span>
            </>
          )}
        </button>
      </TooltipTrigger>
      <TooltipContent side="bottom" sideOffset={8} className="max-w-none">
        {tooltipContent}
      </TooltipContent>
    </Tooltip>
  );
});