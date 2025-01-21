import React, { useState, useEffect } from 'react';
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
  Box,
  InputAdornment,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import dayjs, { Dayjs } from 'dayjs';
import { Budget, Category } from '../../types';
import { useAppContext } from '../../context/AppContext';

interface BudgetFormProps {
  open: boolean;
  onClose: () => void;
  budget?: Budget | null;
  onSubmit: (data: Omit<Budget, 'id' | 'createdAt' | 'updatedAt' | 'spent'>) => Promise<void>;
}

const initialFormData: Omit<Budget, 'id' | 'createdAt' | 'updatedAt' | 'spent'> = {
  amount: 0,
  category: '',
  startDate: '',
  endDate: '',
};

export const BudgetForm: React.FC<BudgetFormProps> = ({
  open,
  onClose,
  budget,
  onSubmit,
}) => {
  const { state } = useAppContext();
  const [formData, setFormData] = useState<Omit<Budget, 'id' | 'createdAt' | 'updatedAt' | 'spent'>>(
    initialFormData
  );

  useEffect(() => {
    if (budget) {
      setFormData({
        amount: budget.amount,
        category: budget.category,
        startDate: budget.startDate,
        endDate: budget.endDate,
      });
    } else {
      setFormData(initialFormData);
    }
  }, [budget]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onSubmit(formData);
    onClose();
  };

  const handleDateChange = (field: 'startDate' | 'endDate') => (date: Dayjs | null) => {
    setFormData({
      ...formData,
      [field]: date ? date.toISOString() : '',
    });
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>{budget ? 'Edit Budget' : 'Add Budget'}</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 2 }}>
            <FormControl fullWidth>
              <InputLabel>Category</InputLabel>
              <Select
                value={formData.category}
                onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                label="Category"
                required
              >
                {state.categories
                  .filter((category: Category) => category.type === 'expense')
                  .map((category: Category) => (
                    <MenuItem key={category.id} value={category.name}>
                      {category.name}
                    </MenuItem>
                  ))}
              </Select>
            </FormControl>

            <TextField
              label="Amount"
              type="number"
              value={formData.amount}
              onChange={(e) =>
                setFormData({ ...formData, amount: parseFloat(e.target.value) })
              }
              InputProps={{
                startAdornment: <InputAdornment position="start">$</InputAdornment>,
              }}
              required
            />

            <DatePicker
              label="Start Date"
              value={formData.startDate ? dayjs(formData.startDate) : null}
              onChange={handleDateChange('startDate')}
              slotProps={{
                textField: {
                  fullWidth: true,
                  required: true,
                },
              }}
            />

            <DatePicker
              label="End Date"
              value={formData.endDate ? dayjs(formData.endDate) : null}
              onChange={handleDateChange('endDate')}
              slotProps={{
                textField: {
                  fullWidth: true,
                  required: true,
                },
              }}
            />
          </Box>
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
