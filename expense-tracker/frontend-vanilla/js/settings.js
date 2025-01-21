// Initialize Notyf
const notyf = new Notyf({
    duration: 3000,
    position: { x: 'right', y: 'top' },
    types: [
        {
            type: 'warning',
            background: '#fb923c',
            icon: false
        }
    ]
});

// DOM Elements
const monthlyLimitInput = document.getElementById('monthly-limit');
const alertThresholdInput = document.getElementById('alert-threshold');
const saveLimitsBtn = document.querySelector('.save-limits-btn');
const emailNotifications = document.getElementById('email-notifications');
const pushNotifications = document.getElementById('push-notifications');
const exportButtons = document.querySelectorAll('.export-btn');

// Load saved settings
function loadSettings() {
    const settings = JSON.parse(localStorage.getItem('settings') || '{}');
    
    // Set spending limits
    monthlyLimitInput.value = settings.monthlyLimit || '';
    alertThresholdInput.value = settings.alertThreshold || 80;
    
    // Set notification preferences
    emailNotifications.checked = settings.emailNotifications || false;
    pushNotifications.checked = settings.pushNotifications || false;
}

// Save spending limits
saveLimitsBtn.addEventListener('click', () => {
    const monthlyLimit = parseFloat(monthlyLimitInput.value);
    const alertThreshold = parseInt(alertThresholdInput.value);
    
    if (!monthlyLimit || monthlyLimit <= 0) {
        notyf.error('Please enter a valid monthly limit');
        return;
    }
    
    if (alertThreshold < 1 || alertThreshold > 100) {
        notyf.error('Alert threshold must be between 1 and 100');
        return;
    }
    
    const settings = JSON.parse(localStorage.getItem('settings') || '{}');
    settings.monthlyLimit = monthlyLimit;
    settings.alertThreshold = alertThreshold;
    localStorage.setItem('settings', JSON.stringify(settings));
    
    notyf.success('Spending limits saved successfully');
});

// Handle notification toggles
[emailNotifications, pushNotifications].forEach(toggle => {
    toggle.addEventListener('change', () => {
        const settings = JSON.parse(localStorage.getItem('settings') || '{}');
        settings[toggle.id] = toggle.checked;
        localStorage.setItem('settings', JSON.stringify(settings));
        
        notyf.success(`${toggle.checked ? 'Enabled' : 'Disabled'} ${toggle.id.replace('-', ' ')}`);
    });
});

// Check spending against limit
function checkSpendingLimit(amount) {
    const settings = JSON.parse(localStorage.getItem('settings') || '{}');
    if (!settings.monthlyLimit) return;
    
    const totalSpent = getCurrentMonthTotal() + amount;
    const limitThreshold = (settings.monthlyLimit * settings.alertThreshold) / 100;
    
    if (totalSpent >= settings.monthlyLimit) {
        notyf.error('Monthly spending limit exceeded!');
    } else if (totalSpent >= limitThreshold) {
        notyf.warning(`Approaching monthly spending limit (${settings.alertThreshold}%)`);
    }
}

// Get current month's total (mock function)
function getCurrentMonthTotal() {
    // This should be replaced with actual API call to get current month's total
    return 0;
}

// Export data handlers
exportButtons.forEach(btn => {
    btn.addEventListener('click', () => {
        const format = btn.dataset.format;
        exportData(format);
    });
});

// Export data function (mock)
function exportData(format) {
    // This should be replaced with actual export functionality
    notyf.success(`Exporting data as ${format.toUpperCase()}...`);
}

// Load settings on page load
loadSettings();

// Export functions for use in other modules
window.expenseTracker = window.expenseTracker || {};
window.expenseTracker.checkSpendingLimit = checkSpendingLimit;
