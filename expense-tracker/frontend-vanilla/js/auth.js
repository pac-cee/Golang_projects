// DOM Elements
const authSection = document.getElementById('auth-section');
const dashboardSection = document.getElementById('dashboard-section');
const loginForm = document.getElementById('login-form');
const registerForm = document.getElementById('register-form');
const tabBtns = document.querySelectorAll('.tab-btn');
const authTabs = document.querySelector('.auth-tabs');
const togglePasswordBtns = document.querySelectorAll('.toggle-password');

// Password strength configuration
const passwordStrengthConfig = {
    minLength: 8,
    patterns: {
        hasNumber: /\d/,
        hasLetter: /[a-zA-Z]/,
        hasSpecial: /[!@#$%^&*(),.?":{}|<>]/,
        hasUppercase: /[A-Z]/,
        hasLowercase: /[a-z]/
    }
};

// Password validation rules
const passwordRules = {
    minLength: 8,
    maxLength: 32,
    requireUppercase: true,
    requireLowercase: true,
    requireNumber: true,
    requireSpecial: true,
    allowedSpecial: "!@#$%^&*(),.?\":{}|<>",
    recommendedLength: 12,
    uniqueCharsMin: 6,
    noCommonPatterns: true,
    commonPatterns: [
        '123', '234', '345', '456', '567', '678', '789', '987', '876', '765', '654', '543', '432', '321',
        'password', 'qwerty', 'asdfgh', 'zxcvbn'
    ]
};

function calculatePasswordStrength(password) {
    let score = 0;
    const maxScore = 100;
    
    // Length score (30%)
    const lengthScore = Math.min((password.length - passwordRules.minLength) / 
        (passwordRules.recommendedLength - passwordRules.minLength), 1) * 30;
    score += lengthScore;
    
    // Character variety score (40%)
    const hasUpper = /[A-Z]/.test(password);
    const hasLower = /[a-z]/.test(password);
    const hasNumber = /\d/.test(password);
    const hasSpecial = new RegExp(`[${passwordRules.allowedSpecial}]`).test(password);
    
    score += hasUpper ? 10 : 0;
    score += hasLower ? 10 : 0;
    score += hasNumber ? 10 : 0;
    score += hasSpecial ? 10 : 0;
    
    // Unique characters score (15%)
    const uniqueChars = new Set(password).size;
    const uniqueScore = Math.min((uniqueChars - passwordRules.uniqueCharsMin) / 
        (passwordRules.recommendedLength - passwordRules.uniqueCharsMin), 1) * 15;
    score += uniqueScore;
    
    // Pattern penalty (15%)
    let patternPenalty = 0;
    if (passwordRules.noCommonPatterns) {
        for (const pattern of passwordRules.commonPatterns) {
            if (password.toLowerCase().includes(pattern)) {
                patternPenalty += 5;
            }
        }
    }
    score = Math.max(0, score - patternPenalty);
    
    // Additional bonuses
    if (password.length >= passwordRules.recommendedLength && 
        hasUpper && hasLower && hasNumber && hasSpecial && 
        uniqueChars >= passwordRules.recommendedLength) {
        score += 10; // Bonus for excellent password
    }
    
    return Math.min(Math.round(score), maxScore);
}

function getPasswordFeedback(password) {
    const feedback = [];
    const strength = calculatePasswordStrength(password);
    
    // Basic requirements
    if (password.length < passwordRules.minLength) {
        feedback.push({
            type: 'error',
            message: `Add ${passwordRules.minLength - password.length} more characters`
        });
    }
    
    if (!/[A-Z]/.test(password)) {
        feedback.push({
            type: 'error',
            message: 'Add an uppercase letter'
        });
    }
    
    if (!/[a-z]/.test(password)) {
        feedback.push({
            type: 'error',
            message: 'Add a lowercase letter'
        });
    }
    
    if (!/\d/.test(password)) {
        feedback.push({
            type: 'error',
            message: 'Add a number'
        });
    }
    
    if (!new RegExp(`[${passwordRules.allowedSpecial}]`).test(password)) {
        feedback.push({
            type: 'error',
            message: 'Add a special character'
        });
    }
    
    // Improvement suggestions
    if (password.length < passwordRules.recommendedLength) {
        feedback.push({
            type: 'suggestion',
            message: `Adding ${passwordRules.recommendedLength - password.length} more characters would make this stronger`
        });
    }
    
    const uniqueChars = new Set(password).size;
    if (uniqueChars < passwordRules.uniqueCharsMin) {
        feedback.push({
            type: 'suggestion',
            message: 'Try using more unique characters'
        });
    }
    
    // Check for common patterns
    for (const pattern of passwordRules.commonPatterns) {
        if (password.toLowerCase().includes(pattern)) {
            feedback.push({
                type: 'warning',
                message: 'Avoid using common patterns'
            });
            break;
        }
    }
    
    return {
        score: strength,
        feedback: feedback,
        strengthText: strength < 40 ? 'Weak' : 
                     strength < 60 ? 'Fair' :
                     strength < 80 ? 'Good' :
                     strength < 90 ? 'Strong' : 'Excellent'
    };
}

function validatePassword(password) {
    const errors = [];
    
    if (password.length < passwordRules.minLength) {
        errors.push(`Password must be at least ${passwordRules.minLength} characters long`);
    }
    
    if (password.length > passwordRules.maxLength) {
        errors.push(`Password must be less than ${passwordRules.maxLength} characters`);
    }
    
    if (passwordRules.requireUppercase && !/[A-Z]/.test(password)) {
        errors.push('Password must contain at least one uppercase letter');
    }
    
    if (passwordRules.requireLowercase && !/[a-z]/.test(password)) {
        errors.push('Password must contain at least one lowercase letter');
    }
    
    if (passwordRules.requireNumber && !/\d/.test(password)) {
        errors.push('Password must contain at least one number');
    }
    
    if (passwordRules.requireSpecial && !new RegExp(`[${passwordRules.allowedSpecial}]`).test(password)) {
        errors.push('Password must contain at least one special character (!@#$%^&*(),.?":{}|<>)');
    }
    
    return {
        isValid: errors.length === 0,
        errors: errors
    };
}

// Tab switching
tabBtns.forEach(btn => {
    btn.addEventListener('click', () => {
        const targetTab = btn.dataset.tab;
        
        // Update active tab button
        tabBtns.forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        
        // Update tab indicator
        authTabs.dataset.active = targetTab;
        
        // Show/hide forms with fade animation
        if (targetTab === 'login') {
            registerForm.classList.add('fade-out');
            setTimeout(() => {
                registerForm.classList.add('hidden');
                loginForm.classList.remove('hidden');
                loginForm.classList.add('fade-in');
            }, 200);
        } else {
            loginForm.classList.add('fade-out');
            setTimeout(() => {
                loginForm.classList.add('hidden');
                registerForm.classList.remove('hidden');
                registerForm.classList.add('fade-in');
            }, 200);
        }
    });
});

// Password visibility toggle
document.querySelectorAll('.password-toggle').forEach(button => {
    button.addEventListener('click', (e) => {
        const input = e.currentTarget.previousElementSibling;
        const icon = e.currentTarget.querySelector('i');
        
        if (input.type === 'password') {
            input.type = 'text';
            icon.classList.remove('fa-eye');
            icon.classList.add('fa-eye-slash');
        } else {
            input.type = 'password';
            icon.classList.remove('fa-eye-slash');
            icon.classList.add('fa-eye');
        }
    });
});

// Password strength update
const registerPassword = document.getElementById('register-password');
if (registerPassword) {
    registerPassword.addEventListener('input', (e) => {
        updatePasswordStrength(e.target.value);
    });
}

// Update password strength meter
function updatePasswordStrength(password) {
    const strengthMeter = document.querySelector('.strength-meter');
    const strengthText = document.querySelector('.strength-text');
    const feedbackList = document.querySelector('.password-feedback');
    
    if (!strengthMeter || !strengthText) return;
    
    const result = getPasswordFeedback(password);
    
    // Update strength meter
    strengthMeter.style.setProperty('--strength', `${result.score}%`);
    strengthMeter.className = `strength-meter ${result.strengthText.toLowerCase()}`;
    
    // Update strength text
    strengthText.textContent = result.strengthText;
    strengthText.className = `strength-text ${result.strengthText.toLowerCase()}`;
    
    // Update feedback list
    if (feedbackList) {
        feedbackList.innerHTML = '';
        result.feedback.forEach(item => {
            const li = document.createElement('li');
            li.className = `feedback-item ${item.type}`;
            li.innerHTML = `
                <i class="fas fa-${item.type === 'error' ? 'times' : 
                                  item.type === 'warning' ? 'exclamation' : 
                                  'info'}-circle"></i>
                ${item.message}
            `;
            feedbackList.appendChild(li);
        });
    }
}

// Password input handler
const passwordInput = registerForm.querySelector('input[type="password"]');
passwordInput.addEventListener('input', (e) => {
    updatePasswordStrength(e.target.value);
});

// Update password input validation
const passwordInputs = document.querySelectorAll('input[type="password"]');
passwordInputs.forEach(input => {
    const errorDiv = document.createElement('div');
    errorDiv.className = 'password-error';
    input.parentElement.appendChild(errorDiv);
    
    input.addEventListener('input', (e) => {
        const result = validatePassword(e.target.value);
        const errorDiv = e.target.parentElement.querySelector('.password-error');
        
        if (!result.isValid) {
            errorDiv.textContent = result.errors[0];
            errorDiv.style.display = 'block';
            e.target.setCustomValidity(result.errors[0]);
        } else {
            errorDiv.style.display = 'none';
            e.target.setCustomValidity('');
        }
        
        // Update password strength if this is a registration password
        if (e.target.closest('#register-form')) {
            updatePasswordStrength(e.target.value);
        }
    });
});

// Form submission animations
function startLoading(form) {
    const button = form.querySelector('.auth-button');
    button.classList.add('loading');
    button.disabled = true;
}

function stopLoading(form) {
    const button = form.querySelector('.auth-button');
    button.classList.remove('loading');
    button.disabled = false;
}

// Register form submission
registerForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    startLoading(registerForm);

    const name = document.getElementById('register-name').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;
    const confirmPassword = document.getElementById('register-confirm-password').value;

    // Validate password match
    if (password !== confirmPassword) {
        notyf.error('Passwords do not match');
        stopLoading(registerForm);
        return;
    }

    try {
        const response = await fetch(`${window.location.origin}/api/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                name,
                email,
                password
            })
        });

        const data = await response.text(); // First get the response as text
        let jsonData;
        try {
            jsonData = JSON.parse(data); // Try to parse it as JSON
        } catch (err) {
            console.error('Server response:', data); // Log the actual response
            throw new Error('Invalid server response');
        }

        if (!response.ok) {
            throw new Error(jsonData.error || 'Registration failed');
        }

        // Registration successful
        notyf.success('Registration successful! Please login.');
        switchTab('login'); // Switch to login tab
        registerForm.reset();
    } catch (error) {
        console.error('Registration error:', error);
        notyf.error(error.message || 'Registration failed. Please try again.');
    } finally {
        stopLoading(registerForm);
    }
});

// Login form submission
loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    startLoading(loginForm);

    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;
    const rememberMe = document.getElementById('remember-me').checked;

    try {
        const response = await fetch(`${window.location.origin}/api/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                email,
                password
            })
        });

        const data = await response.text(); // First get the response as text
        let jsonData;
        try {
            jsonData = JSON.parse(data); // Try to parse it as JSON
        } catch (err) {
            console.error('Server response:', data); // Log the actual response
            throw new Error('Invalid server response');
        }

        if (!response.ok) {
            throw new Error(jsonData.error || 'Login failed');
        }

        // Store the token
        const storage = rememberMe ? localStorage : sessionStorage;
        storage.setItem('token', jsonData.token);
        storage.setItem('user', JSON.stringify(jsonData.user));

        // Update UI
        authSection.style.display = 'none';
        dashboardSection.style.display = 'block';
        const dashboardLink = document.getElementById('dashboard-link');
        if (dashboardLink) {
            dashboardLink.style.display = 'block';
        }

        notyf.success('Login successful!');
        loginForm.reset();
    } catch (error) {
        console.error('Login error:', error);
        notyf.error(error.message || 'Login failed. Please try again.');
    } finally {
        stopLoading(loginForm);
    }
});

