// Utility functions for the application

// Format currency based on user's locale and currency preference
export function formatCurrency(amount, currency = 'USD') {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: currency
    }).format(amount);
}

// Format date to local string
export function formatDate(date) {
    return new Date(date).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
    });
}

// Show notification using Notyf
export const notify = {
    success(message) {
        const notyf = new Notyf({
            duration: 3000,
            position: { x: 'right', y: 'top' }
        });
        notyf.success(message);
    },
    error(message) {
        const notyf = new Notyf({
            duration: 3000,
            position: { x: 'right', y: 'top' }
        });
        notyf.error(message);
    }
};

// Check if user is authenticated
export function isAuthenticated() {
    const token = localStorage.getItem('token');
    return !!token;
}

// Redirect if not authenticated
export function requireAuth() {
    if (!isAuthenticated()) {
        window.location.href = '/pages/login.html';
        return false;
    }
    return true;
}

// Handle API errors
export function handleError(error) {
    console.error('Error:', error);
    notify.error(error.message || 'An unexpected error occurred');
}

// Theme management
export const theme = {
    toggle() {
        const isDark = localStorage.getItem('darkMode') === 'true';
        localStorage.setItem('darkMode', !isDark);
        this.apply();
    },

    apply() {
        const isDark = localStorage.getItem('darkMode') === 'true';
        document.body.classList.toggle('dark-theme', isDark);
        const themeIcon = document.querySelector('#theme-toggle i');
        if (themeIcon) {
            themeIcon.className = isDark ? 'fas fa-sun' : 'fas fa-moon';
        }
    },

    init() {
        // Apply theme on page load
        this.apply();
        // Add event listener to theme toggle button
        const themeToggle = document.getElementById('theme-toggle');
        if (themeToggle) {
            themeToggle.addEventListener('click', () => this.toggle());
        }
    }
};

// Password strength checker
export function checkPasswordStrength(password) {
    let strength = 0;
    const feedback = [];

    // Length check
    if (password.length < 8) {
        feedback.push('Password should be at least 8 characters long');
    } else {
        strength += 1;
    }

    // Uppercase check
    if (!/[A-Z]/.test(password)) {
        feedback.push('Add uppercase letters');
    } else {
        strength += 1;
    }

    // Lowercase check
    if (!/[a-z]/.test(password)) {
        feedback.push('Add lowercase letters');
    } else {
        strength += 1;
    }

    // Number check
    if (!/\d/.test(password)) {
        feedback.push('Add numbers');
    } else {
        strength += 1;
    }

    // Special character check
    if (!/[!@#$%^&*(),.?":{}|<>]/.test(password)) {
        feedback.push('Add special characters');
    } else {
        strength += 1;
    }

    let strengthText = '';
    if (strength < 2) strengthText = 'weak';
    else if (strength < 4) strengthText = 'fair';
    else strengthText = 'strong';

    return {
        score: strength,
        feedback,
        strengthText
    };
}

// Form validation
export const validate = {
    email(email) {
        const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return re.test(email);
    },

    password(password) {
        return password.length >= 8;
    },

    confirmPassword(password, confirmPassword) {
        return password === confirmPassword;
    },

    required(value) {
        return value.trim().length > 0;
    }
};

// Local storage wrapper
export const storage = {
    set(key, value) {
        localStorage.setItem(key, JSON.stringify(value));
    },

    get(key) {
        const item = localStorage.getItem(key);
        try {
            return JSON.parse(item);
        } catch {
            return item;
        }
    },

    remove(key) {
        localStorage.removeItem(key);
    },

    clear() {
        localStorage.clear();
    }
};

// Chart utilities
export const chartUtils = {
    colors: [
        '#4361ee',
        '#3a0ca3',
        '#7209b7',
        '#f72585',
        '#4cc9f0',
        '#4895ef',
        '#560bad',
        '#480ca8',
        '#3f37c9',
        '#4361ee'
    ],

    createChart(ctx, type, data, options = {}) {
        return new Chart(ctx, {
            type,
            data,
            options: {
                responsive: true,
                maintainAspectRatio: false,
                ...options
            }
        });
    }
};
