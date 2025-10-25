import { useState } from 'react';
import {
  LineChart,
  Line,
  AreaChart,
  Area,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { getTimeline, TimelineResponse, TimelinePoint } from '../../services/api';
import { formatCurrency, formatMonthYear } from '../../utils/formatters';
import { useApiData } from '../../hooks/useFetch';
import { ToggleButtonGroup, LoadingSkeleton, ErrorAlert } from '../common';

interface CustomTooltipProps {
  active?: boolean;
  payload?: Array<{
    name: string;
    value: number;
    dataKey: string;
    color: string;
  }>;
  label?: string;
}

function CustomTooltip({ active, payload, label }: CustomTooltipProps) {
  if (active && payload && payload.length > 0) {
    return (
      <div className="bg-white p-4 rounded-lg shadow-lg border border-gray-200">
        <p className="font-semibold text-gray-900 mb-2">{label && formatMonthYear(label)}</p>
        {payload.map((entry, index) => (
          <p key={index} className="text-sm" style={{ color: entry.color }}>
            {entry.name}: {formatCurrency(entry.value)}
          </p>
        ))}
      </div>
    );
  }
  return null;
}

type ChartType = 'line' | 'area';

export default function TimelineChart() {
  const { data: timeline, loading, error } = useApiData<TimelineResponse>(getTimeline);
  const [chartType, setChartType] = useState<ChartType>('area');

  if (loading) {
    return <LoadingSkeleton variant="chart" />;
  }

  if (error || !timeline) {
    return <ErrorAlert message={error || 'No data available'} />;
  }

  // Format data for display
  const formattedData = timeline.timeline.map((point: TimelinePoint) => ({
    ...point,
    displayPeriod: formatMonthYear(point.period),
  }));

  // Calculate stats
  const avgIncome =
    timeline.timeline.reduce((sum, p) => sum + p.income, 0) / timeline.timeline.length;
  const avgExpenses =
    timeline.timeline.reduce((sum, p) => sum + p.expenses, 0) / timeline.timeline.length;
  const totalNet = timeline.timeline.reduce((sum, p) => sum + p.net, 0);

  return (
    <div className="card">
      {/* Header with Controls */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between mb-6 gap-4">
        <h2 className="text-xl font-semibold text-gray-900">
          📈 Income vs Expenses Timeline
        </h2>

        {/* Chart Type Toggle */}
        <ToggleButtonGroup
          options={[
            { value: 'area' as const, label: 'Area Chart' },
            { value: 'line' as const, label: 'Line Chart' },
          ]}
          value={chartType}
          onChange={setChartType}
          className="w-fit"
        />
      </div>

      {/* Chart */}
      <div className="h-80">
        <ResponsiveContainer width="100%" height="100%">
          {chartType === 'area' ? (
            <AreaChart
              data={formattedData}
              margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
            >
              <defs>
                <linearGradient id="colorIncome" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#10b981" stopOpacity={0.8} />
                  <stop offset="95%" stopColor="#10b981" stopOpacity={0} />
                </linearGradient>
                <linearGradient id="colorExpenses" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor="#ef4444" stopOpacity={0.8} />
                  <stop offset="95%" stopColor="#ef4444" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
              <XAxis
                dataKey="displayPeriod"
                stroke="#6b7280"
                style={{ fontSize: '12px' }}
              />
              <YAxis
                stroke="#6b7280"
                style={{ fontSize: '12px' }}
                tickFormatter={(value) => `$${(value / 1000).toFixed(0)}k`}
              />
              <Tooltip content={<CustomTooltip />} />
              <Legend
                wrapperStyle={{ fontSize: '14px' }}
                formatter={(value) => <span className="font-medium">{value}</span>}
              />
              <Area
                type="monotone"
                dataKey="income"
                name="Income"
                stroke="#10b981"
                strokeWidth={2}
                fillOpacity={1}
                fill="url(#colorIncome)"
              />
              <Area
                type="monotone"
                dataKey="expenses"
                name="Expenses"
                stroke="#ef4444"
                strokeWidth={2}
                fillOpacity={1}
                fill="url(#colorExpenses)"
              />
            </AreaChart>
          ) : (
            <LineChart
              data={formattedData}
              margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
            >
              <CartesianGrid strokeDasharray="3 3" stroke="#e5e7eb" />
              <XAxis
                dataKey="displayPeriod"
                stroke="#6b7280"
                style={{ fontSize: '12px' }}
              />
              <YAxis
                stroke="#6b7280"
                style={{ fontSize: '12px' }}
                tickFormatter={(value) => `$${(value / 1000).toFixed(0)}k`}
              />
              <Tooltip content={<CustomTooltip />} />
              <Legend
                wrapperStyle={{ fontSize: '14px' }}
                formatter={(value) => <span className="font-medium">{value}</span>}
              />
              <Line
                type="monotone"
                dataKey="income"
                name="Income"
                stroke="#10b981"
                strokeWidth={3}
                dot={{ fill: '#10b981', strokeWidth: 2, r: 4 }}
                activeDot={{ r: 6 }}
              />
              <Line
                type="monotone"
                dataKey="expenses"
                name="Expenses"
                stroke="#ef4444"
                strokeWidth={3}
                dot={{ fill: '#ef4444', strokeWidth: 2, r: 4 }}
                activeDot={{ r: 6 }}
              />
              <Line
                type="monotone"
                dataKey="net"
                name="Net"
                stroke="#3b82f6"
                strokeWidth={2}
                strokeDasharray="5 5"
                dot={{ fill: '#3b82f6', strokeWidth: 2, r: 3 }}
              />
            </LineChart>
          )}
        </ResponsiveContainer>
      </div>

      {/* Summary Stats */}
      <div className="mt-6 pt-6 border-t grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="text-center">
          <p className="text-sm text-gray-600 mb-1">Average Income</p>
          <p className="text-2xl font-bold text-green-600">{formatCurrency(avgIncome)}</p>
        </div>
        <div className="text-center">
          <p className="text-sm text-gray-600 mb-1">Average Expenses</p>
          <p className="text-2xl font-bold text-red-600">{formatCurrency(avgExpenses)}</p>
        </div>
        <div className="text-center">
          <p className="text-sm text-gray-600 mb-1">Total Net</p>
          <p
            className={`text-2xl font-bold ${
              totalNet >= 0 ? 'text-blue-600' : 'text-orange-600'
            }`}
          >
            {formatCurrency(totalNet)}
          </p>
        </div>
      </div>
    </div>
  );
}