// Logout functionality
document.querySelector('.logout').addEventListener('click', () => {
    // Clear storage
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    sessionStorage.removeItem('token');
    sessionStorage.removeItem('user');
    
    // Animate transition
    dashboardSection.classList.add('fade-out');
    setTimeout(() => {
        dashboardSection.classList.add('hidden');
        authSection.classList.remove('hidden');
        authSection.classList.add('fade-in');
    }, 300);
    
    notyf.success('Logged out successfully');
});

// Check authentication status on page load
document.addEventListener('DOMContentLoaded', () => {
    const token = localStorage.getItem('token');
    const dashboardLink = document.getElementById('dashboard-link');
    
    if (token) {
        authSection.style.display = 'none';
        dashboardSection.style.display = 'block';
        if (dashboardLink) {
            dashboardLink.style.display = 'block';
        }
    } else {
        authSection.style.display = 'block';
        dashboardSection.style.display = 'none';
        if (dashboardLink) {
            dashboardLink.style.display = 'none';
        }
    }
});

// Check if user is already logged in
window.addEventListener('load', () => {
    const token = localStorage.getItem('token') || sessionStorage.getItem('token');
    if (token) {
        authSection.classList.add('hidden');
        dashboardSection.classList.remove('hidden');
        const user = JSON.parse(localStorage.getItem('user') || sessionStorage.getItem('user'));
        document.getElementById('user-name').textContent = user.name;
        loadDashboardData();
    }
});

