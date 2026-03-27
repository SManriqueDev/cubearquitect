import { useState } from 'react';
import { apiFetch } from '@/services/api';
import { useAccountStore, type Project, type SSHKey } from '@/stores/accountStore';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Checkbox } from '@/components/ui/checkbox';
import { Loader2, Key, Server, ArrowRight, Check } from 'lucide-react';

export function AccountSetup() {
  const [step, setStep] = useState<'token' | 'project' | 'ssh'>('token');
  const [tokenInput, setTokenInput] = useState('');
  const [projects, setProjects] = useState<Project[]>([]);
  const [sshKeys, setSSHKeys] = useState<SSHKey[]>([]);
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null);
  const [selectedProjectName, setSelectedProjectName] = useState<string>('');
  const [selectedSSHKeys, setSelectedSSHKeys] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const { configureWithProject, setSSHKeys: saveSSHKeys, setSelectedSSHKeys: saveSelectedSSHKeys } = useAccountStore();

  const handleTokenSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!tokenInput.trim()) return;

    setLoading(true);
    setError(null);

    try {
      const data = await apiFetch<{ project: Project; networks: unknown[]; vps: unknown[] }[]>('/api/projects', {
        authToken: tokenInput,
      });
      
      setProjects(data.map(p => p.project));
      setStep('project');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to validate token');
    } finally {
      setLoading(false);
    }
  };

  const handleProjectSelect = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedProjectId) return;

    setLoading(true);
    setError(null);

    try {
      const keys = await apiFetch<SSHKey[]>('/api/ssh-keys', {
        authToken: tokenInput,
      });
      setSSHKeys(keys);
      setStep('ssh');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch SSH keys');
    } finally {
      setLoading(false);
    }
  };

  const handleSSHKeysSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    configureWithProject(tokenInput, selectedProjectId!, selectedProjectName);
    saveSSHKeys(sshKeys);
    saveSelectedSSHKeys(selectedSSHKeys);
    
    window.location.reload();
  };

  const handleSkipSSH = () => {
    configureWithProject(tokenInput, selectedProjectId!, selectedProjectName);
    saveSSHKeys([]);
    saveSelectedSSHKeys([]);
    
    window.location.reload();
  };

  const toggleSSHKey = (keyName: string) => {
    if (selectedSSHKeys.includes(keyName)) {
      setSelectedSSHKeys(selectedSSHKeys.filter(k => k !== keyName));
    } else {
      setSelectedSSHKeys([...selectedSSHKeys, keyName]);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center space-y-2">
          <div className="mx-auto w-12 h-12 rounded-full bg-primary/10 flex items-center justify-center">
            <Server className="w-6 h-6 text-primary" />
          </div>
          <CardTitle className="text-2xl">CubeArchitect</CardTitle>
          <CardDescription>Connect your CubePath account to get started</CardDescription>
        </CardHeader>
        
        <CardContent>
          {error && (
            <div className="mb-6 p-3 rounded-md bg-destructive/10 border border-destructive/20 text-destructive text-sm">
              {error}
            </div>
          )}

          {step === 'token' && (
            <form onSubmit={handleTokenSubmit} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="api-token">API Token</Label>
                <div className="relative">
                  <Key className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                  <Input
                    id="api-token"
                    type="password"
                    placeholder="Enter your CubePath API token"
                    value={tokenInput}
                    onChange={(e) => setTokenInput(e.target.value)}
                    className="pl-10"
                    autoComplete="off"
                    required
                  />
                </div>
                <p className="text-xs text-muted-foreground">
                  Get your token from the CubePath dashboard
                </p>
              </div>
              
              <Button type="submit" className="w-full" disabled={loading || !tokenInput.trim()}>
                {loading ? (
                  <>
                    <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                    Validating…
                  </>
                ) : (
                  <>
                    Connect
                    <ArrowRight className="w-4 h-4 ml-2" />
                  </>
                )}
              </Button>
            </form>
          )}

          {step === 'project' && (
            <form onSubmit={handleProjectSelect} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="project-select">Select Project</Label>
                <select
                  id="project-select"
                  value={selectedProjectId || ''}
                  onChange={(e) => {
                    const id = parseInt(e.target.value);
                    const proj = projects.find(p => p.id === id);
                    setSelectedProjectId(id);
                    setSelectedProjectName(proj?.name || '');
                  }}
                  className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm transition-colors focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50"
                  required
                >
                  <option value="">Select a project…</option>
                  {projects.map((project) => (
                    <option key={project.id} value={project.id}>
                      {project.name}
                    </option>
                  ))}
                </select>
              </div>
              
              <Button type="submit" className="w-full" disabled={loading || !selectedProjectId}>
                {loading ? (
                  <>
                    <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                    Loading…
                  </>
                ) : (
                  <>
                    Continue
                    <ArrowRight className="w-4 h-4 ml-2" />
                  </>
                )}
              </Button>
            </form>
          )}

          {step === 'ssh' && (
            <div className="space-y-4">
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <Label>SSH Keys</Label>
                  <span className="text-xs text-muted-foreground">Optional</span>
                </div>
                
                {sshKeys.length === 0 ? (
                  <div className="text-center py-6 text-muted-foreground text-sm border border-dashed rounded-md">
                    No SSH keys found in your account
                  </div>
                ) : (
                  <div className="space-y-2 max-h-64 overflow-y-auto border rounded-md p-2">
                    {sshKeys.map((key) => (
                      <label
                        key={key.id}
                        className="flex items-center gap-3 p-3 rounded-md border transition-colors hover:bg-muted cursor-pointer"
                        style={{
                          backgroundColor: selectedSSHKeys.includes(key.name) ? 'var(--accent)' : undefined,
                          borderColor: selectedSSHKeys.includes(key.name) ? 'var(--primary)' : undefined,
                        }}
                      >
                        <Checkbox
                          checked={selectedSSHKeys.includes(key.name)}
                          onCheckedChange={() => toggleSSHKey(key.name)}
                        />
                        <div className="flex-1 min-w-0">
                          <div className="font-medium text-sm truncate">{key.name}</div>
                          <div className="text-xs text-muted-foreground truncate">
                            {key.fingerprint}
                          </div>
                        </div>
                        {selectedSSHKeys.includes(key.name) && (
                          <Check className="w-4 h-4 text-primary shrink-0" />
                        )}
                      </label>
                    ))}
                  </div>
                )}
              </div>
              
              <div className="flex gap-2">
                <Button type="button" onClick={handleSkipSSH} variant="outline" className="flex-1">
                  Skip
                </Button>
                <Button type="submit" onClick={handleSSHKeysSubmit} className="flex-1">
                  <Check className="w-4 h-4 mr-2" />
                  Save
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
