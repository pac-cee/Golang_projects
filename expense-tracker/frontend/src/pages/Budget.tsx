import { useState } from 'react';
import {
  Grid,
  Paper,
  Typography,
  Button,
  LinearProgress,
  Card,
  CardContent,
  Box,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';

const Budget = () => {
  const [budgets] = useState([
    {
      id: '1',
      category: 'Food & Dining',
      amount: 500,
      spent: 350,
      period: 'monthly',
    },
    {
      id: '2',
      category: 'Transportation',
      amount: 300,
      spent: 200,
      period: 'monthly',
    },
  ]);

  const calculateProgress = (spent: number, total: number) => {
    return (spent / total) * 100;
  };

  return (
    <Grid container spacing={3}>
      <Grid item xs={12} display="flex" justifyContent="space-between" alignItems="center">
        <Typography variant="h4">Budget</Typography>
        <Button variant="contained" startIcon={<AddIcon />}>
          Set New Budget
        </Button>
      </Grid>

      {budgets.map((budget) => (
        <Grid item xs={12} md={6} lg={4} key={budget.id}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                {budget.category}
              </Typography>
              <Box sx={{ mb: 2 }}>
                <Typography variant="body2" color="text.secondary">
                  Monthly Budget: ${budget.amount}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Spent: ${budget.spent}
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Remaining: ${budget.amount - budget.spent}
                </Typography>
              </Box>
              <LinearProgress
                variant="determinate"
                value={calculateProgress(budget.spent, budget.amount)}
                color={calculateProgress(budget.spent, budget.amount) > 90 ? 'error' : 'primary'}
                sx={{ height: 10, borderRadius: 5 }}
              />
              <Box sx={{ mt: 2, display: 'flex', justifyContent: 'flex-end', gap: 1 }}>
                <Button size="small" variant="outlined">
                  Edit
                </Button>
                <Button size="small" variant="outlined" color="error">
                  Delete
                </Button>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      ))}
    </Grid>
  );
};

export default Budget;
