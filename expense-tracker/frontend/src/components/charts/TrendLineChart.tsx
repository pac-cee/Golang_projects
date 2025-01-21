import React from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  ChartOptions,
  Scale,
  Filler,
} from 'chart.js';
import { Line } from 'react-chartjs-2';
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

interface TrendLineChartProps {
  data: {
    labels: string[];
    datasets: {
      label: string;
      data: number[];
      borderColor: string;
      backgroundColor: string;
      fill: boolean;
      tension: number;
    }[];
  };
}

export const TrendLineChart: React.FC<TrendLineChartProps> = ({ data }) => {
  const theme = useTheme();

  const options: ChartOptions<'line'> = {
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
        text: 'Spending Trend',
        color: theme.palette.text.primary,
        font: {
          size: 16,
        },
      },
    },
    interaction: {
      mode: 'index' as const,
      intersect: false,
    },
  };

  return (
    <Box width="100%" height="100%">
      <Line data={data} options={options} />
    </Box>
  );
};