// Social Login Configuration
const socialConfig = {
    google: {
        client_id: 'YOUR_GOOGLE_CLIENT_ID',
        redirect_uri: `${window.location.origin}/auth/google/callback`,
        scope: 'email profile'
    },
    github: {
        client_id: 'YOUR_GITHUB_CLIENT_ID',
        redirect_uri: `${window.location.origin}/auth/github/callback`,
        scope: 'user:email'
    },
    twitter: {
        client_id: 'YOUR_TWITTER_CLIENT_ID',
        redirect_uri: `${window.location.origin}/auth/twitter/callback`
    }
};

// Social Login Handlers
function initSocialLogin() {
    // Google Login
    document.querySelectorAll('#google-login, #google-signup').forEach(btn => {
        btn.addEventListener('click', () => {
            const url = `https://accounts.google.com/o/oauth2/v2/auth?` +
                `client_id=${socialConfig.google.client_id}&` +
                `redirect_uri=${encodeURIComponent(socialConfig.google.redirect_uri)}&` +
                `response_type=code&` +
                `scope=${encodeURIComponent(socialConfig.google.scope)}&` +
                `access_type=offline`;
            
            openAuthWindow('Google', url);
        });
    });

    // GitHub Login
    document.querySelectorAll('#github-login, #github-signup').forEach(btn => {
        btn.addEventListener('click', () => {
            const url = `https://github.com/login/oauth/authorize?` +
                `client_id=${socialConfig.github.client_id}&` +
                `redirect_uri=${encodeURIComponent(socialConfig.github.redirect_uri)}&` +
                `scope=${encodeURIComponent(socialConfig.github.scope)}`;
            
            openAuthWindow('GitHub', url);
        });
    });

    // Twitter Login
    document.querySelectorAll('#twitter-login, #twitter-signup').forEach(btn => {
        btn.addEventListener('click', () => {
            const url = `https://twitter.com/i/oauth2/authorize?` +
                `client_id=${socialConfig.twitter.client_id}&` +
                `redirect_uri=${encodeURIComponent(socialConfig.twitter.redirect_uri)}&` +
                `response_type=code&` +
                `scope=tweet.read%20users.read`;
            
            openAuthWindow('Twitter', url);
        });
    });
}

