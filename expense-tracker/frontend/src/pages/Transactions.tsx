import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Button,
  Card,
  CardContent,
  Grid,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  TablePagination,
  Chip,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import dayjs, { Dayjs } from 'dayjs';
import { Transaction, TransactionType } from '../types';
import { useAppContext } from '../context/AppContext';
import { transactionApi } from '../services/api';
import { TransactionForm } from '../components/forms/TransactionForm';

interface FilterValues {
  startDate: Dayjs | null;
  endDate: Dayjs | null;
  type: string;
  category: string;
  account: string;
}

export const Transactions: React.FC = () => {
  const { state, dispatch } = useAppContext();
  const [openForm, setOpenForm] = useState(false);
  const [selectedTransaction, setSelectedTransaction] = useState<Transaction | null>(null);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [filters, setFilters] = useState<FilterValues>({
    startDate: dayjs().startOf('month'),
    endDate: dayjs().endOf('month'),
    type: 'all',
    category: 'all',
    account: 'all',
  });

  useEffect(() => {
    const fetchTransactions = async () => {
      try {
        const response = await transactionApi.getAll();
        dispatch({ type: 'SET_TRANSACTIONS', payload: response });
      } catch (error) {
        console.error('Error fetching transactions:', error);
      }
    };

    fetchTransactions();
  }, [dispatch]);

  const handleUpdateTransaction = async (data: Omit<Transaction, 'id' | 'createdAt' | 'updatedAt'>) => {
    if (!selectedTransaction) return;
    
    try {
      const response = await transactionApi.update(selectedTransaction.id, data);
      dispatch({ type: 'UPDATE_TRANSACTION', payload: response });
      setSelectedTransaction(null);
      setOpenForm(false);
    } catch (error) {
      console.error('Error updating transaction:', error);
    }
  };

  const handleAddTransaction = async (data: Omit<Transaction, 'id' | 'createdAt' | 'updatedAt'>) => {
    try {
      const response = await transactionApi.create(data);
      dispatch({ type: 'ADD_TRANSACTION', payload: response });
      setOpenForm(false);
    } catch (error) {
      console.error('Error adding transaction:', error);
    }
  };

  const handleDeleteTransaction = async (id: string) => {
    try {
      await transactionApi.delete(id);
      dispatch({ type: 'DELETE_TRANSACTION', payload: id });
    } catch (error) {
      console.error('Error deleting transaction:', error);
    }
  };

  const filteredTransactions = state.transactions.filter((transaction) => {
    const matchesType =
      filters.type === 'all' || transaction.type === filters.type;
    const matchesCategory =
      filters.category === 'all' || transaction.category === filters.category;
    const matchesAccount =
      filters.account === 'all' || transaction.account === filters.account;
    const matchesDate =
      (!filters.startDate || dayjs(transaction.date).isAfter(filters.startDate, 'day') || dayjs(transaction.date).isSame(filters.startDate, 'day')) &&
      (!filters.endDate || dayjs(transaction.date).isBefore(filters.endDate, 'day') || dayjs(transaction.date).isSame(filters.endDate, 'day'));

    return matchesType && matchesCategory && matchesAccount && matchesDate;
  });

  const sortedTransactions = [...filteredTransactions].sort(
    (a, b) => dayjs(b.date).valueOf() - dayjs(a.date).valueOf()
  );

  const paginatedTransactions = sortedTransactions.slice(
    page * rowsPerPage,
    page * rowsPerPage + rowsPerPage
  );

  return (
    <Box>
      <Box sx={{ mb: 3, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <Typography variant="h4">Transactions</Typography>
        <Button
          variant="contained"
          color="primary"
          startIcon={<AddIcon />}
          onClick={() => {
            setSelectedTransaction(null);
            setOpenForm(true);
          }}
        >
          Add Transaction
        </Button>
      </Box>

      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Grid container spacing={2} alignItems="center">
            <Grid item xs={12} sm={6} md={2}>
              <FormControl fullWidth size="small">
                <InputLabel>Type</InputLabel>
                <Select
                  value={filters.type}
                  onChange={(e) => setFilters({ ...filters, type: e.target.value })}
                  label="Type"
                >
                  <MenuItem value="all">All</MenuItem>
                  <MenuItem value="expense">Expense</MenuItem>
                  <MenuItem value="income">Income</MenuItem>
                </Select>
              </FormControl>
            </Grid>

            <Grid item xs={12} sm={6} md={2}>
              <FormControl fullWidth size="small">
                <InputLabel>Category</InputLabel>
                <Select
                  value={filters.category}
                  onChange={(e) => setFilters({ ...filters, category: e.target.value })}
                  label="Category"
                >
                  <MenuItem value="all">All</MenuItem>
                  {state.categories.map((category) => (
                    <MenuItem key={category.id} value={category.name}>
                      {category.name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>

            <Grid item xs={12} sm={6} md={2}>
              <FormControl fullWidth size="small">
                <InputLabel>Account</InputLabel>
                <Select
                  value={filters.account}
                  onChange={(e) => setFilters({ ...filters, account: e.target.value })}
                  label="Account"
                >
                  <MenuItem value="all">All</MenuItem>
                  {state.accounts.map((account) => (
                    <MenuItem key={account.id} value={account.name}>
                      {account.name}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Grid>

            <Grid item xs={12} sm={6} md={3}>
              <DatePicker
                label="Start Date"
                value={filters.startDate}
                onChange={(newValue: Dayjs | null) => setFilters({ ...filters, startDate: newValue })}
                slotProps={{
                  textField: {
                    size: 'small',
                    fullWidth: true,
                  },
                }}
              />
            </Grid>

            <Grid item xs={12} sm={6} md={3}>
              <DatePicker
                label="End Date"
                value={filters.endDate}
                onChange={(newValue: Dayjs | null) => setFilters({ ...filters, endDate: newValue })}
                slotProps={{
                  textField: {
                    size: 'small',
                    fullWidth: true,
                  },
                }}
              />
            </Grid>
          </Grid>
        </CardContent>
      </Card>

      <TableContainer component={Paper}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Date</TableCell>
              <TableCell>Type</TableCell>
              <TableCell>Category</TableCell>
              <TableCell>Account</TableCell>
              <TableCell>Description</TableCell>
              <TableCell align="right">Amount</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {paginatedTransactions.map((transaction) => (
              <TableRow key={transaction.id}>
                <TableCell>{dayjs(transaction.date).format('YYYY-MM-DD')}</TableCell>
                <TableCell>
                  <Chip
                    label={transaction.type}
                    color={transaction.type === 'expense' ? 'error' : 'success'}
                    size="small"
                  />
                </TableCell>
                <TableCell>{transaction.category}</TableCell>
                <TableCell>{transaction.account}</TableCell>
                <TableCell>{transaction.description}</TableCell>
                <TableCell align="right">
                  <Typography
                    color={transaction.type === 'expense' ? 'error' : 'success'}
                  >
                    ${Math.abs(transaction.amount).toFixed(2)}
                  </Typography>
                </TableCell>
                <TableCell align="right">
                  <IconButton
                    size="small"
                    onClick={() => {
                      setSelectedTransaction(transaction);
                      setOpenForm(true);
                    }}
                  >
                    <EditIcon />
                  </IconButton>
                  <IconButton
                    size="small"
                    onClick={() => handleDeleteTransaction(transaction.id)}
                  >
                    <DeleteIcon />
                  </IconButton>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
        <TablePagination
          rowsPerPageOptions={[5, 10, 25]}
          component="div"
          count={filteredTransactions.length}
          rowsPerPage={rowsPerPage}
          page={page}
          onPageChange={(_, newPage) => setPage(newPage)}
          onRowsPerPageChange={(e) => {
            setRowsPerPage(parseInt(e.target.value, 10));
            setPage(0);
          }}
        />
      </TableContainer>

      <TransactionForm
        open={openForm}
        onClose={() => {
          setOpenForm(false);
          setSelectedTransaction(null);
        }}
        transaction={selectedTransaction}
        onSubmit={selectedTransaction ? handleUpdateTransaction : handleAddTransaction}
      />
    </Box>
  );
};
