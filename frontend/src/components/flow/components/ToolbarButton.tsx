import { memo } from 'react';
import { cn } from '@/lib/utils';
import { TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { Tooltip } from '@/components/ui/tooltip';

interface ToolbarButtonProps {
  icon: React.ComponentType<{ className?: string }>;
  label: string;
  shortcut?: string;
  onClick?: () => void;
  disabled?: boolean;
  loading?: boolean;
  variant?: 'default' | 'outline' | 'deploy' | 'ghost';
  className?: string;
}

export const ToolbarButton = memo(function ToolbarButton({
  icon: Icon,
  label,
  shortcut,
  onClick,
  disabled,
  loading,
  variant = 'outline',
  className,
}: ToolbarButtonProps) {
  const tooltipContent = shortcut ? `${label} (${shortcut})` : label;

  const variantStyles = {
    default: 'hover:bg-muted',
    outline: 'hover:bg-muted dark:hover:bg-muted/50',
    deploy: 'bg-emerald-600 hover:bg-emerald-700 text-white shadow-sm',
    ghost: 'hover:bg-muted/50',
  };

  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <button
          type="button"
          onClick={onClick}
          disabled={disabled || loading}
          className={cn(
            'inline-flex items-center justify-center gap-1.5 rounded-md',
            'h-8 px-2.5 text-sm font-medium',
            'transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2',
            'disabled:pointer-events-none disabled:opacity-50',
            variant === 'deploy' 
              ? variantStyles.deploy
              : variant === 'ghost'
                ? variantStyles.ghost
                : 'bg-card border border-border shadow-xs',
            variantStyles[variant],
            className
          )}
          aria-label={label}
        >
          {loading ? (
            <span className="animate-spin size-4" />
          ) : (
            <Icon className="size-4" />
          )}
          <span className="sr-only sm:not-sr-only">{label}</span>
        </button>
      </TooltipTrigger>
      <TooltipContent side="bottom" sideOffset={8}>
        <p>{tooltipContent}</p>
      </TooltipContent>
    </Tooltip>
  );
});
