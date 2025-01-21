import React from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
  ChartOptions,
  Scale,
} from 'chart.js';
import { Bar } from 'react-chartjs-2';
import { Box, useTheme } from '@mui/material';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

interface IncomeExpenseChartProps {
  data: {
    labels: string[];
    datasets: {
      label: string;
      data: number[];
      backgroundColor: string[];
      borderColor: string[];
      borderWidth: number;
    }[];
  };
}

export const IncomeExpenseChart: React.FC<IncomeExpenseChartProps> = ({ data }) => {
  const theme = useTheme();

  const options: ChartOptions<'bar'> = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      y: {
        beginAtZero: true,
        grid: {
          color: theme.palette.divider,
        },
        ticks: {
          color: theme.palette.text.secondary,
          callback: function(this: Scale, value: string | number) {
            return `$${typeof value === 'number' ? value.toFixed(0) : value}`;
          },
        },
      },
      x: {
        grid: {
          display: false,
          color: theme.palette.divider,
        },
        ticks: {
          color: theme.palette.text.secondary,
        },
      },
    },
    plugins: {
      legend: {
        position: 'top' as const,
        labels: {
          color: theme.palette.text.secondary,
        },
      },
      title: {
        display: true,
        text: 'Income vs Expenses',
        color: theme.palette.text.primary,
        font: {
          size: 16,
        },
      },
    },
  };

  return (
    <Box width="100%" height="100%">
      <Bar data={data} options={options} />
    </Box>
  );
};
