import { useState } from 'react';
import {
  Grid,
  TextField,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Button,
  Box,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import { useAppContext } from '../../context/AppContext';

interface FilterValues {
  startDate: Date | null;
  endDate: Date | null;
  type: string;
  category: string;
  account: string;
  minAmount: string;
  maxAmount: string;
}

interface TransactionFiltersProps {
  onFilter: (filters: FilterValues) => void;
}

const TransactionFilters = ({ onFilter }: TransactionFiltersProps) => {
  const { state } = useAppContext();
  const [filters, setFilters] = useState<FilterValues>({
    startDate: null,
    endDate: null,
    type: '',
    category: '',
    account: '',
    minAmount: '',
    maxAmount: '',
  });

  const handleChange = (field: keyof FilterValues) => (event: any) => {
    const value = event.target.value;
    setFilters((prev) => ({
      ...prev,
      [field]: value,
    }));
  };

  const handleDateChange = (field: 'startDate' | 'endDate') => (date: Date | null) => {
    setFilters((prev) => ({
      ...prev,
      [field]: date,
    }));
  };

  const handleApplyFilters = () => {
    onFilter(filters);
  };

  const handleResetFilters = () => {
    setFilters({
      startDate: null,
      endDate: null,
      type: '',
      category: '',
      account: '',
      minAmount: '',
      maxAmount: '',
    });
    onFilter({
      startDate: null,
      endDate: null,
      type: '',
      category: '',
      account: '',
      minAmount: '',
      maxAmount: '',
    });
  };

  return (
    <Box sx={{ p: 2 }}>
      <Grid container spacing={2}>
        <Grid item xs={12} sm={6} md={3}>
          <DatePicker
            label="Start Date"
            value={filters.startDate}
            onChange={handleDateChange('startDate')}
            slotProps={{ textField: { fullWidth: true } }}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <DatePicker
            label="End Date"
            value={filters.endDate}
            onChange={handleDateChange('endDate')}
            slotProps={{ textField: { fullWidth: true } }}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <FormControl fullWidth>
            <InputLabel>Type</InputLabel>
            <Select
              value={filters.type}
              label="Type"
              onChange={handleChange('type')}
            >
              <MenuItem value="">All</MenuItem>
              <MenuItem value="income">Income</MenuItem>
              <MenuItem value="expense">Expense</MenuItem>
            </Select>
          </FormControl>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <FormControl fullWidth>
            <InputLabel>Category</InputLabel>
            <Select
              value={filters.category}
              label="Category"
              onChange={handleChange('category')}
            >
              <MenuItem value="">All</MenuItem>
              {state.categories.map((category) => (
                <MenuItem key={category.id} value={category.id}>
                  {category.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <FormControl fullWidth>
            <InputLabel>Account</InputLabel>
            <Select
              value={filters.account}
              label="Account"
              onChange={handleChange('account')}
            >
              <MenuItem value="">All</MenuItem>
              {state.accounts.map((account) => (
                <MenuItem key={account.id} value={account.id}>
                  {account.name}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <TextField
            fullWidth
            label="Min Amount"
            type="number"
            value={filters.minAmount}
            onChange={handleChange('minAmount')}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <TextField
            fullWidth
            label="Max Amount"
            type="number"
            value={filters.maxAmount}
            onChange={handleChange('maxAmount')}
          />
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Box sx={{ display: 'flex', gap: 1 }}>
            <Button
              variant="contained"
              onClick={handleApplyFilters}
              fullWidth
            >
              Apply Filters
            </Button>
            <Button
              variant="outlined"
              onClick={handleResetFilters}
              fullWidth
            >
              Reset
            </Button>
          </Box>
        </Grid>
      </Grid>
    </Box>
  );
};

export default TransactionFilters;
