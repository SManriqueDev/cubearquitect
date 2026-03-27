import { useEffect } from 'react';
import { FlowCanvas } from '@/components/flow/FlowCanvas';
import { AccountSetup } from '@/components/AccountSetup';
import { Toaster } from '@/components/ui/sonner';
import { TooltipProvider } from '@/components/ui/tooltip';
import { useAccountStore } from '@/stores/accountStore';

export default function App() {
  const { isConfigured, initialize } = useAccountStore();

  useEffect(() => {
    initialize();
  }, []);

  if (!isConfigured) {
    return (
      <>
        <AccountSetup />
        <Toaster />
      </>
    );
  }

  return (
    <TooltipProvider delayDuration={300}>
      <FlowCanvas />
      <Toaster />
    </TooltipProvider>
  );
}