import Link from 'next/link'

export default function Home() {
  return (
    <main className="min-h-screen">
      <nav className="bg-white border-b px-6 py-4">
        <div className="max-w-6xl mx-auto flex justify-between items-center">
          <h1 className="text-2xl font-bold">API Platform</h1>
          <div className="space-x-4">
            <Link href="/marketplace" className="text-gray-600 hover:text-gray-900">
              Marketplace
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

      <div className="max-w-6xl mx-auto px-6 py-20">
        <div className="text-center mb-16">
          <h1 className="text-5xl font-bold mb-4">
            API Marketplace & Hosting Platform
          </h1>
          <p className="text-xl text-gray-600 mb-8">
            Host, monetize, and consume APIs with ease
          </p>
          <div className="space-x-4">
            <Link href="/signup" className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700">
              Start Building
            </Link>
            <Link href="/marketplace" className="px-6 py-3 border border-gray-300 rounded-lg hover:bg-gray-50">
              Browse APIs
            </Link>
          </div>
        </div>
        
        <div className="grid md:grid-cols-3 gap-8">
          <div className="p-8 border rounded-lg">
            <h2 className="text-2xl font-semibold mb-3">For Developers</h2>
            <p className="text-gray-600">
              Deploy your APIs instantly and start earning from your code
            </p>
          </div>
          
          <div className="p-8 border rounded-lg">
            <h2 className="text-2xl font-semibold mb-3">For Consumers</h2>
            <p className="text-gray-600">
              Browse and integrate ready-to-use APIs from our marketplace
            </p>
          </div>
          
          <div className="p-8 border rounded-lg">
            <h2 className="text-2xl font-semibold mb-3">Regional Focus</h2>
            <p className="text-gray-600">
              MENA-first platform with local hosting and payment options
            </p>
          </div>
        </div>
      </div>
    </main>
  )
}
