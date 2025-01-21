import { settings as settingsApi } from './api.js';
import { notify, requireAuth } from './utils.js';

// Initialize when document is ready
document.addEventListener('DOMContentLoaded', () => {
    if (!requireAuth()) return;
    initializeSettings();
});

async function initializeSettings() {
    initializeThemeSettings();
    initializeNotificationSettings();
    initializeSpendingLimits();
    initializeDataManagement();
    await loadSettings();
}

function initializeThemeSettings() {
    const themeToggle = document.getElementById('themeToggle');
    const currentTheme = localStorage.getItem('theme') || 'light';
    
    if (themeToggle) {
        themeToggle.checked = currentTheme === 'dark';
        themeToggle.addEventListener('change', (e) => {
            const newTheme = e.target.checked ? 'dark' : 'light';
            localStorage.setItem('theme', newTheme);
            document.documentElement.classList.toggle('dark', e.target.checked);
        });
    }
}

function initializeNotificationSettings() {
    const emailToggle = document.getElementById('emailNotifications');
    const pushToggle = document.getElementById('pushNotifications');
    const weeklyToggle = document.getElementById('weeklyReports');
    const monthlyToggle = document.getElementById('monthlyReports');

    const toggles = [emailToggle, pushToggle, weeklyToggle, monthlyToggle];
    toggles.forEach(toggle => {
        if (toggle) {
            toggle.addEventListener('change', handleNotificationChange);
        }
    });
}

function initializeSpendingLimits() {
    const limitForm = document.getElementById('spendingLimitsForm');
    if (limitForm) {
        limitForm.addEventListener('submit', handleSpendingLimitsSubmit);
    }
}

function initializeDataManagement() {
    const exportBtn = document.getElementById('exportData');
    const importBtn = document.getElementById('importData');
    const deleteBtn = document.getElementById('deleteAccount');

    if (exportBtn) {
        exportBtn.addEventListener('click', handleDataExport);
    }

    if (importBtn) {
        importBtn.addEventListener('change', handleDataImport);
    }

    if (deleteBtn) {
        deleteBtn.addEventListener('click', handleAccountDeletion);
    }
}

async function loadSettings() {
    try {
        const settings = await settingsApi.getProfile();
        populateSettings(settings);
    } catch (error) {
        notify.error('Failed to load settings');
        console.error(error);
    }
}

function populateSettings(settings) {
    // Spending Limits
    const monthlyLimit = document.getElementById('monthlyLimit');
    const alertThreshold = document.getElementById('alertThreshold');
    
    if (monthlyLimit) {
        monthlyLimit.value = settings.monthlyLimit || '';
    }
    if (alertThreshold) {
        alertThreshold.value = settings.alertThreshold || '';
    }

    // Notification Settings
    const emailToggle = document.getElementById('emailNotifications');
    const pushToggle = document.getElementById('pushNotifications');
    const weeklyToggle = document.getElementById('weeklyReports');
    const monthlyToggle = document.getElementById('monthlyReports');

    if (emailToggle) emailToggle.checked = settings.notifications?.email || false;
    if (pushToggle) pushToggle.checked = settings.notifications?.push || false;
    if (weeklyToggle) weeklyToggle.checked = settings.notifications?.weeklyReport || false;
    if (monthlyToggle) monthlyToggle.checked = settings.notifications?.monthlyReport || false;
}

async function handleSpendingLimitsSubmit(e) {
    e.preventDefault();
    const form = e.target;
    const formData = new FormData(form);

    const settings = {
        monthlyLimit: parseFloat(formData.get('monthlyLimit')),
        alertThreshold: parseInt(formData.get('alertThreshold'))
    };

    try {
        await settingsApi.updateSettings(settings);
        notify.success('Spending limits updated successfully');
    } catch (error) {
        notify.error(error.message || 'Failed to update spending limits');
    }
}

async function handleNotificationChange(e) {
    const setting = e.target.id;
    const value = e.target.checked;

    try {
        await settingsApi.updateSettings({
            notifications: {
                [setting]: value
            }
        });
        notify.success('Notification settings updated');
    } catch (error) {
        // Revert toggle if update fails
        e.target.checked = !value;
        notify.error('Failed to update notification settings');
    }
}

async function handleDataExport() {
    try {
        const data = await settingsApi.exportData();
        const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `expense-tracker-export-${new Date().toISOString().split('T')[0]}.json`;
        a.click();
        window.URL.revokeObjectURL(url);
        notify.success('Data exported successfully');
    } catch (error) {
        notify.error('Failed to export data');
    }
}

async function handleDataImport(e) {
    const file = e.target.files[0];
    if (!file) return;

    const reader = new FileReader();
    reader.onload = async (event) => {
        try {
            const data = JSON.parse(event.target.result);
            await settingsApi.importData(data);
            notify.success('Data imported successfully');
            // Reload settings to reflect imported data
            await loadSettings();
        } catch (error) {
            notify.error('Failed to import data. Please check the file format.');
        }
    };
    reader.readAsText(file);
}

async function handleAccountDeletion() {
    const confirmed = confirm(
        'Are you sure you want to delete your account? This action cannot be undone. ' +
        'All your data will be permanently deleted.'
    );

    if (confirmed) {
        try {
            await settingsApi.deleteAccount();
            notify.success('Account deleted successfully');
            // Redirect to login page
            window.location.href = '/pages/login.html';
        } catch (error) {
            notify.error('Failed to delete account');
        }
    }
}
