'use client';

import Link from 'next/link';
import { useAuth } from '@/lib/auth-context';
import { useRouter } from 'next/navigation';

export default function Home() {
  const { user, loading, logout } = useAuth();
  const router = useRouter();

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  const handleLogout = () => {
    logout();
    router.push('/');
  };

  return (
    <main className="min-h-screen bg-gradient-to-b from-white to-gray-50">
      {/* Responsive Navigation */}
      <nav className="bg-white/80 backdrop-blur-sm border-b border-gray-200 px-4 sm:px-6 py-4 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <Link href="/" className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-sm">AP</span>
            </div>
            <span className="text-lg sm:text-xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">
              API Platform
            </span>
          </Link>
          
          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center space-x-6">
            <Link href="/marketplace" className="text-gray-600 hover:text-gray-900 font-medium transition-colors">
              Marketplace
            </Link>
            {user ? (
              <>
                <Link href="/dashboard" className="text-gray-600 hover:text-gray-900 font-medium transition-colors">
                  Dashboard
                </Link>
                <span className="text-gray-600 text-sm">Hi, {user.name}</span>
                <button
                  onClick={handleLogout}
                  className="px-4 py-2 text-gray-600 hover:text-red-600 font-medium transition-colors"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link href="/login" className="text-gray-600 hover:text-gray-900 font-medium transition-colors">
                  Sign in
                </Link>
                <Link href="/signup" className="px-5 py-2.5 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg hover:from-blue-700 hover:to-indigo-700 font-medium transition-all shadow-md hover:shadow-lg">
                  Get Started →
                </Link>
              </>
            )}
          </div>

          {/* Mobile Menu Button */}
          <div className="md:hidden">
            {user ? (
              <Link href="/dashboard" className="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg">
                Dashboard
              </Link>
            ) : (
              <Link href="/login" className="px-4 py-2 text-sm bg-blue-600 text-white rounded-lg">
                Sign in
              </Link>
            )}
          </div>
        </div>
      </nav>

      {/* Hero Section - Responsive */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 py-12 sm:py-24">
        <div className="text-center mb-12 sm:mb-20">
          <div className="inline-block mb-4 px-4 py-2 bg-blue-100 text-blue-700 rounded-full text-xs sm:text-sm font-semibold">
            🚀 MENA-First API Platform
          </div>
          <h1 className="text-4xl sm:text-5xl md:text-6xl lg:text-7xl font-bold mb-4 sm:mb-6 bg-gradient-to-r from-gray-900 via-blue-900 to-indigo-900 bg-clip-text text-transparent px-4">
            Deploy & Monetize
            <br />Your APIs Instantly
          </h1>
          <p className="text-lg sm:text-xl md:text-2xl text-gray-600 mb-6 sm:mb-10 max-w-3xl mx-auto px-4">
            Serverless API hosting platform with built-in marketplace, analytics, and billing
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center px-4">
            {user ? (
              <>
                <Link href="/dashboard" className="w-full sm:w-auto px-6 sm:px-8 py-3 sm:py-4 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg hover:from-blue-700 hover:to-indigo-700 font-semibold text-base sm:text-lg transition-all shadow-lg hover:shadow-xl text-center">
                  Go to Dashboard →
                </Link>
                <Link href="/marketplace" className="w-full sm:w-auto px-6 sm:px-8 py-3 sm:py-4 border-2 border-gray-300 rounded-lg hover:bg-white hover:border-blue-600 font-semibold text-base sm:text-lg transition-all text-center">
                  Browse APIs
                </Link>
              </>
            ) : (
              <>
                <Link href="/signup" className="w-full sm:w-auto px-6 sm:px-8 py-3 sm:py-4 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg hover:from-blue-700 hover:to-indigo-700 font-semibold text-base sm:text-lg transition-all shadow-lg hover:shadow-xl text-center">
                  Start Building Free →
                </Link>
                <Link href="/marketplace" className="w-full sm:w-auto px-6 sm:px-8 py-3 sm:py-4 border-2 border-gray-300 rounded-lg hover:bg-white hover:border-blue-600 font-semibold text-base sm:text-lg transition-all text-center">
                  Explore Marketplace
                </Link>
              </>
            )}
          </div>
          <p className="mt-4 sm:mt-6 text-xs sm:text-sm text-gray-500">No credit card required • Deploy in minutes</p>
        </div>
        
        {/* Feature Cards - Responsive Grid */}
        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6 sm:gap-8 px-4">
          <div className="group p-6 sm:p-8 bg-white border-2 border-gray-200 rounded-2xl hover:border-blue-500 hover:shadow-xl transition-all">
            <div className="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center mb-4 group-hover:bg-blue-500 transition-colors">
              <span className="text-2xl group-hover:text-white transition-colors">👨‍💻</span>
            </div>
            <h2 className="text-xl sm:text-2xl font-bold mb-3 text-gray-900">For Developers</h2>
            <p className="text-sm sm:text-base text-gray-600 leading-relaxed">
              Deploy Python, Node.js, or Go APIs instantly. Start earning from your code with built-in billing and analytics.
            </p>
            <Link href="/signup" className="inline-block mt-4 text-blue-600 font-semibold hover:text-blue-700 text-sm sm:text-base">
              Learn more →
            </Link>
          </div>
          
          <div className="group p-6 sm:p-8 bg-white border-2 border-gray-200 rounded-2xl hover:border-indigo-500 hover:shadow-xl transition-all">
            <div className="w-12 h-12 bg-indigo-100 rounded-xl flex items-center justify-center mb-4 group-hover:bg-indigo-500 transition-colors">
              <span className="text-2xl group-hover:text-white transition-colors">🛒</span>
            </div>
            <h2 className="text-xl sm:text-2xl font-bold mb-3 text-gray-900">For Consumers</h2>
            <p className="text-sm sm:text-base text-gray-600 leading-relaxed">
              Browse and integrate production-ready APIs from our marketplace. Pay only for what you use.
            </p>
            <Link href="/marketplace" className="inline-block mt-4 text-indigo-600 font-semibold hover:text-indigo-700 text-sm sm:text-base">
              Browse APIs →
            </Link>
          </div>
          
          <div className="group p-6 sm:p-8 bg-white border-2 border-gray-200 rounded-2xl hover:border-purple-500 hover:shadow-xl transition-all sm:col-span-2 lg:col-span-1">
            <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center mb-4 group-hover:bg-purple-500 transition-colors">
              <span className="text-2xl group-hover:text-white transition-colors">🌍</span>
            </div>
            <h2 className="text-xl sm:text-2xl font-bold mb-3 text-gray-900">MENA-First</h2>
            <p className="text-sm sm:text-base text-gray-600 leading-relaxed">
              Built for the Middle East with local hosting, multi-currency support, and regional payment gateways.
            </p>
            <span className="inline-block mt-4 text-purple-600 font-semibold text-lg sm:text-xl">🇦🇪 🇸🇦 🇪🇬 🇱🇧</span>
          </div>
        </div>
        
        {/* Features Section - Responsive */}
        <div className="mt-16 sm:mt-32 text-center px-4">
          <h2 className="text-3xl sm:text-4xl font-bold mb-4 text-gray-900">Everything you need to succeed</h2>
          <p className="text-lg sm:text-xl text-gray-600 mb-8 sm:mb-16">Focus on building great APIs, we'll handle the rest</p>
          
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
            <div className="p-4 sm:p-6 text-left">
              <div className="text-2xl sm:text-3xl mb-2 sm:mb-3">⚡</div>
              <h3 className="font-bold text-base sm:text-lg mb-1 sm:mb-2">Instant Deploy</h3>
              <p className="text-gray-600 text-xs sm:text-sm">Upload code and go live in seconds</p>
            </div>
            <div className="p-4 sm:p-6 text-left">
              <div className="text-2xl sm:text-3xl mb-2 sm:mb-3">📊</div>
              <h3 className="font-bold text-base sm:text-lg mb-1 sm:mb-2">Real-time Analytics</h3>
              <p className="text-gray-600 text-xs sm:text-sm">Track usage, performance, and revenue</p>
            </div>
            <div className="p-4 sm:p-6 text-left">
              <div className="text-2xl sm:text-3xl mb-2 sm:mb-3">💳</div>
              <h3 className="font-bold text-base sm:text-lg mb-1 sm:mb-2">Built-in Billing</h3>
              <p className="text-gray-600 text-xs sm:text-sm">Automatic usage-based invoicing</p>
            </div>
            <div className="p-4 sm:p-6 text-left">
              <div className="text-2xl sm:text-3xl mb-2 sm:mb-3">🔒</div>
              <h3 className="font-bold text-base sm:text-lg mb-1 sm:mb-2">Secure by Default</h3>
              <p className="text-gray-600 text-xs sm:text-sm">API keys, rate limiting, and HTTPS</p>
            </div>
          </div>
        </div>
      </div>
    </main>
  );
}
