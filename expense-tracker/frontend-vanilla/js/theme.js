// Theme management module
const theme = {
    // Get current theme
    get() {
        return localStorage.getItem('theme') || 'light';
    },

    // Set theme
    set(mode) {
        localStorage.setItem('theme', mode);
        this.apply();
    },

    // Toggle between light and dark
    toggle() {
        const current = this.get();
        this.set(current === 'light' ? 'dark' : 'light');
    },

    // Apply current theme to document
    apply() {
        const isDark = this.get() === 'dark';
        document.documentElement.classList.toggle('dark', isDark);
        
        // Update theme toggle button if it exists
        const themeToggle = document.getElementById('themeToggle');
        if (themeToggle) {
            const icon = themeToggle.querySelector('i');
            if (icon) {
                icon.className = isDark ? 'fas fa-sun' : 'fas fa-moon';
            }
        }

        // Update theme color meta tag
        const metaThemeColor = document.querySelector('meta[name="theme-color"]');
        if (metaThemeColor) {
            metaThemeColor.content = isDark ? '#1a1a1a' : '#ffffff';
        }

        // Dispatch theme change event
        window.dispatchEvent(new CustomEvent('themechange', {
            detail: { theme: isDark ? 'dark' : 'light' }
        }));
    },

    // Initialize theme system
    init() {
        // Apply theme on page load
        this.apply();

        // Add event listener to theme toggle button
        const themeToggle = document.getElementById('themeToggle');
        if (themeToggle) {
            themeToggle.addEventListener('click', () => this.toggle());
        }

        // Listen for system theme changes
        if (window.matchMedia) {
            const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
            mediaQuery.addEventListener('change', e => {
                if (!localStorage.getItem('theme')) {
                    this.set(e.matches ? 'dark' : 'light');
                }
            });
        }

        // Add keyboard shortcut for theme toggle (Ctrl/Cmd + J)
        document.addEventListener('keydown', e => {
            if ((e.ctrlKey || e.metaKey) && e.key === 'j') {
                e.preventDefault();
                this.toggle();
            }
        });
    }
};

// Initialize theme on page load
document.addEventListener('DOMContentLoaded', () => {
    theme.init();
});

export default theme;
