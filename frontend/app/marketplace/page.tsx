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

export default function MarketplacePage() {
  const router = useRouter();
  const { user, logout } = useAuth();
  const [apis, setAPIs] = useState<API[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');
  const [filterRuntime, setFilterRuntime] = useState('all');
  const [selectedAPI, setSelectedAPI] = useState<API | null>(null);

  useEffect(() => {
    loadAPIs();
  }, []);

  const loadAPIs = async () => {
    try {
      const data: any = await api.getPublicAPIs();
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

  // Filter APIs based on search and runtime
  const filteredAPIs = apis.filter((apiItem) => {
    const matchesSearch = apiItem.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         apiItem.description.toLowerCase().includes(searchTerm.toLowerCase());
    const matchesRuntime = filterRuntime === 'all' || apiItem.runtime === filterRuntime;
    return matchesSearch && matchesRuntime;
  });

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      {/* Navigation */}
      <nav className="bg-white shadow-md sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6">
          <div className="flex justify-between h-16 items-center">
            <Link href="/" className="flex items-center space-x-2">
              <div className="w-8 h-8 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-lg flex items-center justify-center">
                <span className="text-white font-bold text-sm">AP</span>
              </div>
              <span className="text-lg sm:text-xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
                API Marketplace
              </span>
            </Link>
            <div className="flex items-center space-x-2 sm:space-x-4">
              {user ? (
                <>
                  <Link href="/dashboard" className="hidden sm:inline text-gray-600 hover:text-gray-900 font-medium">
                    Dashboard
                  </Link>
                  <span className="hidden md:inline text-gray-600 text-sm">Hi, {user.name}</span>
                  <button
                    onClick={handleLogout}
                    className="px-3 sm:px-4 py-2 text-sm text-red-600 hover:text-red-700 font-medium"
                  >
                    Logout
                  </button>
                </>
              ) : (
                <>
                  <Link href="/login" className="px-4 py-2 text-sm text-gray-600 hover:text-gray-900 font-medium">
                    Sign in
                  </Link>
                  <Link href="/signup" className="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium">
                    Get Started
                  </Link>
                </>
              )}
            </div>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 py-8">
        {/* Header */}
        <div className="text-center mb-8">
          <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 mb-4">
            Discover Public APIs
          </h1>
          <p className="text-lg text-gray-600 max-w-2xl mx-auto">
            Browse and integrate production-ready APIs from our community
          </p>
        </div>

        {/* Search and Filters */}
        <div className="bg-white rounded-lg shadow-md p-4 sm:p-6 mb-6">
          <div className="flex flex-col sm:flex-row gap-4">
            <div className="flex-1">
              <input
                type="text"
                placeholder="Search APIs by name or description..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
            <select
              value={filterRuntime}
              onChange={(e) => setFilterRuntime(e.target.value)}
              className="px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="all">All Runtimes</option>
              <option value="python">Python</option>
              <option value="nodejs">Node.js</option>
              <option value="go">Go</option>
            </select>
          </div>
          {(searchTerm || filterRuntime !== 'all') && (
            <div className="mt-4 text-sm text-gray-600">
              Showing {filteredAPIs.length} API{filteredAPIs.length !== 1 ? 's' : ''}
            </div>
          )}
        </div>

        {/* API Grid */}
        {filteredAPIs.length === 0 ? (
          <div className="text-center py-16 bg-white rounded-lg shadow-md">
            <div className="text-6xl mb-4">🔍</div>
            <h3 className="text-xl font-semibold mb-2 text-gray-900">No APIs Found</h3>
            <p className="text-gray-600">
              {searchTerm || filterRuntime !== 'all'
                ? 'Try adjusting your search or filters'
                : 'Be the first to publish a public API!'}
            </p>
          </div>
        ) : (
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredAPIs.map((apiItem) => (
              <div
                key={apiItem.id}
                onClick={() => setSelectedAPI(apiItem)}
                className="bg-white rounded-lg border-2 border-gray-200 p-6 hover:shadow-xl hover:border-blue-500 transition-all cursor-pointer group"
              >
                <div className="mb-4">
                  <h3 className="text-xl font-bold text-gray-900 group-hover:text-blue-600 transition-colors mb-2">
                    {apiItem.name}
                  </h3>
                  <p className="text-gray-600 text-sm line-clamp-2">
                    {apiItem.description || 'No description available'}
                  </p>
                </div>

                <div className="flex flex-wrap gap-2 mb-4">
                  <span className="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-xs font-medium">
                    {apiItem.runtime}
                  </span>
                  <span className="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-xs font-medium">
                    v{apiItem.version}
                  </span>
                  <span className="px-3 py-1 bg-green-100 text-green-700 rounded-full text-xs font-medium">
                    {apiItem.status}
                  </span>
                </div>

                <button className="w-full py-2 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg hover:from-blue-700 hover:to-indigo-700 font-medium text-sm transition-all">
                  View Details →
                </button>
              </div>
            ))}
          </div>
        )}
      </div>

      {/* API Detail Modal */}
      {selectedAPI && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4" onClick={() => setSelectedAPI(null)}>
          <div className="bg-white rounded-lg p-6 sm:p-8 max-w-2xl w-full max-h-[90vh] overflow-y-auto" onClick={(e) => e.stopPropagation()}>
            <div className="flex justify-between items-start mb-6">
              <div>
                <h2 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">{selectedAPI.name}</h2>
                <p className="text-gray-600">{selectedAPI.description || 'No description available'}</p>
              </div>
              <button
                onClick={() => setSelectedAPI(null)}
                className="text-gray-500 hover:text-gray-700 text-3xl leading-none"
              >
                ×
              </button>
            </div>

            <div className="space-y-4 mb-6">
              <div>
                <h3 className="font-semibold text-gray-900 mb-2">Runtime</h3>
                <span className="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm font-medium">
                  {selectedAPI.runtime}
                </span>
              </div>

              <div>
                <h3 className="font-semibold text-gray-900 mb-2">Version</h3>
                <p className="text-gray-700">v{selectedAPI.version}</p>
              </div>

              <div>
                <h3 className="font-semibold text-gray-900 mb-2">Endpoint</h3>
                <code className="block bg-gray-100 px-4 py-3 rounded-lg text-sm break-all">
                  {selectedAPI.endpoint}
                </code>
              </div>

              <div>
                <h3 className="font-semibold text-gray-900 mb-2">Status</h3>
                <span className="px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm font-medium">
                  {selectedAPI.status}
                </span>
              </div>
            </div>

            {user ? (
              <div className="flex gap-3">
                <button className="flex-1 py-3 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg hover:from-blue-700 hover:to-indigo-700 font-medium transition-all shadow-md">
                  Subscribe to API
                </button>
                <button 
                  onClick={() => router.push(`/marketplace/${selectedAPI.id}`)}
                  className="px-6 py-3 border-2 border-gray-300 text-gray-700 rounded-lg hover:border-blue-600 hover:text-blue-600 font-medium transition-all"
                >
                  Test API
                </button>
              </div>
            ) : (
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 text-center">
                <p className="text-blue-900 mb-3">Sign in to use this API</p>
                <Link
                  href="/login"
                  className="inline-block px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium"
                >
                  Sign In
                </Link>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
