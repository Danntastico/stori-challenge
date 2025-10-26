import { useState } from 'react';
import { getAIAdvice, AIAdviceResponse } from '../../services/api';
import { Card, Button, LoadingSkeleton, ErrorAlert } from '../common';

export default function AIAdvisor() {
  const [advice, setAdvice] = useState<AIAdviceResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleGetAdvice = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const response = await getAIAdvice({ context: 'general' });
      setAdvice(response);
    } catch (err) {
      setError('Failed to generate financial advice. Please try again.');
      console.error('Error fetching advice:', err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <Card>
        <div className="text-center py-8">
          <LoadingSkeleton variant="text" count={5} />
          <p className="text-gray-600 mt-4">Analyzing your finances...</p>
        </div>
      </Card>
    );
  }

  if (!advice) {
    return (
      <Card gradient className="text-center py-12">
        <div className="max-w-2xl mx-auto">
          <div className="text-6xl mb-4">🤖</div>
          <h3 className="text-2xl font-bold text-gray-900 mb-3">
            AI Financial Advisor
          </h3>
          <p className="text-gray-700 mb-6">
            Get personalized insights and recommendations based on your spending patterns.
            Our AI analyzes your financial data to provide actionable advice.
          </p>
          <Button
            variant="primary"
            size="lg"
            onClick={handleGetAdvice}
            className="shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 transition-all"
          >
            Get Financial Advice
          </Button>
        </div>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header with Refresh Button */}
      <div className="flex items-center justify-between">
        <h3 className="text-xl font-semibold text-gray-900 flex items-center gap-2">
          <span>🤖</span>
          <span>Your Personalized Financial Advice</span>
        </h3>
        <Button
          variant="outline"
          size="sm"
          onClick={handleGetAdvice}
          className="flex items-center gap-2"
        >
          <span>🔄</span>
          <span>Refresh Advice</span>
        </Button>
      </div>

      {error && <ErrorAlert message={error} onRetry={handleGetAdvice} />}

      {/* Insights Section */}
      <Card hover>
        <div className="flex items-start gap-3 mb-4">
          <span className="text-3xl">💡</span>
          <div className="flex-1">
            <h4 className="text-lg font-semibold text-gray-900 mb-3">
              Key Insights
            </h4>
            <ul className="space-y-3">
              {advice.insights.map((insight, index) => (
                <li key={index} className="flex items-start gap-2">
                  <span className="text-blue-600 font-bold mt-1">•</span>
                  <span className="text-gray-700 flex-1">{insight}</span>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </Card>

      {/* Recommendations Section */}
      <Card hover>
        <div className="flex items-start gap-3 mb-4">
          <span className="text-3xl">🎯</span>
          <div className="flex-1">
            <h4 className="text-lg font-semibold text-gray-900 mb-3">
              Recommendations
            </h4>
            <ul className="space-y-3">
              {advice.recommendations.map((recommendation, index) => (
                <li key={index} className="flex items-start gap-3">
                  <span className="shrink-0 w-6 h-6 rounded-full bg-green-100 text-green-700 flex items-center justify-center text-sm font-bold">
                    {index + 1}
                  </span>
                  <span className="text-gray-700 flex-1">{recommendation}</span>
                </li>
              ))}
            </ul>
          </div>
        </div>
      </Card>

      {/* Full Advice Section (Collapsible) */}
      <Card>
        <details className="group">
          <summary className="cursor-pointer flex items-center justify-between p-4 hover:bg-gray-50 rounded-lg transition-colors">
            <span className="font-medium text-gray-900 flex items-center gap-2">
              <span>📋</span>
              <span>View Full Analysis</span>
            </span>
            <span className="text-gray-500 group-open:rotate-180 transition-transform">
              ▼
            </span>
          </summary>
          <div className="mt-4 pt-4 border-t border-gray-200">
            <pre className="whitespace-pre-wrap text-sm text-gray-700 font-sans">
              {advice.advice}
            </pre>
            <p className="text-xs text-gray-500 mt-4 italic">
              Generated: {new Date(advice.timestamp).toLocaleString()}
            </p>
          </div>
        </details>
      </Card>
    </div>
  );
}