// Open OAuth Window
function openAuthWindow(provider, url) {
    const width = 500;
    const height = 600;
    const left = window.screen.width / 2 - width / 2;
    const top = window.screen.height / 2 - height / 2;
    
    const authWindow = window.open(
        url,
        `${provider} Login`,
        `width=${width},height=${height},left=${left},top=${top},toolbar=0,scrollbars=0,status=0,resizable=0,location=0,menuBar=0`
    );
    
    // Handle OAuth callback
    window.addEventListener('message', async (event) => {
        if (event.origin !== window.location.origin) return;
        
        if (event.data.type === 'social-auth') {
            authWindow.close();
            
            try {
                const response = await fetch(`${window.location.origin}/api/auth/${provider.toLowerCase()}/callback`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(event.data.code)
                });
                
                if (!response.ok) throw new Error('Authentication failed');
                
                const data = await response.json();
                localStorage.setItem('token', data.token);
                localStorage.setItem('user', JSON.stringify(data.user));
                
                notyf.success(`Successfully logged in with ${provider}`);
                
                // Animate transition to dashboard
                authSection.classList.add('fade-out');
                setTimeout(() => {
                    authSection.classList.add('hidden');
                    dashboardSection.classList.remove('hidden');
                    dashboardSection.classList.add('fade-in');
                    
                    // Update user info
                    document.getElementById('user-name').textContent = data.user.name;
                    
                    // Load dashboard data
                    loadDashboardData();
                }, 300);
            } catch (error) {
                notyf.error(`Failed to authenticate with ${provider}`);
            }
        }
    });
}

