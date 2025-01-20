import React from 'react';
import { Bar } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Box, useTheme } from '@mui/material';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

interface ChartData {
  date: string;
  income: number;
  expenses: number;
}

interface IncomeExpenseChartProps {
  data: ChartData[];
}

const IncomeExpenseChart: React.FC<IncomeExpenseChartProps> = ({ data }) => {
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
        label: 'Income',
        data: data.map((item) => item.income),
        backgroundColor: theme.palette.success.main,
        borderColor: theme.palette.success.dark,
        borderWidth: 1,
      },
      {
        label: 'Expenses',
        data: data.map((item) => item.expenses),
        backgroundColor: theme.palette.error.main,
        borderColor: theme.palette.error.dark,
        borderWidth: 1,
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
        position: 'top' as const,
        align: 'end' as const,
        labels: {
          color: theme.palette.text.primary,
          font: {
            size: 12,
          },
          boxWidth: 12,
        },
      },
      tooltip: {
        callbacks: {
          label: (context: any) => {
            const label = context.dataset.label || '';
            const value = context.raw || 0;
            return `${label}: $${value.toFixed(2)}`;
          },
        },
      },
    },
  };

  return (
    <Box width="100%" height="100%">
      <Bar data={chartData} options={options} />
    </Box>
  );
};

export default IncomeExpenseChart;
