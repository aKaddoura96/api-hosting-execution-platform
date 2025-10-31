'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
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
  const [apis, setAPIs] = useState<API[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');

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

  const filteredAPIs = apis.filter((apiItem) =>
    apiItem.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    apiItem.description.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white border-b px-6 py-4">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <Link href="/" className="text-2xl font-bold">
            API Platform
          </Link>
          <div className="flex items-center space-x-4">
            <Link href="/dashboard" className="text-gray-600 hover:text-gray-900">
              Dashboard
            </Link>
            <Link href="/login" className="text-gray-600 hover:text-gray-900">
              Sign in
            </Link>
            <Link href="/signup" className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
              Get Started
            </Link>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-6 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-4">API Marketplace</h1>
          <p className="text-gray-600 mb-6">
            Browse and integrate ready-to-use APIs from our community
          </p>
          
          <input
            type="text"
            placeholder="Search APIs..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full max-w-lg px-4 py-2 border rounded-lg"
          />
        </div>

        {loading ? (
          <div className="text-center py-12">Loading...</div>
        ) : filteredAPIs.length === 0 ? (
          <div className="text-center py-12 bg-white rounded-lg border">
            <h3 className="text-xl font-semibold mb-2">No public APIs available yet</h3>
            <p className="text-gray-600">Be the first to publish a public API!</p>
          </div>
        ) : (
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredAPIs.map((apiItem) => (
              <div key={apiItem.id} className="bg-white p-6 rounded-lg border hover:shadow-lg transition">
                <div className="mb-4">
                  <h3 className="text-xl font-semibold mb-2">{apiItem.name}</h3>
                  <p className="text-gray-600 text-sm">{apiItem.description}</p>
                </div>
                
                <div className="space-y-2 text-sm text-gray-600 mb-4">
                  <div className="flex items-center">
                    <span className="font-medium mr-2">Runtime:</span>
                    <span className="px-2 py-1 bg-gray-100 rounded">{apiItem.runtime}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="font-medium mr-2">Version:</span>
                    <span>{apiItem.version}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="font-medium mr-2">Status:</span>
                    <span className="px-2 py-1 bg-green-100 text-green-800 rounded">
                      {apiItem.status}
                    </span>
                  </div>
                </div>

                <div className="mb-4">
                  <code className="text-xs bg-gray-100 px-2 py-1 rounded block overflow-x-auto">
                    {apiItem.endpoint}
                  </code>
                </div>

                <Link
                  href={`/marketplace/${apiItem.id}`}
                  className="block w-full text-center px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
                >
                  View Details
                </Link>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
