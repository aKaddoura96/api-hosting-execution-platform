import Link from 'next/link'

export default function Home() {
  return (
    <main className="min-h-screen bg-gradient-to-b from-white to-gray-50">
      <nav className="bg-white/80 backdrop-blur-sm border-b border-gray-200 px-6 py-4 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <Link href="/" className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-gradient-to-br from-blue-600 to-indigo-600 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-sm">AP</span>
            </div>
            <span className="text-xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent">API Platform</span>
          </Link>
          <div className="flex items-center space-x-6">
            <Link href="/marketplace" className="text-gray-600 hover:text-gray-900 font-medium transition-colors">
              Marketplace
            </Link>
            <Link href="/login" className="text-gray-600 hover:text-gray-900 font-medium transition-colors">
              Sign in
            </Link>
            <Link href="/signup" className="px-5 py-2.5 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg hover:from-blue-700 hover:to-indigo-700 font-medium transition-all shadow-md hover:shadow-lg">
              Get Started →
            </Link>
          </div>
        </div>
      </nav>

      <div className="max-w-7xl mx-auto px-6 py-24">
        <div className="text-center mb-20">
          <div className="inline-block mb-4 px-4 py-2 bg-blue-100 text-blue-700 rounded-full text-sm font-semibold">
            🚀 MENA-First API Platform
          </div>
          <h1 className="text-6xl md:text-7xl font-bold mb-6 bg-gradient-to-r from-gray-900 via-blue-900 to-indigo-900 bg-clip-text text-transparent">
            Deploy & Monetize
            <br />Your APIs Instantly
          </h1>
          <p className="text-xl md:text-2xl text-gray-600 mb-10 max-w-3xl mx-auto">
            Serverless API hosting platform with built-in marketplace, analytics, and billing
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
            <Link href="/signup" className="px-8 py-4 bg-gradient-to-r from-blue-600 to-indigo-600 text-white rounded-lg hover:from-blue-700 hover:to-indigo-700 font-semibold text-lg transition-all shadow-lg hover:shadow-xl">
              Start Building Free →
            </Link>
            <Link href="/marketplace" className="px-8 py-4 border-2 border-gray-300 rounded-lg hover:bg-white hover:border-blue-600 font-semibold text-lg transition-all">
              Explore Marketplace
            </Link>
          </div>
          <p className="mt-6 text-sm text-gray-500">No credit card required • Deploy in minutes</p>
        </div>
        
        <div className="grid md:grid-cols-3 gap-8">
          <div className="group p-8 bg-white border-2 border-gray-200 rounded-2xl hover:border-blue-500 hover:shadow-xl transition-all">
            <div className="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center mb-4 group-hover:bg-blue-500 transition-colors">
              <span className="text-2xl group-hover:text-white transition-colors">👨‍💻</span>
            </div>
            <h2 className="text-2xl font-bold mb-3 text-gray-900">For Developers</h2>
            <p className="text-gray-600 leading-relaxed">
              Deploy Python, Node.js, or Go APIs instantly. Start earning from your code with built-in billing and analytics.
            </p>
            <Link href="/signup" className="inline-block mt-4 text-blue-600 font-semibold hover:text-blue-700">Learn more →</Link>
          </div>
          
          <div className="group p-8 bg-white border-2 border-gray-200 rounded-2xl hover:border-indigo-500 hover:shadow-xl transition-all">
            <div className="w-12 h-12 bg-indigo-100 rounded-xl flex items-center justify-center mb-4 group-hover:bg-indigo-500 transition-colors">
              <span className="text-2xl group-hover:text-white transition-colors">🛒</span>
            </div>
            <h2 className="text-2xl font-bold mb-3 text-gray-900">For Consumers</h2>
            <p className="text-gray-600 leading-relaxed">
              Browse and integrate production-ready APIs from our marketplace. Pay only for what you use.
            </p>
            <Link href="/marketplace" className="inline-block mt-4 text-indigo-600 font-semibold hover:text-indigo-700">Browse APIs →</Link>
          </div>
          
          <div className="group p-8 bg-white border-2 border-gray-200 rounded-2xl hover:border-purple-500 hover:shadow-xl transition-all">
            <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center mb-4 group-hover:bg-purple-500 transition-colors">
              <span className="text-2xl group-hover:text-white transition-colors">🌍</span>
            </div>
            <h2 className="text-2xl font-bold mb-3 text-gray-900">MENA-First</h2>
            <p className="text-gray-600 leading-relaxed">
              Built for the Middle East with local hosting, multi-currency support, and regional payment gateways.
            </p>
            <span className="inline-block mt-4 text-purple-600 font-semibold">🇦🇪 🇸🇦 🇪🇬 🇱🇧</span>
          </div>
        </div>
        
        {/* Features Section */}
        <div className="mt-32 text-center">
          <h2 className="text-4xl font-bold mb-4 text-gray-900">Everything you need to succeed</h2>
          <p className="text-xl text-gray-600 mb-16">Focus on building great APIs, we'll handle the rest</p>
          
          <div className="grid md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div className="p-6 text-left">
              <div className="text-3xl mb-3">⚡</div>
              <h3 className="font-bold text-lg mb-2">Instant Deploy</h3>
              <p className="text-gray-600 text-sm">Upload code and go live in seconds</p>
            </div>
            <div className="p-6 text-left">
              <div className="text-3xl mb-3">📊</div>
              <h3 className="font-bold text-lg mb-2">Real-time Analytics</h3>
              <p className="text-gray-600 text-sm">Track usage, performance, and revenue</p>
            </div>
            <div className="p-6 text-left">
              <div className="text-3xl mb-3">💳</div>
              <h3 className="font-bold text-lg mb-2">Built-in Billing</h3>
              <p className="text-gray-600 text-sm">Automatic usage-based invoicing</p>
            </div>
            <div className="p-6 text-left">
              <div className="text-3xl mb-3">🔒</div>
              <h3 className="font-bold text-lg mb-2">Secure by Default</h3>
              <p className="text-gray-600 text-sm">API keys, rate limiting, and HTTPS</p>
            </div>
          </div>
        </div>
      </div>
    </main>
  )
}
