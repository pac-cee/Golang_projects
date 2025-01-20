import { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Grid,
  InputAdornment,
} from '@mui/material';
import { Budget } from '../types';
import { useAppContext } from '../context/AppContext';
import { budgetApi } from '../services/api';

interface BudgetFormProps {
  open: boolean;
  onClose: () => void;
  budget?: Budget;
}

const BudgetForm = ({ open, onClose, budget }: BudgetFormProps) => {
  const { state, dispatch } = useAppContext();
  const [formData, setFormData] = useState({
    category: '',
    amount: '',
    period: 'monthly',
  });
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (budget) {
      setFormData({
        category: budget.category,
        amount: String(budget.amount),
        period: budget.period,
      });
    }
  }, [budget]);

  const handleChange = (field: string) => (event: any) => {
    setFormData((prev) => ({
      ...prev,
      [field]: event.target.value,
    }));
  };

  const handleSubmit = async () => {
    if (!formData.category || !formData.amount) return;

    try {
      setLoading(true);
      const data = {
        ...formData,
        amount: Number(formData.amount),
        spent: budget?.spent || 0,
      };

      if (budget) {
        const response = await budgetApi.update(budget.id, data);
        dispatch({ type: 'UPDATE_BUDGET', payload: response.data });
      } else {
        const response = await budgetApi.create(data);
        dispatch({ type: 'SET_BUDGETS', payload: [...state.budgets, response.data] });
      }
      onClose();
    } catch (error) {
      console.error('Error saving budget:', error);
      dispatch({ type: 'SET_ERROR', payload: 'Error saving budget' });
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setFormData({
      category: '',
      amount: '',
      period: 'monthly',
    });
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        {budget ? 'Edit Budget' : 'Set New Budget'}
      </DialogTitle>
      <DialogContent>
        <Grid container spacing={2} sx={{ mt: 1 }}>
          <Grid item xs={12}>
            <FormControl fullWidth>
              <InputLabel>Category</InputLabel>
              <Select
                value={formData.category}
                label="Category"
                onChange={handleChange('category')}
              >
                {state.categories.map((category) => (
                  <MenuItem key={category.id} value={category.id}>
                    {category.name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          <Grid item xs={12}>
            <TextField
              fullWidth
              label="Budget Amount"
              type="number"
              value={formData.amount}
              onChange={handleChange('amount')}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">$</InputAdornment>
                ),
              }}
            />
          </Grid>
          <Grid item xs={12}>
            <FormControl fullWidth>
              <InputLabel>Period</InputLabel>
              <Select
                value={formData.period}
                label="Period"
                onChange={handleChange('period')}
              >
                <MenuItem value="weekly">Weekly</MenuItem>
                <MenuItem value="monthly">Monthly</MenuItem>
                <MenuItem value="yearly">Yearly</MenuItem>
              </Select>
            </FormControl>
          </Grid>
        </Grid>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Cancel</Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={loading || !formData.category || !formData.amount}
        >
          {loading ? 'Saving...' : 'Save'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default BudgetForm;
