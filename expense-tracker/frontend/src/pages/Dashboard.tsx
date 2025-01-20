import React, { useState, useEffect } from 'react';
import {
  Grid,
  Paper,
  Typography,
  Box,
  useTheme,
} from '@mui/material';
import { styled } from '@mui/material/styles';

import { useAppContext } from '../context/AppContext';
import { reportApi } from '../services/api';
import ExpensePieChart from '../components/charts/ExpensePieChart';
import IncomeExpenseChart from '../components/charts/IncomeExpenseChart';
import TrendLineChart from '../components/charts/TrendLineChart';

const StyledPaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(3),
  height: '100%',
}));

const Dashboard: React.FC = () => {
  const theme = useTheme();
  const { state } = useAppContext();
  const [totalIncome, setTotalIncome] = useState<number>(0);
  const [totalExpenses, setTotalExpenses] = useState<number>(0);
  const [netAmount, setNetAmount] = useState<number>(0);
  const [categoryExpenses, setCategoryExpenses] = useState<any[]>([]);
  const [monthlyData, setMonthlyData] = useState<any[]>([]);

  useEffect(() => {
    const fetchDashboardData = async () => {
      try {
        const startDate = new Date();
        startDate.setMonth(startDate.getMonth() - 1);
        
        const [transactionReport, categoryReport] = await Promise.all([
          reportApi.getTransactionReport({ startDate: startDate.toISOString() }),
          reportApi.getCategoryReport({ startDate: startDate.toISOString() })
        ]);

        setTotalIncome(transactionReport.totalIncome);
        setTotalExpenses(transactionReport.totalExpenses);
        setNetAmount(transactionReport.netAmount);
        setCategoryExpenses(categoryReport);

        // Process monthly data for charts
        const monthlyTransactions = transactionReport.transactions.reduce((acc: any, transaction: any) => {
          const date = new Date(transaction.date);
          const monthYear = `${date.getFullYear()}-${date.getMonth() + 1}`;
          
          if (!acc[monthYear]) {
            acc[monthYear] = {
              income: 0,
              expenses: 0,
            };
          }

          if (transaction.amount >= 0) {
            acc[monthYear].income += transaction.amount;
          } else {
            acc[monthYear].expenses += Math.abs(transaction.amount);
          }

          return acc;
        }, {});

        const chartData = Object.entries(monthlyTransactions).map(([date, data]: [string, any]) => ({
          date,
          ...data,
        }));

        setMonthlyData(chartData);
      } catch (error) {
        console.error('Error fetching dashboard data:', error);
      }
    };

    fetchDashboardData();
  }, []);

  return (
    <Grid container spacing={3}>
      <Grid item xs={12}>
        <Typography variant="h4" gutterBottom>
          Dashboard
        </Typography>
      </Grid>

      {/* Summary Cards */}
      <Grid item xs={12} md={4}>
        <StyledPaper>
          <Typography variant="h6" gutterBottom>
            Total Income
          </Typography>
          <Typography variant="h4" color="success.main">
            ${totalIncome.toFixed(2)}
          </Typography>
        </StyledPaper>
      </Grid>

      <Grid item xs={12} md={4}>
        <StyledPaper>
          <Typography variant="h6" gutterBottom>
            Total Expenses
          </Typography>
          <Typography variant="h4" color="error.main">
            ${totalExpenses.toFixed(2)}
          </Typography>
        </StyledPaper>
      </Grid>

      <Grid item xs={12} md={4}>
        <StyledPaper>
          <Typography variant="h6" gutterBottom>
            Net Amount
          </Typography>
          <Typography
            variant="h4"
            color={netAmount >= 0 ? 'success.main' : 'error.main'}
          >
            ${netAmount.toFixed(2)}
          </Typography>
        </StyledPaper>
      </Grid>

      {/* Charts */}
      <Grid item xs={12} md={6}>
        <StyledPaper>
          <Typography variant="h6" gutterBottom>
            Expense Distribution
          </Typography>
          <Box sx={{ height: 300 }}>
            <ExpensePieChart data={categoryExpenses} />
          </Box>
        </StyledPaper>
      </Grid>

      <Grid item xs={12} md={6}>
        <StyledPaper>
          <Typography variant="h6" gutterBottom>
            Income vs Expenses
          </Typography>
          <Box sx={{ height: 300 }}>
            <IncomeExpenseChart data={monthlyData} />
          </Box>
        </StyledPaper>
      </Grid>

      <Grid item xs={12}>
        <StyledPaper>
          <Typography variant="h6" gutterBottom>
            Balance Trend
          </Typography>
          <Box sx={{ height: 300 }}>
            <TrendLineChart data={monthlyData} />
          </Box>
        </StyledPaper>
      </Grid>
    </Grid>
  );
};

export default Dashboard;
