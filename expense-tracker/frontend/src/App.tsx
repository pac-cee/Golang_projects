import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ThemeProvider } from '@mui/material/styles';
import { LocalizationProvider } from '@mui/x-date-pickers';
import { AdapterDayjs } from '@mui/x-date-pickers/AdapterDayjs';
import CssBaseline from '@mui/material/CssBaseline';
import theme from './theme';
import { AppProvider } from './context/AppContext';
import MainLayout from './layouts/MainLayout';
import Dashboard from './pages/Dashboard';
import { Transactions } from './pages/Transactions';
import { Categories } from './pages/Categories';
import { BudgetPage } from './pages/Budget';
import Reports from './pages/Reports';
import { Accounts } from './pages/Accounts';

function App() {
  return (
    <ThemeProvider theme={theme}>
      <LocalizationProvider dateAdapter={AdapterDayjs}>
        <CssBaseline />
        <AppProvider>
          <Router>
            <MainLayout>
              <Routes>
                <Route path="/" element={<Dashboard />} />
                <Route path="/transactions" element={<Transactions />} />
                <Route path="/categories" element={<Categories />} />
                <Route path="/budget" element={<BudgetPage />} />
                <Route path="/reports" element={<Reports />} />
                <Route path="/accounts" element={<Accounts />} />
              </Routes>
            </MainLayout>
          </Router>
        </AppProvider>
      </LocalizationProvider>
    </ThemeProvider>
  );
}

export default App;