// Handle OAuth Callback
function handleOAuthCallback() {
    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get('code');
    const provider = urlParams.get('provider');
    
    if (code && provider) {
        window.opener.postMessage({
            type: 'social-auth',
            code,
            provider
        }, window.location.origin);
        window.close();
    }
}

// Initialize social login
initSocialLogin();

// Check if this is an OAuth callback
handleOAuthCallback();

// Initialize Notyf
const notyf = new Notyf({
    duration: 3000,
    position: { x: 'right', y: 'top' },
    types: [
        {
            type: 'success',
            background: 'var(--success-color)',
            icon: {
                className: 'fas fa-check-circle',
                tagName: 'i'
            }
        },
        {
            type: 'error',
            background: 'var(--danger-color)',
            icon: {
                className: 'fas fa-times-circle',
                tagName: 'i'
            }
        }
    ]
});

function switchTab(tab) {
    tabBtns.forEach(btn => btn.classList.remove('active'));
    document.querySelector(`.tab-btn[data-tab="${tab}"]`).classList.add('active');
    authTabs.dataset.active = tab;
    if (tab === 'login') {
        registerForm.classList.add('fade-out');
        setTimeout(() => {
            registerForm.classList.add('hidden');
            loginForm.classList.remove('hidden');
            loginForm.classList.add('fade-in');
        }, 200);
    } else {
        loginForm.classList.add('fade-out');
        setTimeout(() => {
            loginForm.classList.add('hidden');
            registerForm.classList.remove('hidden');
            registerForm.classList.add('fade-in');
        }, 200);
    }
}
