// Theme handling
const themeButtons = document.querySelectorAll('.theme-btn');
const body = document.body;

// Initialize theme
function initTheme() {
    const savedTheme = localStorage.getItem('theme') || 'system';
    setTheme(savedTheme);
    themeButtons.forEach(btn => {
        if (btn.dataset.theme === savedTheme) {
            btn.classList.add('active');
        }
    });
}

// Set theme
function setTheme(theme) {
    localStorage.setItem('theme', theme);
    
    if (theme === 'system') {
        // Check system preference
        if (window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches) {
            body.className = 'dark-theme';
        } else {
            body.className = 'light-theme';
        }
    } else {
        body.className = `${theme}-theme`;
    }

    // Update active button
    themeButtons.forEach(btn => btn.classList.remove('active'));
    document.querySelector(`[data-theme="${theme}"]`).classList.add('active');
}

// Theme button click handlers
themeButtons.forEach(btn => {
    btn.addEventListener('click', () => {
        const theme = btn.dataset.theme;
        setTheme(theme);
    });
});

// Listen for system theme changes
if (window.matchMedia) {
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', e => {
        if (localStorage.getItem('theme') === 'system') {
            setTheme('system');
        }
    });
}

// Initialize theme on page load
initTheme();
