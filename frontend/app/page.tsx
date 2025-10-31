export default function Home() {
  return (
    <main className="min-h-screen p-8">
      <div className="max-w-6xl mx-auto">
        <h1 className="text-4xl font-bold mb-4">
          API Marketplace & Hosting Platform
        </h1>
        <p className="text-xl text-gray-600 mb-8">
          Host, monetize, and consume APIs with ease
        </p>
        
        <div className="grid md:grid-cols-3 gap-6">
          <div className="p-6 border rounded-lg">
            <h2 className="text-2xl font-semibold mb-2">For Developers</h2>
            <p className="text-gray-600">
              Deploy your APIs instantly and start earning from your code
            </p>
          </div>
          
          <div className="p-6 border rounded-lg">
            <h2 className="text-2xl font-semibold mb-2">For Consumers</h2>
            <p className="text-gray-600">
              Browse and integrate ready-to-use APIs from our marketplace
            </p>
          </div>
          
          <div className="p-6 border rounded-lg">
            <h2 className="text-2xl font-semibold mb-2">Regional Focus</h2>
            <p className="text-gray-600">
              MENA-first platform with local hosting and payment options
            </p>
          </div>
        </div>
      </div>
    </main>
  )
}
