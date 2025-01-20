import { Grid, Paper, Typography, Box, Select, MenuItem, FormControl, InputLabel } from '@mui/material';
import { useState } from 'react';

const Reports = () => {
  const [timeRange, setTimeRange] = useState('month');

  return (
    <Grid container spacing={3}>
      <Grid item xs={12} display="flex" justifyContent="space-between" alignItems="center">
        <Typography variant="h4">Reports</Typography>
        <FormControl sx={{ minWidth: 200 }}>
          <InputLabel>Time Range</InputLabel>
          <Select
            value={timeRange}
            label="Time Range"
            onChange={(e) => setTimeRange(e.target.value)}
          >
            <MenuItem value="week">Last Week</MenuItem>
            <MenuItem value="month">Last Month</MenuItem>
            <MenuItem value="quarter">Last Quarter</MenuItem>
            <MenuItem value="year">Last Year</MenuItem>
          </Select>
        </FormControl>
      </Grid>

      {/* Income vs Expenses Chart */}
      <Grid item xs={12} md={6}>
        <Paper sx={{ p: 2 }}>
          <Typography variant="h6" gutterBottom>
            Income vs Expenses
          </Typography>
          <Box sx={{ height: 300 }}>
            {/* Chart component will be added here */}
          </Box>
        </Paper>
      </Grid>

      {/* Expense by Category Chart */}
      <Grid item xs={12} md={6}>
        <Paper sx={{ p: 2 }}>
          <Typography variant="h6" gutterBottom>
            Expenses by Category
          </Typography>
          <Box sx={{ height: 300 }}>
            {/* Chart component will be added here */}
          </Box>
        </Paper>
      </Grid>

      {/* Monthly Trend Chart */}
      <Grid item xs={12}>
        <Paper sx={{ p: 2 }}>
          <Typography variant="h6" gutterBottom>
            Monthly Trend
          </Typography>
          <Box sx={{ height: 300 }}>
            {/* Chart component will be added here */}
          </Box>
        </Paper>
      </Grid>
    </Grid>
  );
};

export default Reports;
