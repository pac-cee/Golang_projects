import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  Grid,
  IconButton,
  Typography,
  TextField,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Chip,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import { Account } from '../types';
import { useAppContext } from '../context/AppContext';
import { accountApi } from '../services/api';
import AccountForm from '../components/forms/AccountForm';

const Accounts: React.FC = () => {
  const { state, dispatch } = useAppContext();
  const [openForm, setOpenForm] = useState(false);
  const [selectedAccount, setSelectedAccount] = useState<Account | undefined>();
  const [search, setSearch] = useState('');

  useEffect(() => {
    const fetchAccounts = async () => {
      try {
        dispatch({ type: 'SET_LOADING', payload: true });
        const accounts = await accountApi.getAll();
        dispatch({ type: 'SET_ACCOUNTS', payload: accounts });
      } catch (error) {
        dispatch({ type: 'SET_ERROR', payload: 'Error fetching accounts' });
      } finally {
        dispatch({ type: 'SET_LOADING', payload: false });
      }
    };

    fetchAccounts();
  }, [dispatch]);

  const handleAddAccount = async (
    account: Omit<Account, 'id' | 'createdAt' | 'updatedAt'>
  ) => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true });
      const newAccount = await accountApi.create(account);
      dispatch({ type: 'ADD_ACCOUNT', payload: newAccount });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Error adding account' });
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  };

  const handleUpdateAccount = async (
    account: Omit<Account, 'createdAt' | 'updatedAt'>
  ) => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true });
      const updatedAccount = await accountApi.update(account.id, account);
      dispatch({ type: 'UPDATE_ACCOUNT', payload: updatedAccount });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Error updating account' });
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  };

  const handleDeleteAccount = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this account?')) {
      try {
        dispatch({ type: 'SET_LOADING', payload: true });
        await accountApi.delete(id);
        dispatch({ type: 'DELETE_ACCOUNT', payload: id });
      } catch (error) {
        dispatch({ type: 'SET_ERROR', payload: 'Error deleting account' });
      } finally {
        dispatch({ type: 'SET_LOADING', payload: false });
      }
    }
  };

  const filteredAccounts = state.accounts.filter((account) =>
    account.name.toLowerCase().includes(search.toLowerCase())
  );

  const getAccountTypeLabel = (type: string) => {
    switch (type) {
      case 'bank':
        return 'Bank Account';
      case 'cash':
        return 'Cash';
      case 'mobile_money':
        return 'Mobile Money';
      default:
        return type;
    }
  };

  return (
    <Box>
      <Grid container spacing={3}>
        <Grid item xs={12} display="flex" justifyContent="space-between" alignItems="center">
          <Typography variant="h4">Accounts</Typography>
          <Button
            variant="contained"
            color="primary"
            startIcon={<AddIcon />}
            onClick={() => {
              setSelectedAccount(undefined);
              setOpenForm(true);
            }}
          >
            Add Account
          </Button>
        </Grid>

        <Grid item xs={12}>
          <Card>
            <CardContent>
              <TextField
                fullWidth
                size="small"
                label="Search Accounts"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
              />
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12}>
          <TableContainer component={Paper}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Account Name</TableCell>
                  <TableCell>Type</TableCell>
                  <TableCell align="right">Balance</TableCell>
                  <TableCell align="right">Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {filteredAccounts.map((account) => (
                  <TableRow key={account.id}>
                    <TableCell>
                      <Typography variant="subtitle1">{account.name}</Typography>
                    </TableCell>
                    <TableCell>
                      <Chip
                        label={getAccountTypeLabel(account.type)}
                        size="small"
                        color={account.type === 'bank' ? 'primary' : 'default'}
                        variant={account.type === 'bank' ? 'filled' : 'outlined'}
                      />
                    </TableCell>
                    <TableCell align="right">
                      <Typography
                        color={account.balance >= 0 ? 'success.main' : 'error.main'}
                      >
                        ${Math.abs(account.balance).toFixed(2)}
                      </Typography>
                    </TableCell>
                    <TableCell align="right">
                      <IconButton
                        size="small"
                        onClick={() => {
                          setSelectedAccount(account);
                          setOpenForm(true);
                        }}
                      >
                        <EditIcon />
                      </IconButton>
                      <IconButton
                        size="small"
                        onClick={() => handleDeleteAccount(account.id)}
                      >
                        <DeleteIcon />
                      </IconButton>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </Grid>
      </Grid>

      <AccountForm
        open={openForm}
        onClose={() => {
          setOpenForm(false);
          setSelectedAccount(undefined);
        }}
        account={selectedAccount}
        onSubmit={selectedAccount ? handleUpdateAccount : handleAddAccount}
      />
    </Box>
  );
};

export default Accounts;
