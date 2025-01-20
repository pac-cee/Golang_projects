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
import { Account } from '../../types';

interface AccountFormProps {
  open: boolean;
  onClose: () => void;
  account?: Account;
  onSubmit: (account: Omit<Account, 'id' | 'balance' | 'createdAt' | 'updatedAt'>) => Promise<void>;
}

const AccountForm: React.FC<AccountFormProps> = ({
  open,
  onClose,
  account,
  onSubmit,
}) => {
  const [formData, setFormData] = useState({
    name: '',
    type: 'bank' as 'bank' | 'cash' | 'mobile_money',
    initialBalance: '0',
  });
  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (account) {
      setFormData({
        name: account.name,
        type: account.type,
        initialBalance: account.balance.toString(),
      });
    } else {
      setFormData({
        name: '',
        type: 'bank',
        initialBalance: '0',
      });
    }
  }, [account]);

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) {
      newErrors.name = 'Please enter an account name';
    }
    if (
      !formData.initialBalance ||
      isNaN(Number(formData.initialBalance))
    ) {
      newErrors.initialBalance = 'Please enter a valid initial balance';
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
        name: formData.name.trim(),
        type: formData.type,
      });
      onClose();
    } catch (error) {
      console.error('Error submitting account:', error);
    }
  };

  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>{account ? 'Edit Account' : 'Add Account'}</DialogTitle>
      <form onSubmit={handleSubmit}>
        <DialogContent>
          <Grid container spacing={2}>
            <Grid item xs={12}>
              <TextField
                fullWidth
                label="Account Name"
                value={formData.name}
                onChange={(e) => {
                  setFormData({ ...formData, name: e.target.value });
                  if (errors.name) {
                    setErrors({ ...errors, name: '' });
                  }
                }}
                error={!!errors.name}
                helperText={errors.name}
              />
            </Grid>

            <Grid item xs={12}>
              <FormControl fullWidth>
                <InputLabel>Account Type</InputLabel>
                <Select
                  value={formData.type}
                  onChange={(e) =>
                    setFormData({
                      ...formData,
                      type: e.target.value as 'bank' | 'cash' | 'mobile_money',
                    })
                  }
                  label="Account Type"
                >
                  <MenuItem value="bank">Bank Account</MenuItem>
                  <MenuItem value="cash">Cash</MenuItem>
                  <MenuItem value="mobile_money">Mobile Money</MenuItem>
                </Select>
              </FormControl>
            </Grid>

            {!account && (
              <Grid item xs={12}>
                <TextField
                  fullWidth
                  label="Initial Balance"
                  type="number"
                  value={formData.initialBalance}
                  onChange={(e) => {
                    setFormData({ ...formData, initialBalance: e.target.value });
                    if (errors.initialBalance) {
                      setErrors({ ...errors, initialBalance: '' });
                    }
                  }}
                  error={!!errors.initialBalance}
                  helperText={errors.initialBalance}
                  InputProps={{
                    startAdornment: <InputAdornment position="start">$</InputAdornment>,
                  }}
                />
              </Grid>
            )}
          </Grid>
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

export default AccountForm;
