import { FlowCanvas } from '@/components/flow/FlowCanvas';
import { Toaster } from '@/components/ui/sonner';
import { TooltipProvider } from '@/components/ui/tooltip';

export default function App() {
  return (
    <TooltipProvider delayDuration={300}>
      <FlowCanvas />
      <Toaster />
    </TooltipProvider>
  );
}