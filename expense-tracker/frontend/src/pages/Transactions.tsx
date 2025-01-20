import React, { useState, useEffect } from 'react';
import {
  Box,
  Button,
  Card,
  CardContent,
  Grid,
  IconButton,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination,
  TextField,
  Typography,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Chip,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import AddIcon from '@mui/icons-material/Add';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';
import { Transaction } from '../types';
import { useAppContext } from '../context/AppContext';
import { transactionApi } from '../services/api';
import TransactionForm from '../components/forms/TransactionForm';

const Transactions: React.FC = () => {
  const { state, dispatch } = useAppContext();
  const [openForm, setOpenForm] = useState(false);
  const [selectedTransaction, setSelectedTransaction] = useState<Transaction | undefined>();
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [filters, setFilters] = useState({
    type: 'all',
    category: 'all',
    account: 'all',
    startDate: null as Date | null,
    endDate: null as Date | null,
    search: '',
  });

  useEffect(() => {
    const fetchTransactions = async () => {
      try {
        dispatch({ type: 'SET_LOADING', payload: true });
        const transactions = await transactionApi.getAll();
        dispatch({ type: 'SET_TRANSACTIONS', payload: transactions });
      } catch (error) {
        dispatch({ type: 'SET_ERROR', payload: 'Error fetching transactions' });
      } finally {
        dispatch({ type: 'SET_LOADING', payload: false });
      }
    };

    fetchTransactions();
  }, [dispatch]);

  const handleAddTransaction = async (
    transaction: Omit<Transaction, 'id' | 'createdAt' | 'updatedAt'>
  ) => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true });
      const newTransaction = await transactionApi.create(transaction);
      dispatch({ type: 'ADD_TRANSACTION', payload: newTransaction });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Error adding transaction' });
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  };

  const handleUpdateTransaction = async (
    transaction: Omit<Transaction, 'createdAt' | 'updatedAt'>
  ) => {
    try {
      dispatch({ type: 'SET_LOADING', payload: true });
      const updatedTransaction = await transactionApi.update(transaction.id, transaction);
      dispatch({ type: 'UPDATE_TRANSACTION', payload: updatedTransaction });
    } catch (error) {
      dispatch({ type: 'SET_ERROR', payload: 'Error updating transaction' });
    } finally {
      dispatch({ type: 'SET_LOADING', payload: false });
    }
  };

  const handleDeleteTransaction = async (id: string) => {
    if (window.confirm('Are you sure you want to delete this transaction?')) {
      try {
        dispatch({ type: 'SET_LOADING', payload: true });
        await transactionApi.delete(id);
        dispatch({ type: 'DELETE_TRANSACTION', payload: id });
      } catch (error) {
        dispatch({ type: 'SET_ERROR', payload: 'Error deleting transaction' });
      } finally {
        dispatch({ type: 'SET_LOADING', payload: false });
      }
    }
  };

  const filteredTransactions = state.transactions.filter((transaction) => {
    const matchesType =
      filters.type === 'all' ||
      (filters.type === 'income' && transaction.amount > 0) ||
      (filters.type === 'expense' && transaction.amount < 0);

    const matchesCategory =
      filters.category === 'all' || transaction.category === filters.category;

    const matchesAccount =
      filters.account === 'all' || transaction.account === filters.account;

    const matchesDate =
      (!filters.startDate || new Date(transaction.date) >= filters.startDate) &&
      (!filters.endDate || new Date(transaction.date) <= filters.endDate);

    const matchesSearch = filters.search
      ? transaction.description.toLowerCase().includes(filters.search.toLowerCase()) ||
        transaction.category.toLowerCase().includes(filters.search.toLowerCase()) ||
        transaction.account.toLowerCase().includes(filters.search.toLowerCase())
      : true;

    return matchesType && matchesCategory && matchesAccount && matchesDate && matchesSearch;
  });

  const sortedTransactions = [...filteredTransactions].sort(
    (a, b) => new Date(b.date).getTime() - new Date(a.date).getTime()
  );

  const paginatedTransactions = sortedTransactions.slice(
    page * rowsPerPage,
    (page + 1) * rowsPerPage
  );

  return (
    <Box>
      <Grid container spacing={3}>
        <Grid item xs={12} display="flex" justifyContent="space-between" alignItems="center">
          <Typography variant="h4">Transactions</Typography>
          <Button
            variant="contained"
            color="primary"
            startIcon={<AddIcon />}
            onClick={() => {
              setSelectedTransaction(undefined);
              setOpenForm(true);
            }}
          >
            Add Transaction
          </Button>
        </Grid>

        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6} md={2}>
                  <FormControl fullWidth size="small">
                    <InputLabel>Type</InputLabel>
                    <Select
                      value={filters.type}
                      onChange={(e) => setFilters({ ...filters, type: e.target.value })}
                      label="Type"
                    >
                      <MenuItem value="all">All</MenuItem>
                      <MenuItem value="income">Income</MenuItem>
                      <MenuItem value="expense">Expense</MenuItem>
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
                      <MenuItem value="all">All Categories</MenuItem>
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
                      <MenuItem value="all">All Accounts</MenuItem>
                      {state.accounts.map((account) => (
                        <MenuItem key={account.id} value={account.name}>
                          {account.name}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>
                </Grid>

                <Grid item xs={12} sm={6} md={2}>
                  <DatePicker
                    label="Start Date"
                    value={filters.startDate}
                    onChange={(newValue) => setFilters({ ...filters, startDate: newValue })}
                    slotProps={{
                      textField: {
                        size: 'small',
                        fullWidth: true,
                      },
                    }}
                  />
                </Grid>

                <Grid item xs={12} sm={6} md={2}>
                  <DatePicker
                    label="End Date"
                    value={filters.endDate}
                    onChange={(newValue) => setFilters({ ...filters, endDate: newValue })}
                    slotProps={{
                      textField: {
                        size: 'small',
                        fullWidth: true,
                      },
                    }}
                  />
                </Grid>

                <Grid item xs={12} sm={6} md={2}>
                  <TextField
                    fullWidth
                    size="small"
                    label="Search"
                    value={filters.search}
                    onChange={(e) => setFilters({ ...filters, search: e.target.value })}
                  />
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12}>
          <TableContainer component={Paper}>
            <Table>
              <TableHead>
                <TableRow>
                  <TableCell>Date</TableCell>
                  <TableCell>Description</TableCell>
                  <TableCell>Category</TableCell>
                  <TableCell>Account</TableCell>
                  <TableCell align="right">Amount</TableCell>
                  <TableCell align="right">Actions</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {paginatedTransactions.map((transaction) => (
                  <TableRow key={transaction.id}>
                    <TableCell>
                      {new Date(transaction.date).toLocaleDateString()}
                    </TableCell>
                    <TableCell>{transaction.description}</TableCell>
                    <TableCell>
                      <Chip
                        label={transaction.category}
                        size="small"
                        color={transaction.amount >= 0 ? 'success' : 'error'}
                      />
                      {transaction.subcategory && (
                        <Chip
                          label={transaction.subcategory}
                          size="small"
                          variant="outlined"
                          sx={{ ml: 1 }}
                        />
                      )}
                    </TableCell>
                    <TableCell>{transaction.account}</TableCell>
                    <TableCell align="right">
                      <Typography
                        color={transaction.amount >= 0 ? 'success.main' : 'error.main'}
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
              component="div"
              count={filteredTransactions.length}
              page={page}
              onPageChange={(_, newPage) => setPage(newPage)}
              rowsPerPage={rowsPerPage}
              onRowsPerPageChange={(e) => {
                setRowsPerPage(parseInt(e.target.value, 10));
                setPage(0);
              }}
            />
          </TableContainer>
        </Grid>
      </Grid>

      <TransactionForm
        open={openForm}
        onClose={() => {
          setOpenForm(false);
          setSelectedTransaction(undefined);
        }}
        transaction={selectedTransaction}
        onSubmit={selectedTransaction ? handleUpdateTransaction : handleAddTransaction}
      />
    </Box>
  );
};

export default Transactions;
