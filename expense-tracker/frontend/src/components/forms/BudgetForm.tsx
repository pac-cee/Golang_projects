import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Grid,
  FormHelperText,
  InputAdornment,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { Budget } from '../../types';
import { useAppContext } from '../../context/AppContext';

interface BudgetFormProps {
  open: boolean;
  onClose: () => void;
  budget?: Budget;
  onSubmit: (budget: Omit<Budget, 'id' | 'spent' | 'createdAt' | 'updatedAt'>) => Promise<void>;
}

const BudgetForm: React.FC<BudgetFormProps> = ({
  open,
  onClose,
  budget,
  onSubmit,
}) => {
  const { state } = useAppContext();
  const [formData, setFormData] = useState({
    category: '',
    amount: '',
    period: 'monthly' as 'weekly' | 'monthly' | 'yearly',
    startDate: new Date(),
    endDate: new Date(),
  });
  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (budget) {
      setFormData({
        category: budget.category,
        amount: budget.amount.toString(),
        period: budget.period,
        startDate: new Date(budget.startDate),
        endDate: new Date(budget.endDate),
      });
    } else {
      // Set default end date based on period
      const endDate = new Date();
      endDate.setMonth(endDate.getMonth() + 1);
      setFormData({
        category: '',
        amount: '',
        period: 'monthly',
        startDate: new Date(),
        endDate,
      });
    }
  }, [budget]);

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.category) {
      newErrors.category = 'Please select a category';
    }
    if (!formData.amount || isNaN(Number(formData.amount)) || Number(formData.amount) <= 0) {
      newErrors.amount = 'Please enter a valid amount';
    }
    if (formData.startDate >= formData.endDate) {
      newErrors.endDate = 'End date must be after start date';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    try {
      await onSubmit({
        category: formData.category,
        amount: Number(formData.amount),
        period: formData.period,
        startDate: formData.startDate.toISOString(),
        endDate: formData.endDate.toISOString(),
      });
      onClose();
    } catch (error) {
      console.error('Error submitting budget:', error);
    }
  };

  const handlePeriodChange = (period: 'weekly' | 'monthly' | 'yearly') => {
    const startDate = new Date();
    const endDate = new Date();

    switch (period) {
      case 'weekly':
        endDate.setDate(endDate.getDate() + 7);
        break;
      case 'monthly':
        endDate.setMonth(endDate.getMonth() + 1);
        break;
      case 'yearly':
        endDate.setFullYear(endDate.getFullYear() + 1);
        break;
    }

    setFormData({
      ...formData,
      period,
      startDate,
      endDate,
    });
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>{budget ? 'Edit Budget' : 'Add Budget'}</DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <FormControl fullWidth error={!!errors.category}>
                <InputLabel>Category</InputLabel>
                <Select
                  value={formData.category}
                  onChange={(e) => {
                    setFormData({ ...formData, category: e.target.value });
                    if (errors.category) {
                      setErrors({ ...errors, category: '' });
                    }
                  }}
                  label="Category"
                >
                  {state.categories.map((category) => (
                    <MenuItem key={category.id} value={category.name}>
                      {category.name}
                    </MenuItem>
                  ))}
                </Select>
                {errors.category && <FormHelperText>{errors.category}</FormHelperText>}
              </FormControl>
            </Grid>

            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Amount"
                type="number"
                value={formData.amount}
                onChange={(e) => {
                  setFormData({ ...formData, amount: e.target.value });
                  if (errors.amount) {
                    setErrors({ ...errors, amount: '' });
                  }
                }}
                error={!!errors.amount}
                helperText={errors.amount}
                InputProps={{
                  startAdornment: <InputAdornment position="start">$</InputAdornment>,
                }}
              />
            </Grid>

            <Grid item xs={12}>
              <FormControl fullWidth>
                <InputLabel>Period</InputLabel>
                <Select
                  value={formData.period}
                  onChange={(e) => handlePeriodChange(e.target.value as 'weekly' | 'monthly' | 'yearly')}
                  label="Period"
                >
                  <MenuItem value="weekly">Weekly</MenuItem>
                  <MenuItem value="monthly">Monthly</MenuItem>
                  <MenuItem value="yearly">Yearly</MenuItem>
                </Select>
              </FormControl>
            </Grid>

            <Grid item xs={12} sm={6}>
              <DatePicker
                label="Start Date"
                value={formData.startDate}
                onChange={(newValue) =>
                  setFormData({ ...formData, startDate: newValue || new Date() })
                }
                slotProps={{
                  textField: {
                    fullWidth: true,
                  },
                }}
              />
            </Grid>

            <Grid item xs={12} sm={6}>
              <DatePicker
                label="End Date"
                value={formData.endDate}
                onChange={(newValue) => {
                  setFormData({ ...formData, endDate: newValue || new Date() });
                  if (errors.endDate) {
                    setErrors({ ...errors, endDate: '' });
                  }
                }}
                slotProps={{
                  textField: {
                    fullWidth: true,
                    error: !!errors.endDate,
                    helperText: errors.endDate,
                  },
                }}
              />
            </Grid>
          </Grid>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Cancel</Button>
          <Button type="submit" variant="contained" color="primary">
            {budget ? 'Update' : 'Add'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};

export default BudgetForm;
