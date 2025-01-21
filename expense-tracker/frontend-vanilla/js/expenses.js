import { expenses as expensesApi } from './api.js';
import { notify, formatCurrency, formatDate, requireAuth } from './utils.js';

// State management
let expenses = [];
let currentPage = 1;
let totalPages = 1;
let filters = {
    search: '',
    category: '',
    timeRange: '',
    sortBy: 'date',
    sortOrder: 'desc'
};

// Initialize when document is ready
document.addEventListener('DOMContentLoaded', () => {
    if (!requireAuth()) return;
    initializeExpenses();
});

async function initializeExpenses() {
    // Initialize UI elements
    initializeSearchAndFilters();
    initializePagination();
    initializeExpenseForm();
    
    // Load initial data
    await loadExpenses();
}

function initializeSearchAndFilters() {
    const searchInput = document.getElementById('searchExpenses');
    const categoryFilter = document.getElementById('categoryFilter');
    const timeFilter = document.getElementById('timeFilter');

    searchInput?.addEventListener('input', debounce(async (e) => {
        filters.search = e.target.value;
        currentPage = 1;
        await loadExpenses();
    }, 300));

    categoryFilter?.addEventListener('change', async (e) => {
        filters.category = e.target.value;
        currentPage = 1;
        await loadExpenses();
    });

    timeFilter?.addEventListener('change', async (e) => {
        filters.timeRange = e.target.value;
        currentPage = 1;
        await loadExpenses();
    });
}

function initializePagination() {
    const prevPageBtn = document.getElementById('prevPage');
    const nextPageBtn = document.getElementById('nextPage');

    prevPageBtn?.addEventListener('click', async () => {
        if (currentPage > 1) {
            currentPage--;
            await loadExpenses();
        }
    });

    nextPageBtn?.addEventListener('click', async () => {
        if (currentPage < totalPages) {
            currentPage++;
            await loadExpenses();
        }
    });
}

function initializeExpenseForm() {
    const addExpenseBtn = document.getElementById('addExpenseBtn');
    const expenseModal = document.getElementById('expenseModal');
    const expenseForm = document.getElementById('expenseForm');
    const cancelExpenseBtn = document.getElementById('cancelExpense');

    addExpenseBtn?.addEventListener('click', () => {
        currentExpense = null;
        expenseForm.reset();
        expenseModal.classList.remove('hidden');
    });

    cancelExpenseBtn?.addEventListener('click', () => {
        expenseModal.classList.add('hidden');
    });

    expenseForm?.addEventListener('submit', handleExpenseSubmit);
}

async function loadExpenses() {
    try {
        const params = {
            page: currentPage,
            limit: 10,
            ...filters
        };

        const response = await expensesApi.getAll(params);
        expenses = response.expenses;
        totalPages = response.totalPages;
        
        renderExpenses();
        updatePagination();
    } catch (error) {
        notify.error('Failed to load expenses');
        console.error(error);
    }
}

function renderExpenses() {
    const tbody = document.getElementById('expensesTableBody');
    if (!tbody) return;

    tbody.innerHTML = expenses.map(expense => `
        <tr>
            <td>${formatDate(expense.date)}</td>
            <td>${expense.description}</td>
            <td>
                <span class="category-badge ${expense.category.toLowerCase()}">
                    ${expense.category}
                </span>
            </td>
            <td>${formatCurrency(expense.amount)}</td>
            <td>
                <button class="icon-btn" onclick="editExpense('${expense._id}')">
                    <i class="fas fa-edit"></i>
                </button>
                <button class="icon-btn" onclick="deleteExpense('${expense._id}')">
                    <i class="fas fa-trash"></i>
                </button>
            </td>
        </tr>
    `).join('');
}

function updatePagination() {
    const paginationText = document.querySelector('.pagination span');
    if (paginationText) {
        paginationText.textContent = `Page ${currentPage} of ${totalPages}`;
    }

    const prevPageBtn = document.getElementById('prevPage');
    const nextPageBtn = document.getElementById('nextPage');

    if (prevPageBtn) {
        prevPageBtn.disabled = currentPage === 1;
    }
    if (nextPageBtn) {
        nextPageBtn.disabled = currentPage === totalPages;
    }
}

async function handleExpenseSubmit(e) {
    e.preventDefault();
    const form = e.target;
    const formData = new FormData(form);

    const expenseData = {
        description: formData.get('description'),
        amount: parseFloat(formData.get('amount')),
        category: formData.get('category'),
        date: formData.get('date'),
        notes: formData.get('notes')
    };

    try {
        if (currentExpense) {
            await expensesApi.update(currentExpense._id, expenseData);
            notify.success('Expense updated successfully');
        } else {
            await expensesApi.create(expenseData);
            notify.success('Expense added successfully');
        }

        document.getElementById('expenseModal').classList.add('hidden');
        await loadExpenses();
    } catch (error) {
        notify.error(error.message || 'Failed to save expense');
    }
}

async function editExpense(id) {
    const expense = expenses.find(e => e._id === id);
    if (!expense) return;

    currentExpense = expense;
    const form = document.getElementById('expenseForm');
    if (!form) return;

    form.description.value = expense.description;
    form.amount.value = expense.amount;
    form.category.value = expense.category;
    form.date.value = expense.date.split('T')[0];
    form.notes.value = expense.notes || '';

    document.getElementById('expenseModal').classList.remove('hidden');
}

async function deleteExpense(id) {
    if (!confirm('Are you sure you want to delete this expense?')) return;

    try {
        await expensesApi.delete(id);
        notify.success('Expense deleted successfully');
        await loadExpenses();
    } catch (error) {
        notify.error('Failed to delete expense');
    }
}

// Utility function for debouncing
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Export functions that need to be accessed globally
window.editExpense = editExpense;
window.deleteExpense = deleteExpense;
