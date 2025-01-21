// DOM Elements
const addExpenseBtn = document.getElementById('add-expense-btn');
const expenseModal = document.getElementById('expense-modal');
const expenseForm = document.getElementById('expense-form');
const cancelExpenseBtn = document.getElementById('cancel-expense');

// Load dashboard data
async function loadDashboardData() {
    const token = localStorage.getItem('token');
    if (!token) return;

    try {
        // Load summary data
        const summaryResponse = await fetch(`${API_URL}/expenses/summary`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        const summaryData = await summaryResponse.json();
        updateSummaryCards(summaryData);

        // Load recent expenses
        const expensesResponse = await fetch(`${API_URL}/expenses?limit=5`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        const expensesData = await expensesResponse.json();
        updateRecentExpenses(expensesData);

        // Update charts
        updateCharts();
    } catch (error) {
        console.error('Failed to load dashboard data:', error);
    }
}

// Update summary cards with data
function updateSummaryCards(data) {
    const cards = document.querySelectorAll('.summary-cards .card');
    
    // Total Expenses
    cards[0].querySelector('.amount').textContent = formatCurrency(data.totalExpenses);
    
    // Average Daily
    cards[1].querySelector('.amount').textContent = formatCurrency(data.averageDaily);
    
    // Highest Expense
    cards[2].querySelector('.amount').textContent = formatCurrency(data.highestExpense);
}

// Update recent expenses list
function updateRecentExpenses(expenses) {
    const expenseList = document.querySelector('.expense-list');
    expenseList.innerHTML = '';

    expenses.forEach(expense => {
        const expenseItem = document.createElement('div');
        expenseItem.className = 'expense-item';
        expenseItem.innerHTML = `
            <div class="expense-info">
                <h4>${expense.description}</h4>
                <span class="category">${expense.category}</span>
                <span class="date">${formatDate(expense.date)}</span>
            </div>
            <div class="expense-amount">${formatCurrency(expense.amount)}</div>
        `;
        expenseList.appendChild(expenseItem);
    });
}

// Modal handlers
addExpenseBtn.addEventListener('click', () => {
    expenseModal.classList.remove('hidden');
});

cancelExpenseBtn.addEventListener('click', () => {
    expenseModal.classList.add('hidden');
    expenseForm.reset();
});

// Close modal when clicking outside
expenseModal.addEventListener('click', (e) => {
    if (e.target === expenseModal) {
        expenseModal.classList.add('hidden');
        expenseForm.reset();
    }
});

// Handle expense form submission
expenseForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const token = localStorage.getItem('token');
    if (!token) return;

    const formData = new FormData(expenseForm);
    const expenseData = {
        description: formData.get('description'),
        amount: parseFloat(formData.get('amount')),
        category: formData.get('category'),
        date: formData.get('date'),
        notes: formData.get('notes')
    };

    try {
        const response = await fetch(`${API_URL}/expenses`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(expenseData)
        });

        if (!response.ok) {
            throw new Error('Failed to add expense');
        }

        // Refresh dashboard data
        loadDashboardData();
        
        // Close modal and reset form
        expenseModal.classList.add('hidden');
        expenseForm.reset();
    } catch (error) {
        alert('Failed to add expense. Please try again.');
    }
});

// Utility functions
function formatCurrency(amount) {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD'
    }).format(amount);
}

function formatDate(dateString) {
    return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
}

// Navigation
const navLinks = document.querySelectorAll('.nav-links li:not(.logout)');
navLinks.forEach(link => {
    link.addEventListener('click', () => {
        navLinks.forEach(l => l.classList.remove('active'));
        link.classList.add('active');
        // TODO: Implement navigation to different sections
    });
});
