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
  const [editingAPI, setEditingAPI] = useState<API | null>(null);

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
      <nav className="bg-white border-b px-4 sm:px-6 py-4">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <Link href="/" className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-sm">AP</span>
            </div>
            <span className="text-lg sm:text-2xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
              API Platform
            </span>
          </Link>
          <div className="flex items-center space-x-2 sm:space-x-4">
            <Link href="/marketplace" className="hidden sm:inline text-gray-600 hover:text-gray-900 font-medium">
              Marketplace
            </Link>
            <Link href="/dashboard/api-keys" className="hidden sm:inline text-gray-600 hover:text-gray-900 font-medium">
              🔑 API Keys
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
        <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6 sm:mb-8">
          <h1 className="text-2xl sm:text-3xl font-bold">My APIs</h1>
          <button
            onClick={() => setShowCreateModal(true)}
            className="w-full sm:w-auto px-4 py-2 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg hover:from-blue-700 hover:to-indigo-700 font-medium transition-all shadow-md hover:shadow-lg"
          >
            + Create API
          </button>
        </div>

        {apis.length === 0 ? (
          <div className="text-center py-12 bg-white rounded-lg border mx-4 sm:mx-0">
            <div className="text-5xl mb-4">🚀</div>
            <h3 className="text-xl font-semibold mb-2">No APIs yet</h3>
            <p className="text-gray-600 mb-4 px-4">Create your first API to get started</p>
            <button
              onClick={() => setShowCreateModal(true)}
              className="px-6 py-3 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg hover:from-blue-700 hover:to-indigo-700 font-medium shadow-md"
            >
              Create Your First API
            </button>
          </div>
        ) : (
          <div className="grid gap-4 sm:gap-6">
            {apis.map((apiItem) => (
              <div key={apiItem.id} className="relative bg-white p-4 sm:p-6 rounded-lg border-2 border-gray-200 hover:shadow-xl transition-all hover:border-blue-500 group">
                {/* Edit Button */}
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    setEditingAPI(apiItem);
                  }}
                  className="absolute top-4 right-4 p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all"
                  title="Edit API"
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                </button>

                <Link href={`/dashboard/api/${apiItem.id}`} className="block">
                  <div className="flex justify-between items-start mb-4 pr-10">
                    <div className="flex-1">
                      <h3 className="text-xl font-semibold group-hover:text-blue-600 transition-colors">
                        {apiItem.name}
                      </h3>
                      <p className="text-gray-600 mt-1">{apiItem.description || 'No description'}</p>
                      {apiItem.status === 'pending' && (
                        <p className="text-xs text-yellow-600 mt-2 flex items-center gap-1">
                          <span>⚠️</span>
                          <span>Upload code to activate this API</span>
                        </p>
                      )}
                    </div>
                    <span
                      className={`px-3 py-1 rounded-full text-sm font-medium whitespace-nowrap ml-4 ${
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
                <div className="flex flex-wrap gap-2 sm:gap-4 text-sm text-gray-600 mb-3">
                  <span className="flex items-center gap-1">
                    <span className="font-medium">Runtime:</span> {apiItem.runtime}
                  </span>
                  <span className="hidden sm:inline text-gray-400">•</span>
                  <span className="flex items-center gap-1">
                    <span className="font-medium">Version:</span> {apiItem.version}
                  </span>
                  <span className="hidden sm:inline text-gray-400">•</span>
                  <span className="flex items-center gap-1">
                    <span className="font-medium">Visibility:</span> {apiItem.visibility}
                  </span>
                </div>
                  <div className="flex justify-between items-center">
                    <code className="text-xs sm:text-sm bg-gray-100 px-3 py-1.5 rounded break-all text-gray-700">
                      {apiItem.endpoint}
                    </code>
                    <span className="ml-4 text-blue-600 font-medium text-sm whitespace-nowrap group-hover:translate-x-1 transition-transform">
                      Open →
                    </span>
                  </div>
                </Link>
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

        {editingAPI && (
          <EditAPIModal
            api={editingAPI}
            onClose={() => setEditingAPI(null)}
            onSuccess={() => {
              setEditingAPI(null);
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

function EditAPIModal({ api: initialAPI, onClose, onSuccess }: { api: API; onClose: () => void; onSuccess: () => void }) {
  const [name, setName] = useState(initialAPI.name);
  const [description, setDescription] = useState(initialAPI.description);
  const [visibility, setVisibility] = useState(initialAPI.visibility);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await api.updateAPI(initialAPI.id, {
        name,
        description,
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
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg p-6 sm:p-8 max-w-md w-full">
        <h2 className="text-2xl font-bold mb-4">Edit API</h2>
        
        <form onSubmit={handleSubmit} className="space-y-4">
          {error && <div className="bg-red-50 text-red-600 p-3 rounded text-sm">{error}</div>}
          
          <div>
            <label className="block text-sm font-medium mb-1">Name</label>
            <input
              type="text"
              required
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Description</label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              rows={3}
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-1">Visibility</label>
            <select
              value={visibility}
              onChange={(e) => setVisibility(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="private">Private</option>
              <option value="public">Public (Free)</option>
              <option value="paid">Paid</option>
            </select>
          </div>

          <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 text-sm text-blue-800">
            <strong>Note:</strong> Runtime cannot be changed after creation.
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
              {loading ? 'Saving...' : 'Save Changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
