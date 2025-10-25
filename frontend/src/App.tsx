import { useState, useEffect } from 'react';
import { checkHealth } from './services/api';
import FinancialOverview from './components/Dashboard/FinancialOverview';
import CategoryChart from './components/Charts/CategoryChart';
import TimelineChart from './components/Charts/TimelineChart';

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
      <header className="bg-white shadow-sm sticky top-0 z-10">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-gray-900">
                💰 Stori Financial Tracker
              </h1>
              <p className="text-sm text-gray-600 mt-1">
                Your personal finance dashboard with AI-powered insights
              </p>
            </div>
            <div className="flex items-center gap-2 text-sm">
              <span className="text-gray-600">API:</span>
              <span className={`font-medium ${apiStatus.includes('✅') ? 'text-green-600' : 'text-red-600'}`}>
                {apiStatus}
              </span>
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Financial Overview Section */}
        <section className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-4">📊 Financial Overview</h2>
          <FinancialOverview />
        </section>

        {/* Charts Section */}
        <section className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-4">📈 Analytics</h2>
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <CategoryChart />
            <TimelineChart />
          </div>
        </section>

        {/* AI Advisor Placeholder */}
        <section className="mb-8">
          <h2 className="text-2xl font-bold text-gray-900 mb-4">🤖 AI Financial Advisor</h2>
          <div className="card bg-linear-to-br from-purple-50 to-blue-50 border-2 border-dashed border-purple-300">
            <div className="text-center py-8">
              <p className="text-lg font-medium text-purple-900 mb-2">
                🚀 Coming Soon!
              </p>
              <p className="text-sm text-purple-700">
                Get personalized financial advice powered by AI
              </p>
            </div>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 mt-8 border-t border-gray-200">
        <div className="text-center text-sm text-gray-500">
          <p>Built with React + TypeScript + Vite + TailwindCSS + Recharts</p>
          <p className="mt-1">Stori Full Stack Challenge © 2025</p>
        </div>
      </footer>
    </div>
  );
}

export default App;

