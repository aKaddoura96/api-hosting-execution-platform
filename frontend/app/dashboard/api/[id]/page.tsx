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
  code_path?: string;
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
  const [uploading, setUploading] = useState(false);
  const [dragActive, setDragActive] = useState(false);
  const [uploadStatus, setUploadStatus] = useState<{type: 'success' | 'error', message: string} | null>(null);

  useEffect(() => {
    loadAPI();
  }, [apiId]);

  const loadAPI = async () => {
    try {
      const data = await api.getAPI(apiId) as API;
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
    if (apiData?.status === 'pending' && !apiData?.code_path) {
      alert('Please upload code before deploying.');
      return;
    }
    
    if (!confirm(`Deploy ${apiData?.name}? This will make your API live.`)) return;
    
    setDeploying(true);
    try {
      await api.deployAPI(apiId);
      alert('✅ API deployed successfully!');
      await loadAPI(); // Reload to get updated status
    } catch (err: any) {
      alert('❌ Failed to deploy: ' + err.message);
    } finally {
      setDeploying(false);
    }
  };

  const handleStop = async () => {
    if (!confirm(`Stop ${apiData?.name}? Users will no longer be able to access it.`)) return;
    
    try {
      await api.stopAPI(apiId);
      alert('API stopped successfully');
      await loadAPI();
    } catch (err: any) {
      alert('Failed to stop: ' + err.message);
    }
  };

  const handleDelete = async () => {
    if (!confirm(`Delete ${apiData?.name}? This action cannot be undone.`)) return;
    
    try {
      await api.deleteAPI(apiId);
      alert('API deleted successfully');
      router.push('/dashboard');
    } catch (err: any) {
      alert('Failed to delete: ' + err.message);
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

  const handleFileUpload = async (file: File) => {
    setUploading(true);
    setUploadStatus(null);
    
    try {
      await api.uploadCode(apiId, file);
      setUploadStatus({ type: 'success', message: `✅ ${file.name} uploaded successfully!` });
      // Reload API data to get updated code_path and status
      await loadAPI();
    } catch (err: any) {
      setUploadStatus({ type: 'error', message: `❌ Failed: ${err.message}` });
    } finally {
      setUploading(false);
    }
  };

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleFileUpload(e.dataTransfer.files[0]);
    }
  };

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      handleFileUpload(e.target.files[0]);
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
        <div className="max-w-7xl mx-auto px-4 sm:px-6">
          <div className="flex justify-between h-14 sm:h-16 items-center">
            <Link href="/dashboard" className="text-blue-600 hover:text-blue-700 text-sm sm:text-base font-medium">
              ← Back
            </Link>
            <button
              onClick={() => {
                localStorage.removeItem('token');
                router.push('/login');
              }}
              className="px-3 sm:px-4 py-2 text-sm sm:text-base text-red-600 hover:text-red-700 font-medium"
            >
              Logout
            </button>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 py-4 sm:py-8">
        {/* Header */}
        <div className="bg-white rounded-lg shadow-md p-4 sm:p-6 mb-4 sm:mb-6">
          <div className="flex flex-col sm:flex-row justify-between items-start gap-4">
            <div className="flex-1">
              <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">{apiData.name}</h1>
              <p className="text-sm sm:text-base text-gray-600 mt-2">{apiData.description}</p>
              <div className="flex flex-wrap items-center gap-2 sm:gap-4 mt-4">
                <span className="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm font-medium">
                  {apiData.runtime}
                </span>
                <span className="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm font-medium">
                  v{apiData.version}
                </span>
                <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                  apiData.status === 'deployed' ? 'bg-green-100 text-green-700' :
                  apiData.status === 'stopped' ? 'bg-gray-100 text-gray-700' :
                  'bg-yellow-100 text-yellow-700'
                }`}>
                  {apiData.status}
                </span>
              </div>
            </div>
            <div className="flex flex-col sm:flex-row gap-2 w-full sm:w-auto">
              {apiData.status === 'deployed' ? (
                <button
                  onClick={handleStop}
                  className="px-6 py-3 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700 font-medium transition-colors shadow-md"
                >
                  ⏸️ Stop
                </button>
              ) : (
                <button
                  onClick={handleDeploy}
                  disabled={deploying}
                  className="px-6 py-3 bg-gradient-to-r from-green-600 to-green-700 text-white rounded-lg hover:from-green-700 hover:to-green-800 font-medium transition-all shadow-md disabled:opacity-50"
                >
                  {deploying ? 'Deploying...' : '🚀 Deploy'}
                </button>
              )}
              <button
                onClick={handleDelete}
                className="px-6 py-3 bg-red-600 text-white rounded-lg hover:bg-red-700 font-medium transition-colors shadow-md"
              >
                🗑️ Delete
              </button>
            </div>
          </div>
        </div>

        {/* Upload Code */}
        <div className="bg-white rounded-lg shadow-md p-4 sm:p-6 mb-4 sm:mb-6">
          <h2 className="text-xl sm:text-2xl font-bold text-gray-900 mb-4">Upload Code</h2>
          
          <div
            onDragEnter={handleDrag}
            onDragLeave={handleDrag}
            onDragOver={handleDrag}
            onDrop={handleDrop}
            className={`border-2 border-dashed rounded-lg p-6 sm:p-8 text-center transition-colors ${
              dragActive ? 'border-blue-500 bg-blue-50' : 'border-gray-300 hover:border-gray-400'
            } ${uploading ? 'opacity-50 pointer-events-none' : ''}`}
          >
            <div className="space-y-4">
              <div className="text-4xl sm:text-5xl">📁</div>
              <div>
                <p className="text-base sm:text-lg font-medium text-gray-700">
                  {uploading ? 'Uploading...' : 'Drag & drop your code file here'}
                </p>
                <p className="text-sm text-gray-500 mt-2">or</p>
              </div>
              <label className="inline-block">
                <input
                  type="file"
                  onChange={handleFileInput}
                  disabled={uploading}
                  className="hidden"
                  accept=".py,.js,.go,.ts"
                />
                <span className="px-4 sm:px-6 py-2 sm:py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 font-medium transition-colors cursor-pointer inline-block text-sm sm:text-base">
                  Browse Files
                </span>
              </label>
              <p className="text-xs sm:text-sm text-gray-500">
                Supported: .py, .js, .go, .ts (max 10MB)
              </p>
            </div>
          </div>

          {uploadStatus && (
            <div className={`mt-4 p-3 sm:p-4 rounded-lg text-sm sm:text-base ${
              uploadStatus.type === 'success' ? 'bg-green-50 text-green-700 border border-green-200' : 'bg-red-50 text-red-700 border border-red-200'
            }`}>
              {uploadStatus.message}
            </div>
          )}
        </div>

        {/* Test Interface */}
        <div className="bg-white rounded-lg shadow-md p-4 sm:p-6">
          <h2 className="text-xl sm:text-2xl font-bold text-gray-900 mb-4">Test API</h2>
          
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
