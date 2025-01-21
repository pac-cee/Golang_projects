import React from 'react';
import {
  Box,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
  InputAdornment,
} from '@mui/material';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';
import dayjs, { Dayjs } from 'dayjs';
import { TransactionType, TransactionFilters, Category, Account } from '../../types';
import { useAppContext } from '../../context/AppContext';

interface TransactionFiltersProps {
  filters: TransactionFilters;
  onFilterChange: (filters: TransactionFilters) => void;
}

export const TransactionFiltersComponent: React.FC<TransactionFiltersProps> = ({
  filters,
  onFilterChange,
}) => {
  const { state } = useAppContext();

  const handleFilterChange = (
    field: keyof TransactionFilters,
    value: string | number | undefined
  ) => {
    onFilterChange({
      ...filters,
      [field]: value,
    });
  };

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mb: 3 }}>
      <Box sx={{ display: 'flex', gap: 2 }}>
        <DatePicker
          label="Start Date"
          value={filters.startDate ? dayjs(filters.startDate) : null}
          onChange={(date: Dayjs | null) =>
            handleFilterChange('startDate', date?.toISOString())
          }
          slotProps={{
            textField: {
              fullWidth: true,
            },
          }}
        />

        <DatePicker
          label="End Date"
          value={filters.endDate ? dayjs(filters.endDate) : null}
          onChange={(date: Dayjs | null) =>
            handleFilterChange('endDate', date?.toISOString())
          }
          slotProps={{
            textField: {
              fullWidth: true,
            },
          }}
        />
      </Box>

      <Box sx={{ display: 'flex', gap: 2 }}>
        <FormControl fullWidth>
          <InputLabel>Type</InputLabel>
          <Select
            value={filters.type || ''}
            onChange={(e) =>
              handleFilterChange('type', e.target.value as TransactionType)
            }
            label="Type"
          >
            <MenuItem value="">All</MenuItem>
            <MenuItem value="expense">Expense</MenuItem>
            <MenuItem value="income">Income</MenuItem>
          </Select>
        </FormControl>

        <FormControl fullWidth>
          <InputLabel>Category</InputLabel>
          <Select
            value={filters.category || ''}
            onChange={(e) => handleFilterChange('category', e.target.value)}
            label="Category"
          >
            <MenuItem value="">All</MenuItem>
            {state.categories
              .filter((category: Category) =>
                filters.type ? category.type === filters.type : true
              )
              .map((category: Category) => (
                <MenuItem key={category.id} value={category.name}>
                  {category.name}
                </MenuItem>
              ))}
          </Select>
        </FormControl>

        <FormControl fullWidth>
          <InputLabel>Account</InputLabel>
          <Select
            value={filters.account || ''}
            onChange={(e) => handleFilterChange('account', e.target.value)}
            label="Account"
          >
            <MenuItem value="">All</MenuItem>
            {state.accounts.map((account: Account) => (
              <MenuItem key={account.id} value={account.name}>
                {account.name}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Box>

      <Box sx={{ display: 'flex', gap: 2 }}>
        <TextField
          label="Min Amount"
          type="number"
          value={filters.minAmount || ''}
          onChange={(e) =>
            handleFilterChange(
              'minAmount',
              e.target.value ? Number(e.target.value) : undefined
            )
          }
          InputProps={{
            startAdornment: <InputAdornment position="start">$</InputAdornment>,
          }}
          fullWidth
        />

        <TextField
          label="Max Amount"
          type="number"
          value={filters.maxAmount || ''}
          onChange={(e) =>
            handleFilterChange(
              'maxAmount',
              e.target.value ? Number(e.target.value) : undefined
            )
          }
          InputProps={{
            startAdornment: <InputAdornment position="start">$</InputAdornment>,
          }}
          fullWidth
        />
      </Box>
    </Box>
  );
};
