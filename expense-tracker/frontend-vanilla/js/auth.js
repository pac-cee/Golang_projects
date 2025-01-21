// DOM Elements
const authSection = document.getElementById('auth-section');
const dashboardSection = document.getElementById('dashboard-section');
const loginForm = document.getElementById('login-form');
const registerForm = document.getElementById('register-form');
const tabBtns = document.querySelectorAll('.tab-btn');
const authTabs = document.querySelector('.auth-tabs');
const togglePasswordBtns = document.querySelectorAll('.toggle-password');

// API endpoints
const API_URL = 'http://localhost:8080/api';

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

// Toggle password visibility
togglePasswordBtns.forEach(btn => {
    btn.addEventListener('click', () => {
        const input = btn.parentElement.querySelector('input');
        const icon = btn.querySelector('i');
        
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

// Password strength meter
function updatePasswordStrength(password) {
    const strengthMeter = document.querySelector('.strength-meter');
    const strengthText = document.querySelector('.strength-text');
    
    let score = 0;
    let feedback = [];
    
    // Check patterns
    if (password.length >= passwordStrengthConfig.minLength) score += 20;
    if (passwordStrengthConfig.patterns.hasNumber.test(password)) score += 20;
    if (passwordStrengthConfig.patterns.hasLetter.test(password)) score += 20;
    if (passwordStrengthConfig.patterns.hasSpecial.test(password)) score += 20;
    if (passwordStrengthConfig.patterns.hasUppercase.test(password) && 
        passwordStrengthConfig.patterns.hasLowercase.test(password)) score += 20;
    
    // Update meter
    strengthMeter.style.setProperty('--strength', `${score}%`);
    
    // Update feedback
    if (score < 40) {
        strengthText.textContent = 'Weak password';
        strengthText.style.color = 'var(--danger-color)';
    } else if (score < 80) {
        strengthText.textContent = 'Medium password';
        strengthText.style.color = 'var(--warning-color)';
    } else {
        strengthText.textContent = 'Strong password';
        strengthText.style.color = 'var(--success-color)';
    }
}

// Password input handler
const passwordInput = registerForm.querySelector('input[type="password"]');
passwordInput.addEventListener('input', (e) => {
    updatePasswordStrength(e.target.value);
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

// Login form submission
loginForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    startLoading(loginForm);
    
    const email = loginForm.querySelector('input[type="email"]').value;
    const password = loginForm.querySelector('input[type="password"]').value;
    const rememberMe = loginForm.querySelector('#remember-me').checked;

    try {
        const response = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password, rememberMe }),
        });

        if (!response.ok) {
            throw new Error('Login failed');
        }

        const data = await response.json();
        
        if (rememberMe) {
            localStorage.setItem('token', data.token);
            localStorage.setItem('user', JSON.stringify(data.user));
        } else {
            sessionStorage.setItem('token', data.token);
            sessionStorage.setItem('user', JSON.stringify(data.user));
        }
        
        // Show success message
        notyf.success('Login successful!');
        
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
        notyf.error('Invalid email or password');
    } finally {
        stopLoading(loginForm);
    }
});

// Register form submission
registerForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    startLoading(registerForm);
    
    const name = registerForm.querySelector('input[type="text"]').value;
    const email = registerForm.querySelector('input[type="email"]').value;
    const password = registerForm.querySelectorAll('input[type="password"]')[0].value;
    const confirmPassword = registerForm.querySelectorAll('input[type="password"]')[1].value;

    if (password !== confirmPassword) {
        notyf.error('Passwords do not match');
        stopLoading(registerForm);
        return;
    }

    try {
        const response = await fetch(`${API_URL}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name, email, password }),
        });

        if (!response.ok) {
            throw new Error('Registration failed');
        }

        // Show success message and switch to login
        notyf.success('Registration successful! Please login.');
        registerForm.reset();
        tabBtns[0].click();
    } catch (error) {
        notyf.error('Registration failed. Please try again.');
    } finally {
        stopLoading(registerForm);
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
                const response = await fetch(`${API_URL}/auth/${provider.toLowerCase()}/callback`, {
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
