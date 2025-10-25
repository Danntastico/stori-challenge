import { getCategorySummary, CategorySummaryResponse } from '../../services/api';
import { formatCurrency, formatPercentage } from '../../utils/formatters';
import { useApiData } from '../../hooks/useFetch';

interface StatCardProps {
  title: string;
  value: string;
  subtitle?: string;
  trend?: 'positive' | 'negative' | 'neutral';
  icon?: string;
}

function StatCard({ title, value, subtitle, trend, icon }: StatCardProps) {
  const trendColors = {
    positive: 'text-green-600 bg-green-50',
    negative: 'text-red-600 bg-red-50',
    neutral: 'text-blue-600 bg-blue-50',
  };

  const bgColor = trend ? trendColors[trend] : 'bg-gray-50';

  return (
    <div className="card hover:shadow-lg transition-shadow duration-200">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <p className="text-sm font-medium text-gray-600 mb-1">{title}</p>
          <p className="text-3xl font-bold text-gray-900 mb-2">{value}</p>
          {subtitle && (
            <p className={`text-sm font-medium ${trend ? trendColors[trend].split(' ')[0] : 'text-gray-500'}`}>
              {subtitle}
            </p>
          )}
        </div>
        {icon && (
          <div className={`text-3xl ${bgColor} p-3 rounded-lg`}>
            {icon}
          </div>
        )}
      </div>
    </div>
  );
}

export default function FinancialOverview() {
  const { data: summary, loading, error } = useApiData<CategorySummaryResponse>(
    getCategorySummary
  );

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {[1, 2, 3, 4].map((i) => (
          <div key={i} className="card animate-pulse">
            <div className="h-4 bg-gray-200 rounded w-1/2 mb-4"></div>
            <div className="h-8 bg-gray-200 rounded w-3/4 mb-2"></div>
            <div className="h-4 bg-gray-200 rounded w-1/3"></div>
          </div>
        ))}
      </div>
    );
  }

  if (error || !summary) {
    return (
      <div className="card bg-red-50 border border-red-200">
        <p className="text-red-800 font-medium">⚠️ {error || 'No data available'}</p>
      </div>
    );
  }

  const { summary: financialSummary, period } = summary;
  const isSavingsPositive = financialSummary.net_savings > 0;

  return (
    <div>
      {/* Period Badge */}
      <div className="mb-4">
        <span className="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800">
          📅 {period.start} to {period.end} ({period.months} months)
        </span>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {/* Total Income */}
        <StatCard
          title="Total Income"
          value={formatCurrency(financialSummary.total_income)}
          subtitle={`${formatCurrency(financialSummary.total_income / (period.months || 1))}/month avg`}
          trend="positive"
          icon="💰"
        />

        {/* Total Expenses */}
        <StatCard
          title="Total Expenses"
          value={formatCurrency(financialSummary.total_expenses)}
          subtitle={`${formatCurrency(financialSummary.total_expenses / (period.months || 1))}/month avg`}
          trend="negative"
          icon="💸"
        />

        {/* Net Savings */}
        <StatCard
          title="Net Savings"
          value={formatCurrency(financialSummary.net_savings)}
          subtitle={isSavingsPositive ? 'Great job saving!' : 'Spending exceeds income'}
          trend={isSavingsPositive ? 'positive' : 'negative'}
          icon="🏦"
        />

        {/* Savings Rate */}
        <StatCard
          title="Savings Rate"
          value={formatPercentage(financialSummary.savings_rate)}
          subtitle={
            financialSummary.savings_rate >= 20
              ? 'Excellent!'
              : financialSummary.savings_rate >= 10
              ? 'Good progress'
              : 'Room to improve'
          }
          trend={
            financialSummary.savings_rate >= 20
              ? 'positive'
              : financialSummary.savings_rate >= 10
              ? 'neutral'
              : 'negative'
          }
          icon="📊"
        />
      </div>
    </div>
  );
}

