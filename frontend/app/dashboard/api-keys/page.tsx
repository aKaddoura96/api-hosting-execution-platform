'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/lib/auth-context';
import { api } from '@/lib/api';

interface APIKey {
  id: string;
  user_id: string;
  api_id?: string;
  key: string;
  name: string;
  is_active: boolean;
  expires_at?: string;
  created_at: string;
}

export default function APIKeysPage() {
  const router = useRouter();
  const { user, logout, loading: authLoading } = useAuth();
  const [apiKeys, setAPIKeys] = useState<APIKey[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [copiedKey, setCopiedKey] = useState<string | null>(null);

  useEffect(() => {
    if (!authLoading && !user) {
      router.push('/login');
      return;
    }

    if (user) {
      loadAPIKeys();
    }
  }, [user, authLoading, router]);

  const loadAPIKeys = async () => {
    try {
      const data: any = await api.getAPIKeys();
      setAPIKeys(data || []);
    } catch (error) {
      console.error('Failed to load API keys:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    logout();
    router.push('/');
  };

  const copyToClipboard = (key: string) => {
    navigator.clipboard.writeText(key);
    setCopiedKey(key);
    setTimeout(() => setCopiedKey(null), 2000);
  };

  const maskKey = (key: string) => {
    if (key.length < 16) return key;
    return key.substring(0, 12) + '••••••••••••••' + key.substring(key.length - 8);
  };

  if (authLoading || loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Navigation */}
      <nav className="bg-white border-b px-4 sm:px-6 py-4">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <Link href="/dashboard" className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-sm">AP</span>
            </div>
            <span className="text-lg sm:text-2xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
              API Keys
            </span>
          </Link>
          <div className="flex items-center space-x-2 sm:space-x-4">
            <Link href="/dashboard" className="text-gray-600 hover:text-gray-900 font-medium">
              Dashboard
            </Link>
            <span className="hidden md:inline text-gray-600 text-sm">Hi, {user?.name}</span>
            <button
              onClick={handleLogout}
              className="px-3 sm:px-4 py-2 text-sm sm:text-base text-red-600 hover:text-red-700 font-medium"
            >
              Logout
            </button>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 py-6 sm:py-8">
        {/* Header */}
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6 sm:mb-8">
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold">API Keys</h1>
            <p className="text-gray-600 mt-2">Manage authentication keys for your APIs</p>
          </div>
          <button
            onClick={() => setShowCreateModal(true)}
            className="w-full sm:w-auto px-4 py-2 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg hover:from-blue-700 hover:to-indigo-700 font-medium transition-all shadow-md hover:shadow-lg"
          >
            + Generate New Key
          </button>
        </div>

        {/* API Keys List */}
        {apiKeys.length === 0 ? (
          <div className="text-center py-12 bg-white rounded-lg border">
            <div className="text-5xl mb-4">🔑</div>
            <h3 className="text-xl font-semibold mb-2">No API Keys Yet</h3>
            <p className="text-gray-600 mb-4">Generate your first API key to start making authenticated requests</p>
            <button
              onClick={() => setShowCreateModal(true)}
              className="px-6 py-3 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg hover:from-blue-700 hover:to-indigo-700 font-medium shadow-md"
            >
              Generate Your First Key
            </button>
          </div>
        ) : (
          <div className="space-y-4">
            {apiKeys.map((apiKey) => (
              <div
                key={apiKey.id}
                className="bg-white p-4 sm:p-6 rounded-lg border-2 border-gray-200 hover:shadow-lg transition-shadow"
              >
                <div className="flex flex-col sm:flex-row justify-between items-start gap-4">
                  <div className="flex-1 min-w-0 w-full">
                    <div className="flex items-center gap-3 mb-3">
                      <h3 className="text-lg font-semibold">{apiKey.name}</h3>
                      <span
                        className={`px-3 py-1 rounded-full text-xs font-medium ${
                          apiKey.is_active
                            ? 'bg-green-100 text-green-800'
                            : 'bg-gray-100 text-gray-800'
                        }`}
                      >
                        {apiKey.is_active ? 'Active' : 'Inactive'}
                      </span>
                    </div>

                    <div className="bg-gray-50 p-3 rounded-lg mb-3 flex items-center gap-2">
                      <code className="flex-1 text-xs sm:text-sm text-gray-700 break-all font-mono">
                        {maskKey(apiKey.key)}
                      </code>
                      <button
                        onClick={() => copyToClipboard(apiKey.key)}
                        className="px-3 py-1.5 bg-blue-600 text-white rounded hover:bg-blue-700 text-xs font-medium whitespace-nowrap transition-colors"
                        title="Copy full key"
                      >
                        {copiedKey === apiKey.key ? '✓ Copied!' : 'Copy'}
                      </button>
                    </div>

                    <div className="text-sm text-gray-600">
                      <p>Created: {new Date(apiKey.created_at).toLocaleDateString()}</p>
                      {apiKey.expires_at && (
                        <p>Expires: {new Date(apiKey.expires_at).toLocaleDateString()}</p>
                      )}
                    </div>
                  </div>

                  {apiKey.is_active && (
                    <button
                      onClick={async () => {
                        if (confirm('Deactivate this API key? This action cannot be undone.')) {
                          try {
                            await api.deactivateAPIKey(apiKey.id);
                            await loadAPIKeys();
                          } catch (err: any) {
                            alert('Failed to deactivate key: ' + err.message);
                          }
                        }
                      }}
                      className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 font-medium text-sm transition-colors"
                    >
                      Revoke
                    </button>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Info Box */}
        <div className="mt-8 bg-blue-50 border border-blue-200 rounded-lg p-4 sm:p-6">
          <h3 className="font-semibold text-blue-900 mb-2">🔐 How to Use API Keys</h3>
          <div className="text-sm text-blue-800 space-y-2">
            <p>Include your API key in requests using one of these methods:</p>
            <div className="bg-white p-3 rounded border border-blue-200 font-mono text-xs sm:text-sm mt-2">
              <p className="mb-2">Method 1: X-API-Key header</p>
              <code className="text-gray-700">curl -H "X-API-Key: your_key_here" ...</code>
            </div>
            <div className="bg-white p-3 rounded border border-blue-200 font-mono text-xs sm:text-sm">
              <p className="mb-2">Method 2: Authorization header</p>
              <code className="text-gray-700">curl -H "Authorization: Bearer your_key_here" ...</code>
            </div>
            <p className="mt-4 font-medium">⚠️ Keep your API keys secure! Never share them publicly or commit them to version control.</p>
          </div>
        </div>
      </div>

      {/* Create API Key Modal */}
      {showCreateModal && (
        <CreateAPIKeyModal
          onClose={() => setShowCreateModal(false)}
          onSuccess={() => {
            setShowCreateModal(false);
            loadAPIKeys();
          }}
        />
      )}
    </div>
  );
}

function CreateAPIKeyModal({ onClose, onSuccess }: { onClose: () => void; onSuccess: () => void }) {
  const [name, setName] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [createdKey, setCreatedKey] = useState<APIKey | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      const result: any = await api.createAPIKey({ name });
      setCreatedKey(result);
    } catch (err: any) {
      setError(err.message);
      setLoading(false);
    }
  };

  const copyKey = () => {
    if (createdKey) {
      navigator.clipboard.writeText(createdKey.key);
      alert('API key copied to clipboard!');
    }
  };

  if (createdKey) {
    return (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
        <div className="bg-white rounded-lg p-6 sm:p-8 max-w-2xl w-full">
          <h2 className="text-2xl font-bold mb-4 text-green-600">✅ API Key Created!</h2>
          
          <div className="bg-yellow-50 border-2 border-yellow-400 rounded-lg p-4 mb-6">
            <p className="font-semibold text-yellow-900 mb-2">⚠️ Important: Save Your Key Now!</p>
            <p className="text-sm text-yellow-800">
              This is the only time you'll see the full key. Copy it now and store it securely.
            </p>
          </div>

          <div className="bg-gray-50 p-4 rounded-lg mb-4">
            <label className="block text-sm font-medium mb-2">Your API Key:</label>
            <code className="block bg-white p-3 rounded border text-sm break-all font-mono">
              {createdKey.key}
            </code>
          </div>

          <div className="flex gap-3">
            <button
              onClick={copyKey}
              className="flex-1 px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium"
            >
              📋 Copy to Clipboard
            </button>
            <button
              onClick={() => {
                setCreatedKey(null);
                onSuccess();
              }}
              className="flex-1 px-4 py-3 border-2 border-gray-300 rounded-lg hover:bg-gray-50 font-medium"
            >
              Done
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg p-6 sm:p-8 max-w-md w-full">
        <h2 className="text-2xl font-bold mb-4">Generate New API Key</h2>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          {error && <div className="bg-red-50 text-red-600 p-3 rounded text-sm">{error}</div>}
          
          <div>
            <label className="block text-sm font-medium mb-1">Key Name</label>
            <input
              type="text"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g., Production Key, Development Key"
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
            <p className="text-xs text-gray-500 mt-1">Choose a descriptive name to identify this key</p>
          </div>

          <div className="flex space-x-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 font-medium transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 font-medium transition-colors"
            >
              {loading ? 'Generating...' : 'Generate Key'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
