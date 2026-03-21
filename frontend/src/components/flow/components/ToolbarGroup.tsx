import { memo } from 'react';
import { cn } from '@/lib/utils';

interface ToolbarGroupProps {
  children: React.ReactNode;
  className?: string;
}

export const ToolbarGroup = memo(function ToolbarGroup({
  children,
  className,
}: ToolbarGroupProps) {
  return (
    <div
      role="group"
      className={cn('flex items-center gap-1', className)}
    >
      {children}
    </div>
  );
});
