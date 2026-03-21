import { memo } from 'react';
import { cn } from '@/lib/utils';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Plus, Server, Database } from 'lucide-react';

export interface NodeTypeOption {
  id: string;
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  onClick: () => void;
  disabled?: boolean;
}

interface ToolbarAddButtonProps {
  onAddNode?: (type: 'app' | 'database') => void;
  options?: NodeTypeOption[];
  className?: string;
}

const defaultOptions: NodeTypeOption[] = [
  { id: 'app', label: 'Add App', icon: Server, onClick: () => {} },
  { id: 'database', label: 'Add Database', icon: Database, onClick: () => {} },
];

export const ToolbarAddButton = memo(function ToolbarAddButton({
  onAddNode,
  options: customOptions,
  className,
}: ToolbarAddButtonProps) {
  const options = customOptions || defaultOptions;

  const handleSelect = (option: NodeTypeOption) => {
    if (option.onClick) {
      option.onClick();
    }
    if (option.id === 'app') onAddNode?.('app');
    if (option.id === 'database') onAddNode?.('database');
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button
          type="button"
          className={cn(
            'inline-flex items-center justify-center gap-1.5 rounded-md',
            'h-8 px-2.5 text-sm font-medium',
            'bg-card border border-border shadow-xs',
            'hover:bg-muted dark:hover:bg-muted/50',
            'transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
            'focus-visible:ring-offset-0',
            className
          )}
          aria-label="Add node"
        >
          <Plus className="size-4" />
          <span className="sr-only sm:not-sr-only">Add</span>
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="start" sideOffset={8} className="w-48">
        <div className="px-2 py-1.5 text-xs font-medium text-muted-foreground">
          Node Types
        </div>
        {options.map((option) => (
          <DropdownMenuItem
            key={option.id}
            onSelect={() => handleSelect(option)}
            disabled={option.disabled}
            className="gap-2"
          >
            <option.icon className="size-4" />
            {option.label}
          </DropdownMenuItem>
        ))}
        <DropdownMenuSeparator />
        <DropdownMenuItem disabled className="text-muted-foreground text-xs">
          More nodes coming soon
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
});