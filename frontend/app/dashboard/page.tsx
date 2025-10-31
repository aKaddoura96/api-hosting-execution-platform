'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/lib/auth-context';
import { api } from '@/lib/api';

interface API {
  id: string;
  name: string;
  description: string;
  version: string;
  runtime: string;
  visibility: string;
  status: string;
  endpoint: string;
  created_at: string;
}

export default function DashboardPage() {
  const router = useRouter();
  const { user, logout, loading: authLoading } = useAuth();
  const [apis, setAPIs] = useState<API[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);

  useEffect(() => {
    if (!authLoading && !user) {
      router.push('/login');
      return;
    }

    if (user) {
      loadAPIs();
    }
  }, [user, authLoading, router]);

  const loadAPIs = async () => {
    try {
      const data: any = await api.getMyAPIs();
      setAPIs(data || []);
    } catch (error) {
      console.error('Failed to load APIs:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    logout();
    router.push('/');
  };

  if (authLoading || loading) {
    return <div className="min-h-screen flex items-center justify-center">Loading...</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white border-b px-6 py-4">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <Link href="/" className="text-2xl font-bold">
            API Platform
          </Link>
          <div className="flex items-center space-x-4">
            <Link href="/marketplace" className="text-gray-600 hover:text-gray-900">
              Marketplace
            </Link>
            <span className="text-gray-600">{user?.name}</span>
            <button
              onClick={handleLogout}
              className="px-4 py-2 text-gray-600 hover:text-gray-900"
            >
              Logout
            </button>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-6 py-8">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">My APIs</h1>
          <button
            onClick={() => setShowCreateModal(true)}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            + Create API
          </button>
        </div>

        {apis.length === 0 ? (
          <div className="text-center py-12 bg-white rounded-lg border">
            <h3 className="text-xl font-semibold mb-2">No APIs yet</h3>
            <p className="text-gray-600 mb-4">Create your first API to get started</p>
            <button
              onClick={() => setShowCreateModal(true)}
              className="px-6 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Create Your First API
            </button>
          </div>
        ) : (
          <div className="grid gap-6">
            {apis.map((apiItem) => (
              <div key={apiItem.id} className="bg-white p-6 rounded-lg border">
                <div className="flex justify-between items-start mb-4">
                  <div>
                    <h3 className="text-xl font-semibold">{apiItem.name}</h3>
                    <p className="text-gray-600">{apiItem.description}</p>
                  </div>
                  <span
                    className={`px-3 py-1 rounded text-sm ${
                      apiItem.status === 'deployed'
                        ? 'bg-green-100 text-green-800'
                        : apiItem.status === 'pending'
                        ? 'bg-yellow-100 text-yellow-800'
                        : 'bg-red-100 text-red-800'
                    }`}
                  >
                    {apiItem.status}
                  </span>
                </div>
                <div className="flex space-x-4 text-sm text-gray-600">
                  <span>Runtime: {apiItem.runtime}</span>
                  <span>•</span>
                  <span>Version: {apiItem.version}</span>
                  <span>•</span>
                  <span>Visibility: {apiItem.visibility}</span>
                </div>
                <div className="mt-4">
                  <code className="text-sm bg-gray-100 px-3 py-1 rounded">
                    {apiItem.endpoint}
                  </code>
                </div>
              </div>
            ))}
          </div>
        )}

        {showCreateModal && (
          <CreateAPIModal
            onClose={() => setShowCreateModal(false)}
            onSuccess={() => {
              setShowCreateModal(false);
              loadAPIs();
            }}
          />
        )}
      </div>
    </div>
  );
}

function CreateAPIModal({ onClose, onSuccess }: { onClose: () => void; onSuccess: () => void }) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [runtime, setRuntime] = useState('python');
  const [visibility, setVisibility] = useState('private');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await api.createAPI({
        name,
        description,
        version: 'v1',
        runtime,
        visibility,
      });
      onSuccess();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-8 max-w-md w-full">
        <h2 className="text-2xl font-bold mb-4">Create New API</h2>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          {error && <div className="bg-red-50 text-red-600 p-3 rounded">{error}</div>}
          
          <div>
            <label className="block text-sm font-medium mb-1">Name</label>
            <input
              type="text"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-3 py-2 border rounded"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Description</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full px-3 py-2 border rounded"
              rows={3}
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Runtime</label>
            <select
              value={runtime}
              onChange={(e) => setRuntime(e.target.value)}
              className="w-full px-3 py-2 border rounded"
            >
              <option value="python">Python</option>
              <option value="nodejs">Node.js</option>
              <option value="go">Go</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Visibility</label>
            <select
              value={visibility}
              onChange={(e) => setVisibility(e.target.value)}
              className="w-full px-3 py-2 border rounded"
            >
              <option value="private">Private</option>
              <option value="public">Public (Free)</option>
              <option value="paid">Paid</option>
            </select>
          </div>

          <div className="flex space-x-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 border rounded hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="flex-1 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:opacity-50"
            >
              {loading ? 'Creating...' : 'Create'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
