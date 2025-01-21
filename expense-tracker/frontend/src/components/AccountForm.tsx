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
import { Account, AccountType } from '../types';

interface AccountFormProps {
  open: boolean;
  onClose: () => void;
  account?: Account | null;
  onSubmit: (data: Omit<Account, 'id' | 'createdAt' | 'updatedAt'>) => Promise<void>;
}

const initialFormData: Omit<Account, 'id' | 'createdAt' | 'updatedAt'> = {
  name: '',
  type: 'checking',
  balance: 0,
  currency: 'USD',
};

export const AccountForm: React.FC<AccountFormProps> = ({
  open,
  onClose,
  account,
  onSubmit,
}) => {
  const [formData, setFormData] = useState<Omit<Account, 'id' | 'createdAt' | 'updatedAt'>>(initialFormData);

  useEffect(() => {
    if (account) {
      setFormData({
        name: account.name,
        type: account.type,
        balance: account.balance,
        currency: account.currency,
      });
    } else {
      setFormData(initialFormData);
    }
  }, [account]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await onSubmit(formData);
    onClose();
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <form onSubmit={handleSubmit}>
        <DialogTitle>{account ? 'Edit Account' : 'Add Account'}</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 2 }}>
            <TextField
              label="Name"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              required
            />

            <FormControl fullWidth>
              <InputLabel>Type</InputLabel>
              <Select
                value={formData.type}
                onChange={(e) => setFormData({ ...formData, type: e.target.value as AccountType })}
                label="Type"
                required
              >
                <MenuItem value="checking">Checking</MenuItem>
                <MenuItem value="savings">Savings</MenuItem>
                <MenuItem value="credit">Credit</MenuItem>
                <MenuItem value="cash">Cash</MenuItem>
                <MenuItem value="investment">Investment</MenuItem>
              </Select>
            </FormControl>

            <TextField
              label="Balance"
              type="number"
              value={formData.balance}
              onChange={(e) => setFormData({ ...formData, balance: parseFloat(e.target.value) })}
              InputProps={{
                startAdornment: <InputAdornment position="start">$</InputAdornment>,
              }}
              required
            />

            <FormControl fullWidth>
              <InputLabel>Currency</InputLabel>
              <Select
                value={formData.currency}
                onChange={(e) => setFormData({ ...formData, currency: e.target.value })}
                label="Currency"
                required
              >
                <MenuItem value="USD">USD</MenuItem>
                <MenuItem value="EUR">EUR</MenuItem>
                <MenuItem value="GBP">GBP</MenuItem>
              </Select>
            </FormControl>
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={onClose}>Cancel</Button>
          <Button type="submit" variant="contained" color="primary">
            {account ? 'Update' : 'Add'}
          </Button>
        </DialogActions>
      </form>
    </Dialog>
  );
};
