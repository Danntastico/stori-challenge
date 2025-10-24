import { useState, useEffect } from 'react';
import { checkHealth } from './services/api';

function App() {
  const [apiStatus, setApiStatus] = useState<string>('checking...');

  useEffect(() => {
    // Check API health on mount
    checkHealth()
      .then((data) => {
        setApiStatus('✅ Connected');
        console.log('API Health:', data);
      })
      .catch((error) => {
        setApiStatus('❌ Disconnected');
        console.error('API Health Check Failed:', error);
      });
  }, []);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white shadow-sm">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold text-gray-900">
              💰 Stori Financial Tracker
            </h1>
            <div className="flex items-center gap-2 text-sm text-gray-600">
              <span>API:</span>
              <span className="font-medium">{apiStatus}</span>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="card">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">
            🚀 Welcome to Stori Financial Tracker
          </h2>
          <p className="text-gray-600 mb-4">
            Your personal finance dashboard with AI-powered insights.
          </p>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
            <div className="p-4 bg-blue-50 rounded-lg">
              <h3 className="font-semibold text-blue-900 mb-2">📊 Dashboard</h3>
              <p className="text-sm text-blue-700">
                View your financial overview and key metrics
              </p>
            </div>
            
            <div className="p-4 bg-green-50 rounded-lg">
              <h3 className="font-semibold text-green-900 mb-2">📈 Analytics</h3>
              <p className="text-sm text-green-700">
                Analyze spending patterns and trends
              </p>
            </div>
            
            <div className="p-4 bg-purple-50 rounded-lg">
              <h3 className="font-semibold text-purple-900 mb-2">🤖 AI Advisor</h3>
              <p className="text-sm text-purple-700">
                Get personalized financial advice
              </p>
            </div>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 mt-8">
        <div className="text-center text-sm text-gray-500">
          Built with React + TypeScript + Vite + TailwindCSS + Recharts
        </div>
      </footer>
    </div>
  );
}

export default App;

