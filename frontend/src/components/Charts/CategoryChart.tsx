import { useState } from 'react';
import { PieChart, Pie, Cell, ResponsiveContainer, Legend, Tooltip, PieLabelRenderProps } from 'recharts';
import { getCategorySummary, CategorySummaryResponse, CategoryDetail } from '../../services/api';
import { formatCurrency, formatPercentage } from '../../utils/formatters';
import { useApiData } from '../../hooks/useFetch';

interface ChartDataItem {
  name: string;
  value: number;
  percentage: number;
  count: number;
  [key: string]: string | number; // Index signature for Recharts compatibility
}

// Color palette for categories
const COLORS = [
  '#3b82f6', // blue-500
  '#10b981', // green-500
  '#f59e0b', // amber-500
  '#ef4444', // red-500
  '#8b5cf6', // violet-500
  '#ec4899', // pink-500
  '#06b6d4', // cyan-500
  '#f97316', // orange-500
  '#6366f1', // indigo-500
  '#14b8a6', // teal-500
];

interface CustomTooltipProps {
  active?: boolean;
  payload?: Array<{
    name: string;
    value: number;
    payload: ChartDataItem;
  }>;
}

function CustomTooltip({ active, payload }: CustomTooltipProps) {
  if (active && payload && payload.length > 0 && payload[0]) {
    const data = payload[0].payload;
    return (
      <div className="bg-white p-4 rounded-lg shadow-lg border border-gray-200">
        <p className="font-semibold text-gray-900 mb-2">{data.name}</p>
        <p className="text-sm text-gray-600">Amount: {formatCurrency(data.value)}</p>
        <p className="text-sm text-gray-600">Percentage: {formatPercentage(data.percentage)}</p>
        <p className="text-sm text-gray-600">Transactions: {data.count}</p>
      </div>
    );
  }
  return null;
}

export default function CategoryChart() {
  const { data: summary, loading, error } = useApiData<CategorySummaryResponse>(
    getCategorySummary
  );
  const [selectedType, setSelectedType] = useState<'expenses' | 'income'>('expenses');

  if (loading) {
    return (
      <div className="card">
        <div className="animate-pulse">
          <div className="h-6 bg-gray-200 rounded w-1/3 mb-4"></div>
          <div className="h-64 bg-gray-200 rounded"></div>
        </div>
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

  // Prepare chart data
  const categoryData: Record<string, CategoryDetail> =
    selectedType === 'expenses' ? summary.expenses : summary.income;

  const chartData: ChartDataItem[] = Object.entries(categoryData).map(([name, detail]) => ({
    name: name.charAt(0).toUpperCase() + name.slice(1),
    value: detail.total,
    percentage: detail.percentage,
    count: detail.count,
  }));

  // Sort by value descending
  chartData.sort((a, b) => b.value - a.value);

  return (
    <div className="card">
      {/* Header with Toggle */}
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-semibold text-gray-900">
          📊 Category Breakdown
        </h2>
        
        {/* Toggle Buttons */}
        <div className="inline-flex rounded-lg border border-gray-300 bg-white">
          <button
            onClick={() => setSelectedType('expenses')}
            className={`px-4 py-2 text-sm font-medium rounded-l-lg transition-colors ${
              selectedType === 'expenses'
                ? 'bg-blue-600 text-white'
                : 'text-gray-700 hover:bg-gray-50'
            }`}
          >
            Expenses
          </button>
          <button
            onClick={() => setSelectedType('income')}
            className={`px-4 py-2 text-sm font-medium rounded-r-lg transition-colors ${
              selectedType === 'income'
                ? 'bg-blue-600 text-white'
                : 'text-gray-700 hover:bg-gray-50'
            }`}
          >
            Income
          </button>
        </div>
      </div>

      {/* Chart */}
      <div className="h-80">
        <ResponsiveContainer width="100%" height="100%">
          <PieChart>
            <Pie
              data={chartData}
              cx="50%"
              cy="50%"
              labelLine={false}
              label={(props: PieLabelRenderProps) => {
                const entry = chartData[props.index ?? 0];
                return entry ? `${entry.name} ${entry.percentage.toFixed(1)}%` : '';
              }}
              outerRadius={100}
              innerRadius={60}
              fill="#8884d8"
              dataKey="value"
              paddingAngle={2}
            >
              {chartData.map((_, index) => (
                <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
              ))}
            </Pie>
            <Tooltip content={<CustomTooltip />} />
            <Legend
              verticalAlign="bottom"
              height={36}
              formatter={(value) => <span className="text-sm text-gray-700">{value}</span>}
            />
          </PieChart>
        </ResponsiveContainer>
      </div>

      {/* Summary Table */}
      <div className="mt-6 border-t pt-4">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {chartData.slice(0, 6).map((item, index) => (
            <div key={item.name} className="flex items-center gap-3">
              <div
                className="w-4 h-4 rounded-full shrink-0"
                style={{ backgroundColor: COLORS[index % COLORS.length] }}
              />
              <div className="flex-1 min-w-0">
                <p className="text-sm font-medium text-gray-900 truncate">{item.name}</p>
                <p className="text-xs text-gray-500">
                  {formatCurrency(item.value)} ({formatPercentage(item.percentage)})
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

