import { expenses as expensesApi } from './api.js';
import { notify, formatCurrency, formatDate, requireAuth, chartUtils } from './utils.js';

// Initialize dashboard when document is ready
document.addEventListener('DOMContentLoaded', () => {
    if (!requireAuth()) return;
    initializeDashboard();
});

async function initializeDashboard() {
    await Promise.all([
        loadSummaryData(),
        loadExpenseCharts(),
        loadRecentExpenses()
    ]);
}

async function loadSummaryData() {
    try {
        const response = await expensesApi.getSummary();
        updateSummaryCards(response);
    } catch (error) {
        notify.error('Failed to load summary data');
        console.error(error);
    }
}

function updateSummaryCards(data) {
    // Update total expenses
    const totalExpenses = document.getElementById('totalExpenses');
    if (totalExpenses) {
        totalExpenses.textContent = formatCurrency(data.totalExpenses);
    }

    // Update monthly spending
    const monthlySpending = document.getElementById('monthlySpending');
    if (monthlySpending) {
        monthlySpending.textContent = formatCurrency(data.monthlyTotal);
        const progress = (data.monthlyTotal / data.monthlyBudget) * 100;
        const progressBar = monthlySpending.querySelector('.progress');
        if (progressBar) {
            progressBar.style.width = `${Math.min(progress, 100)}%`;
            progressBar.style.backgroundColor = progress > 100 ? '#ef4444' : '#10b981';
        }
    }

    // Update category spending
    const categorySpending = document.getElementById('topCategories');
    if (categorySpending) {
        categorySpending.innerHTML = data.categoryTotals
            .slice(0, 3)
            .map(cat => `
                <div class="category-item">
                    <span class="category-name">${cat.name}</span>
                    <span class="category-amount">${formatCurrency(cat.total)}</span>
                </div>
            `).join('');
    }
}

async function loadExpenseCharts() {
    try {
        const [monthlyData, categoryData] = await Promise.all([
            expensesApi.getMonthlyTrend(),
            expensesApi.getCategoryDistribution()
        ]);

        createMonthlyTrendChart(monthlyData);
        createCategoryDistributionChart(categoryData);
    } catch (error) {
        notify.error('Failed to load chart data');
        console.error(error);
    }
}

function createMonthlyTrendChart(data) {
    const ctx = document.getElementById('monthlyTrendChart');
    if (!ctx) return;

    chartUtils.createChart(ctx, 'line', {
        labels: data.map(d => d.month),
        datasets: [{
            label: 'Monthly Expenses',
            data: data.map(d => d.total),
            borderColor: '#3b82f6',
            tension: 0.1
        }]
    }, {
        scales: {
            y: {
                beginAtZero: true,
                ticks: {
                    callback: value => formatCurrency(value)
                }
            }
        },
        plugins: {
            tooltip: {
                callbacks: {
                    label: context => formatCurrency(context.raw)
                }
            }
        }
    });
}

function createCategoryDistributionChart(data) {
    const ctx = document.getElementById('categoryDistributionChart');
    if (!ctx) return;

    chartUtils.createChart(ctx, 'doughnut', {
        labels: data.map(d => d.category),
        datasets: [{
            data: data.map(d => d.total),
            backgroundColor: chartUtils.colors
        }]
    }, {
        plugins: {
            tooltip: {
                callbacks: {
                    label: context => formatCurrency(context.raw)
                }
            }
        }
    });
}

async function loadRecentExpenses() {
    try {
        const response = await expensesApi.getAll({ limit: 5, sortBy: 'date', sortOrder: 'desc' });
        updateRecentExpenses(response.expenses);
    } catch (error) {
        notify.error('Failed to load recent expenses');
        console.error(error);
    }
}

function updateRecentExpenses(expenses) {
    const recentExpensesList = document.getElementById('recentExpenses');
    if (!recentExpensesList) return;

    recentExpensesList.innerHTML = expenses.map(expense => `
        <div class="expense-item">
            <div class="expense-info">
                <div class="expense-primary">
                    <span class="expense-description">${expense.description}</span>
                    <span class="expense-amount">${formatCurrency(expense.amount)}</span>
                </div>
                <div class="expense-secondary">
                    <span class="expense-date">${formatDate(expense.date)}</span>
                    <span class="expense-category" style="color: ${expense.category.color}">
                        <i class="fas ${expense.category.icon}"></i>
                        ${expense.category.name}
                    </span>
                </div>
            </div>
        </div>
    `).join('');
}
