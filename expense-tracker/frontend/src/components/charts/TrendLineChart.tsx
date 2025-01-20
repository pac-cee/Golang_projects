import React from 'react';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js';
import { Box, useTheme } from '@mui/material';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

interface ChartData {
  date: string;
  income: number;
  expenses: number;
}

interface TrendLineChartProps {
  data: ChartData[];
}

const TrendLineChart: React.FC<TrendLineChartProps> = ({ data }) => {
  const theme = useTheme();

  const chartData = {
    labels: data.map((item) => {
      const [year, month] = item.date.split('-');
      return `${new Date(parseInt(year), parseInt(month) - 1).toLocaleString('default', {
        month: 'short',
      })} ${year}`;
    }),
    datasets: [
      {
        label: 'Net Balance',
        data: data.map((item) => item.income - item.expenses),
        fill: true,
        backgroundColor: (context: any) => {
          const ctx = context.chart.ctx;
          const gradient = ctx.createLinearGradient(0, 0, 0, 200);
          gradient.addColorStop(0, theme.palette.primary.light + '40');
          gradient.addColorStop(1, theme.palette.primary.light + '00');
          return gradient;
        },
        borderColor: theme.palette.primary.main,
        borderWidth: 2,
        tension: 0.4,
        pointRadius: 4,
        pointBackgroundColor: theme.palette.primary.main,
        pointBorderColor: theme.palette.background.paper,
        pointBorderWidth: 2,
        pointHoverRadius: 6,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      x: {
        grid: {
          display: false,
          color: theme.palette.divider,
        },
        ticks: {
          color: theme.palette.text.secondary,
        },
      },
      y: {
        grid: {
          color: theme.palette.divider,
        },
        ticks: {
          color: theme.palette.text.secondary,
          callback: (value: number) => `$${value.toFixed(0)}`,
        },
      },
    },
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        backgroundColor: theme.palette.background.paper,
        titleColor: theme.palette.text.primary,
        bodyColor: theme.palette.text.secondary,
        borderColor: theme.palette.divider,
        borderWidth: 1,
        padding: 12,
        callbacks: {
          label: (context: any) => {
            const value = context.raw || 0;
            return `Balance: $${value.toFixed(2)}`;
          },
        },
      },
    },
    interaction: {
      intersect: false,
      mode: 'index' as const,
    },
  };

  return (
    <Box width="100%" height="100%">
      <Line data={chartData} options={options} />
    </Box>
  );
};

export default TrendLineChart;
