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
import { Transaction, TransactionType, Category, Account } from '../../types';
import { useAppContext } from '../../context/AppContext';

interface TransactionFormProps {
  open: boolean;
  onClose: () => void;
  transaction?: Transaction | null;
  onSubmit: (data: Omit<Transaction, 'id' | 'createdAt' | 'updatedAt'>) => Promise<void>;
}

const initialFormData: Omit<Transaction, 'id' | 'createdAt' | 'updatedAt'> = {
  amount: 0,
  type: 'expense',
  category: '',
  subcategory: '',
  account: '',
  description: '',
  date: dayjs().toISOString(),
};

export const TransactionForm: React.FC<TransactionFormProps> = ({
  open,
  onClose,
  transaction,
  onSubmit,
}) => {
  const { state } = useAppContext();
  const [formData, setFormData] = useState<Omit<Transaction, 'id' | 'createdAt' | 'updatedAt'>>(
    initialFormData
  );

  useEffect(() => {
    if (transaction) {
      setFormData({
        amount: transaction.amount,
        type: transaction.type,
        category: transaction.category,
        subcategory: transaction.subcategory,
        account: transaction.account,
        description: transaction.description,
        date: transaction.date,
      });
    } else {
      setFormData(initialFormData);
    }
  }, [transaction]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onSubmit(formData);
    onClose();
  };

  const selectedCategory = state.categories.find(
    (category: Category) => category.name === formData.category
  );

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>{transaction ? 'Edit Transaction' : 'Add Transaction'}</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 2 }}>
            <FormControl fullWidth>
              <InputLabel>Type</InputLabel>
              <Select
                value={formData.type}
                onChange={(e) => setFormData({ ...formData, type: e.target.value as TransactionType })}
                label="Type"
                required
              >
                <MenuItem value="expense">Expense</MenuItem>
                <MenuItem value="income">Income</MenuItem>
              </Select>
            </FormControl>

            <FormControl fullWidth>
              <InputLabel>Category</InputLabel>
              <Select
                value={formData.category}
                onChange={(e) => {
                  setFormData({
                    ...formData,
                    category: e.target.value,
                    subcategory: '', // Reset subcategory when category changes
                  });
                }}
                label="Category"
                required
              >
                {state.categories
                  .filter((category: Category) => category.type === formData.type)
                  .map((category: Category) => (
                    <MenuItem key={category.id} value={category.name}>
                      {category.name}
                    </MenuItem>
                  ))}
              </Select>
            </FormControl>

            {selectedCategory && selectedCategory.subcategories.length > 0 && (
              <FormControl fullWidth>
                <InputLabel>Subcategory</InputLabel>
                <Select
                  value={formData.subcategory}
                  onChange={(e) =>
                    setFormData({ ...formData, subcategory: e.target.value })
                  }
                  label="Subcategory"
                >
                  <MenuItem value="">None</MenuItem>
                  {selectedCategory.subcategories.map((subcategory: string) => (
                    <MenuItem key={subcategory} value={subcategory}>
                      {subcategory}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            )}

            <FormControl fullWidth>
              <InputLabel>Account</InputLabel>
              <Select
                value={formData.account}
                onChange={(e) => setFormData({ ...formData, account: e.target.value })}
                label="Account"
                required
              >
                {state.accounts.map((account: Account) => (
                  <MenuItem key={account.id} value={account.name}>
                    {account.name}
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

            <TextField
              label="Description"
              value={formData.description}
              onChange={(e) =>
                setFormData({ ...formData, description: e.target.value })
              }
              multiline
              rows={2}
            />

            <DatePicker
              label="Date"
              value={formData.date ? dayjs(formData.date) : null}
              onChange={(date: Dayjs | null) =>
                setFormData({
                  ...formData,
                  date: date ? date.toISOString() : dayjs().toISOString(),
                })
              }
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
            {transaction ? 'Update' : 'Add'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};
