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
import { Account } from '../types';
import { useAppContext } from '../context/AppContext';
import { accountApi } from '../services/api';

interface AccountFormProps {
  open: boolean;
  onClose: () => void;
  account?: Account;
}

const AccountForm = ({ open, onClose, account }: AccountFormProps) => {
  const { state, dispatch } = useAppContext();
  const [formData, setFormData] = useState({
    name: '',
    type: 'bank',
    balance: '',
  });
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (account) {
      setFormData({
        name: account.name,
        type: account.type,
        balance: String(account.balance),
      });
    }
  }, [account]);

  const handleChange = (field: string) => (event: any) => {
    setFormData((prev) => ({
      ...prev,
      [field]: event.target.value,
    }));
  };

  const handleSubmit = async () => {
    if (!formData.name || !formData.balance) return;

    try {
      setLoading(true);
      const data = {
        ...formData,
        balance: Number(formData.balance),
      };

      if (account) {
        const response = await accountApi.update(account.id, data);
        dispatch({ type: 'SET_ACCOUNTS', payload: state.accounts.map((a) => 
          a.id === account.id ? response.data : a
        )});
      } else {
        const response = await accountApi.create(data);
        dispatch({ type: 'SET_ACCOUNTS', payload: [...state.accounts, response.data] });
      }
      onClose();
    } catch (error) {
      console.error('Error saving account:', error);
      dispatch({ type: 'SET_ERROR', payload: 'Error saving account' });
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    setFormData({
      name: '',
      type: 'bank',
      balance: '',
    });
    onClose();
  };

  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>
        {account ? 'Edit Account' : 'Add New Account'}
      </DialogTitle>
      <DialogContent>
        <Grid container spacing={2} sx={{ mt: 1 }}>
          <Grid item xs={12}>
            <TextField
              fullWidth
              label="Account Name"
              value={formData.name}
              onChange={handleChange('name')}
            />
          </Grid>
          <Grid item xs={12}>
            <FormControl fullWidth>
              <InputLabel>Account Type</InputLabel>
              <Select
                value={formData.type}
                label="Account Type"
                onChange={handleChange('type')}
              >
                <MenuItem value="bank">Bank Account</MenuItem>
                <MenuItem value="cash">Cash</MenuItem>
                <MenuItem value="mobile_money">Mobile Money</MenuItem>
              </Select>
            </FormControl>
          </Grid>
          <Grid item xs={12}>
            <TextField
              fullWidth
              label="Initial Balance"
              type="number"
              value={formData.balance}
              onChange={handleChange('balance')}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">$</InputAdornment>
                ),
              }}
            />
          </Grid>
        </Grid>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Cancel</Button>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={loading || !formData.name || !formData.balance}
        >
          {loading ? 'Saving...' : 'Save'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default AccountForm;
