import { useState } from 'react';
import { apiFetch } from '@/services/api';
import { useAccountStore, type Project, type SSHKey } from '@/stores/accountStore';

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

  return (
    <div className="min-h-screen flex items-center justify-center bg-zinc-950">
      <div className="w-full max-w-md p-8 bg-zinc-900 rounded-lg border border-zinc-800">
        <div className="text-center mb-8">
          <h1 className="text-2xl font-bold text-white mb-2">CubeArchitect</h1>
          <p className="text-zinc-400">Configure your CubePath account</p>
        </div>

        {step === 'token' && (
          <form onSubmit={handleTokenSubmit}>
            <div className="mb-4">
              <label className="block text-sm font-medium text-zinc-300 mb-2">
                API Token
              </label>
              <input
                type="password"
                value={tokenInput}
                onChange={(e) => setTokenInput(e.target.value)}
                placeholder="Enter your CubePath API token"
                className="w-full px-4 py-2 bg-zinc-800 border border-zinc-700 rounded-md text-white placeholder-zinc-500 focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              />
              <p className="mt-2 text-xs text-zinc-500">
                Get your token from the CubePath dashboard
              </p>
            </div>
            {error && <p className="text-red-400 text-sm mb-4">{error}</p>}
            <button
              type="submit"
              disabled={loading || !tokenInput.trim()}
              className="w-full py-2 px-4 bg-blue-600 hover:bg-blue-700 disabled:bg-zinc-700 disabled:cursor-not-allowed text-white rounded-md font-medium transition-colors"
            >
              {loading ? 'Validating...' : 'Connect'}
            </button>
          </form>
        )}

        {step === 'project' && (
          <form onSubmit={handleProjectSelect}>
            <div className="mb-4">
              <label className="block text-sm font-medium text-zinc-300 mb-2">
                Select Project
              </label>
              <select
                value={selectedProjectId || ''}
                onChange={(e) => {
                  const id = parseInt(e.target.value);
                  const proj = projects.find(p => p.id === id);
                  setSelectedProjectId(id);
                  setSelectedProjectName(proj?.name || '');
                }}
                className="w-full px-4 py-2 bg-zinc-800 border border-zinc-700 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
                required
              >
                <option value="">Select a project...</option>
                {projects.map((project) => (
                  <option key={project.id} value={project.id}>
                    {project.name}
                  </option>
                ))}
              </select>
            </div>
            {error && <p className="text-red-400 text-sm mb-4">{error}</p>}
            <button
              type="submit"
              disabled={loading || !selectedProjectId}
              className="w-full py-2 px-4 bg-blue-600 hover:bg-blue-700 disabled:bg-zinc-700 disabled:cursor-not-allowed text-white rounded-md font-medium transition-colors"
            >
              {loading ? 'Loading...' : 'Continue'}
            </button>
          </form>
        )}

        {step === 'ssh' && (
          <form onSubmit={handleSSHKeysSubmit}>
            <div className="mb-4">
              <label className="block text-sm font-medium text-zinc-300 mb-2">
                SSH Keys (optional)
              </label>
              {sshKeys.length === 0 ? (
                <p className="text-zinc-500 text-sm">No SSH keys found</p>
              ) : (
                <div className="space-y-2 max-h-48 overflow-y-auto">
                  {sshKeys.map((key) => (
                    <label key={key.id} className="flex items-center gap-2 text-zinc-300">
                      <input
                        type="checkbox"
                        checked={selectedSSHKeys.includes(key.name)}
                        onChange={(e) => {
                          if (e.target.checked) {
                            setSelectedSSHKeys([...selectedSSHKeys, key.name]);
                          } else {
                            setSelectedSSHKeys(selectedSSHKeys.filter(k => k !== key.name));
                          }
                        }}
                        className="rounded bg-zinc-800 border-zinc-700"
                      />
                      <span className="text-sm">{key.name}</span>
                    </label>
                  ))}
                </div>
              )}
            </div>
            <div className="flex gap-2">
              <button
                type="submit"
                className="flex-1 py-2 px-4 bg-blue-600 hover:bg-blue-700 text-white rounded-md font-medium transition-colors"
              >
                Save
              </button>
              <button
                type="button"
                onClick={handleSkipSSH}
                className="flex-1 py-2 px-4 bg-zinc-700 hover:bg-zinc-600 text-white rounded-md font-medium transition-colors"
              >
                Skip
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}
