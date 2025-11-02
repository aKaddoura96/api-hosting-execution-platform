'use client';

import { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import Link from 'next/link';
import { api } from '@/lib/api';

interface API {
  id: string;
  name: string;
  description: string;
  version: string;
  runtime: string;
  status: string;
  visibility: string;
  created_at: string;
}

export default function APIDetailPage() {
  const router = useRouter();
  const params = useParams();
  const apiId = params.id as string;

  const [apiData, setApiData] = useState<API | null>(null);
  const [loading, setLoading] = useState(true);
  const [deploying, setDeploying] = useState(false);
  const [testCode, setTestCode] = useState('print("Hello from API!")');
  const [testResult, setTestResult] = useState<any>(null);
  const [testing, setTesting] = useState(false);

  useEffect(() => {
    loadAPI();
  }, [apiId]);

  const loadAPI = async () => {
    try {
      const data = await api.getAPI(apiId);
      setApiData(data);
      
      // Set default test code based on runtime
      if (data.runtime === 'nodejs') {
        setTestCode('console.log("Hello from API!");');
      } else if (data.runtime === 'go') {
        setTestCode('package main\n\nimport "fmt"\n\nfunc main() {\n\tfmt.Println("Hello from API!")\n}');
      }
    } catch (err: any) {
      alert('Failed to load API: ' + err.message);
      router.push('/dashboard');
    } finally {
      setLoading(false);
    }
  };

  const handleDeploy = async () => {
    if (!confirm('Deploy this API?')) return;
    
    setDeploying(true);
    try {
      // TODO: Implement deploy endpoint
      alert('Deployment coming soon!');
    } catch (err: any) {
      alert('Failed to deploy: ' + err.message);
    } finally {
      setDeploying(false);
    }
  };

  const handleTest = async () => {
    setTesting(true);
    setTestResult(null);
    
    try {
      const response = await fetch('http://localhost:8081/execute', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          code: testCode,
          runtime: apiData?.runtime || 'python',
        }),
      });
      
      const result = await response.json();
      setTestResult(result);
    } catch (err: any) {
      setTestResult({ error: err.message });
    } finally {
      setTesting(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (!apiData) return null;

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      {/* Navigation */}
      <nav className="bg-white shadow-md">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between h-16 items-center">
            <div className="flex items-center space-x-4">
              <Link href="/dashboard" className="text-blue-600 hover:text-blue-700">
                ← Back to Dashboard
              </Link>
            </div>
            <button
              onClick={() => {
                localStorage.removeItem('token');
                router.push('/login');
              }}
              className="px-4 py-2 text-red-600 hover:text-red-700"
            >
              Logout
            </button>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <div className="flex justify-between items-start">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">{apiData.name}</h1>
              <p className="text-gray-600 mt-2">{apiData.description}</p>
              <div className="flex items-center gap-4 mt-4">
                <span className="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm">
                  {apiData.runtime}
                </span>
                <span className="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm">
                  v{apiData.version}
                </span>
                <span className={`px-3 py-1 rounded-full text-sm ${
                  apiData.status === 'deployed' ? 'bg-green-100 text-green-700' :
                  apiData.status === 'stopped' ? 'bg-gray-100 text-gray-700' :
                  'bg-yellow-100 text-yellow-700'
                }`}>
                  {apiData.status}
                </span>
              </div>
            </div>
            <button
              onClick={handleDeploy}
              disabled={deploying}
              className="px-6 py-3 bg-gradient-to-r from-green-600 to-green-700 text-white rounded-lg hover:from-green-700 hover:to-green-800 font-medium transition-all shadow-md disabled:opacity-50"
            >
              {deploying ? 'Deploying...' : '🚀 Deploy'}
            </button>
          </div>
        </div>

        {/* Test Interface */}
        <div className="bg-white rounded-lg shadow-md p-6">
          <h2 className="text-2xl font-bold text-gray-900 mb-4">Test API</h2>
          
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Code
              </label>
              <textarea
                value={testCode}
                onChange={(e) => setTestCode(e.target.value)}
                rows={10}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg font-mono text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Enter your code here..."
              />
            </div>

            <button
              onClick={handleTest}
              disabled={testing}
              className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium transition-colors disabled:opacity-50"
            >
              {testing ? 'Running...' : '▶️ Run Code'}
            </button>

            {testResult && (
              <div className="mt-4">
                <h3 className="text-lg font-semibold mb-2">Result:</h3>
                <div className={`p-4 rounded-lg ${
                  testResult.error ? 'bg-red-50 border border-red-200' : 'bg-gray-50 border border-gray-200'
                }`}>
                  <pre className="whitespace-pre-wrap text-sm font-mono">
                    {testResult.error || testResult.output || JSON.stringify(testResult, null, 2)}
                  </pre>
                  {testResult.duration_ms && (
                    <p className="text-sm text-gray-600 mt-2">
                      Duration: {testResult.duration_ms}ms | Exit Code: {testResult.exit_code}
                    </p>
                  )}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
