import { cn } from '@/lib/utils';
import { Separator } from '@/components/ui/separator';

interface ToolbarSeparatorProps {
  orientation?: 'horizontal' | 'vertical';
  className?: string;
}

export function ToolbarSeparator({
  orientation = 'vertical',
  className,
}: ToolbarSeparatorProps) {
  return (
    <Separator
      orientation={orientation}
      className={cn(
        orientation === 'vertical' ? 'h-6 w-px' : 'h-px w-full',
        'bg-border/60',
        className
      )}
    />
  );
}
