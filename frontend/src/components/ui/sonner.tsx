import { Toaster as SonnerToasterComponent } from "sonner";

export function SonnerToaster() {
  return (
    <SonnerToasterComponent
      position="bottom-right"
      toastOptions={{
        unstyled: true,
        classNames: {
          toast: "flex w-full p-4 rounded-lg border bg-background text-foreground shadow-lg",
          title: "text-sm font-semibold",
          description: "text-sm opacity-90",
          actionButton: "bg-primary text-primary-foreground text-sm font-medium px-3 py-1 rounded-md",
          cancelButton: "bg-muted text-muted-foreground text-sm font-medium px-3 py-1 rounded-md",
          success: "border-green-500/20 bg-green-500/10",
          error: "border-red-500/20 bg-red-500/10",
          warning: "border-yellow-500/20 bg-yellow-500/10",
        },
      }}
    />
  );
}
