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
import { Transaction } from '../types';
import { useAppContext } from '../context/AppContext';
import { transactionApi } from '../services/api';

interface TransactionFormProps {
  open: boolean;
  onClose: () => void;
  transaction?: Transaction;
}

const TransactionForm = ({ open, onClose, transaction }: TransactionFormProps) => {
  const { state, dispatch } = useAppContext();
  const [formData, setFormData] = useState({
    amount: '',
    type: 'expense',
    category: '',
    subcategory: '',
    account: '',
    description: '',
    date: new Date().toISOString().split('T')[0],
  });
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (transaction) {
      setFormData({
        amount: String(Math.abs(transaction.amount)),
        type: transaction.amount < 0 ? 'expense' : 'income',
        category: transaction.category,
        subcategory: transaction.subcategory || '',
        account: transaction.account,
        description: transaction.description,
        date: new Date(transaction.date).toISOString().split('T')[0],
      });
    }
  }, [transaction]);

  const handleSubmit = async () => {
    try {
      setLoading(true);
      const amount = formData.type === 'expense' 
        ? -Math.abs(Number(formData.amount))
        : Math.abs(Number(formData.amount));

      const data = {
        ...formData,
        amount,
      };

      if (transaction) {
        const response = await transactionApi.update(transaction.id, data);
        dispatch({ type: 'UPDATE_TRANSACTION', payload: response.data });
      } else {
        const response = await transactionApi.create(data);
        dispatch({ type: 'ADD_TRANSACTION', payload: response.data });
      }
      onClose();
    } catch (error) {
      console.error('Error saving transaction:', error);
      dispatch({ type: 'SET_ERROR', payload: 'Error saving transaction' });
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (field: string) => (event: any) => {
    setFormData((prev) => ({
      ...prev,
      [field]: event.target.value,
    }));
  };

  const selectedCategory = state.categories.find((c) => c.id === formData.category);

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        {transaction ? 'Edit Transaction' : 'Add New Transaction'}
      </DialogTitle>
      <DialogContent>
        <Grid container spacing={2} sx={{ mt: 1 }}>
          <Grid item xs={12} sm={6}>
            <TextField
              fullWidth
              label="Amount"
              type="number"
              value={formData.amount}
              onChange={handleChange('amount')}
              InputProps={{
                startAdornment: <InputAdornment position="start">$</InputAdornment>,
              }}
            />
          </Grid>
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth>
              <InputLabel>Type</InputLabel>
              <Select
                value={formData.type}
                label="Type"
                onChange={handleChange('type')}
              >
                <MenuItem value="expense">Expense</MenuItem>
                <MenuItem value="income">Income</MenuItem>
              </Select>
            </FormControl>
          </Grid>
          <Grid item xs={12} sm={6}>
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
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth>
              <InputLabel>Subcategory</InputLabel>
              <Select
                value={formData.subcategory}
                label="Subcategory"
                onChange={handleChange('subcategory')}
                disabled={!selectedCategory}
              >
                {selectedCategory?.subcategories.map((sub) => (
                  <MenuItem key={sub} value={sub}>
                    {sub}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          <Grid item xs={12} sm={6}>
            <FormControl fullWidth>
              <InputLabel>Account</InputLabel>
              <Select
                value={formData.account}
                label="Account"
                onChange={handleChange('account')}
              >
                {state.accounts.map((account) => (
                  <MenuItem key={account.id} value={account.id}>
                    {account.name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </Grid>
          <Grid item xs={12} sm={6}>
            <TextField
              fullWidth
              type="date"
              label="Date"
              value={formData.date}
              onChange={handleChange('date')}
              InputLabelProps={{ shrink: true }}
            />
          </Grid>
          <Grid item xs={12}>
            <TextField
              fullWidth
              label="Description"
              multiline
              rows={2}
              value={formData.description}
              onChange={handleChange('description')}
            />
          </Grid>
        </Grid>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={loading}
        >
          {loading ? 'Saving...' : 'Save'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default TransactionForm;
