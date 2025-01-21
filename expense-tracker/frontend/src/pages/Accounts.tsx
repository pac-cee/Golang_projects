import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Button,
  Card,
  CardContent,
  Grid,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  Chip,
} from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import { Account, AccountType } from '../types';
import { useAppContext } from '../context/AppContext';
import { accountApi } from '../services/api';
import { AccountForm } from '../components/forms/AccountForm';

export const Accounts: React.FC = () => {
  const { state, dispatch } = useAppContext();
  const [openForm, setOpenForm] = useState(false);
  const [selectedAccount, setSelectedAccount] = useState<Account | null>(null);

  useEffect(() => {
    const fetchAccounts = async () => {
      try {
        const response = await accountApi.getAll();
        dispatch({ type: 'SET_ACCOUNTS', payload: response });
      } catch (error) {
        console.error('Error fetching accounts:', error);
      }
    };

    fetchAccounts();
  }, [dispatch]);

  const handleAddAccount = async (data: Omit<Account, 'id' | 'createdAt' | 'updatedAt'>) => {
    try {
      const response = await accountApi.create(data);
      dispatch({ type: 'ADD_ACCOUNT', payload: response });
      setOpenForm(false);
    } catch (error) {
      console.error('Error adding account:', error);
    }
  };

  const handleUpdateAccount = async (data: Omit<Account, 'id' | 'createdAt' | 'updatedAt'>) => {
    if (!selectedAccount) return;
    
    try {
      const response = await accountApi.update(selectedAccount.id, data);
      dispatch({ type: 'UPDATE_ACCOUNT', payload: response });
      setSelectedAccount(null);
      setOpenForm(false);
    } catch (error) {
      console.error('Error updating account:', error);
    }
  };

  const handleDeleteAccount = async (id: string) => {
    try {
      await accountApi.delete(id);
      dispatch({ type: 'DELETE_ACCOUNT', payload: id });
    } catch (error) {
      console.error('Error deleting account:', error);
    }
  };

  const calculateTotalBalance = () => {
    return state.accounts.reduce((total, account) => total + account.balance, 0);
  };

  return (
    <Box>
      <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4">Accounts</Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => {
            setSelectedAccount(null);
            setOpenForm(true);
          }}
        >
          Add Account
        </Button>
      </Box>

      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Typography variant="h6" gutterBottom>
            Total Balance
          </Typography>
          <Typography variant="h4" color="primary">
            ${calculateTotalBalance().toFixed(2)}
          </Typography>
        </CardContent>
      </Card>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Type</TableCell>
              <TableCell>Currency</TableCell>
              <TableCell align="right">Balance</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {state.accounts.map((account) => (
              <TableRow key={account.id}>
                <TableCell>{account.name}</TableCell>
                <TableCell>
                  <Chip
                    label={account.type}
                    color={account.type === 'savings' ? 'success' : 'default'}
                    size="small"
                  />
                </TableCell>
                <TableCell>{account.currency}</TableCell>
                <TableCell align="right">
                  <Typography
                    color={account.balance >= 0 ? 'success.main' : 'error.main'}
                  >
                    ${account.balance.toFixed(2)}
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

      <AccountForm
        open={openForm}
        onClose={() => {
          setOpenForm(false);
          setSelectedAccount(null);
        }}
        account={selectedAccount}
        onSubmit={selectedAccount ? handleUpdateAccount : handleAddAccount}
      />
    </Box>
  );
};
