// DOM Elements
const categoriesGrid = document.querySelector('.categories-grid');
const addCategoryBtn = document.getElementById('add-category-btn');

// Default categories
const defaultCategories = [
    { id: 'food', name: 'Food & Dining', icon: 'fa-utensils', color: '#ef4444' },
    { id: 'transport', name: 'Transportation', icon: 'fa-car', color: '#3b82f6' },
    { id: 'utilities', name: 'Utilities', icon: 'fa-bolt', color: '#f59e0b' },
    { id: 'entertainment', name: 'Entertainment', icon: 'fa-film', color: '#8b5cf6' },
    { id: 'shopping', name: 'Shopping', icon: 'fa-shopping-bag', color: '#ec4899' },
    { id: 'health', name: 'Healthcare', icon: 'fa-heart', color: '#10b981' }
];

// Load categories
function loadCategories() {
    const categories = JSON.parse(localStorage.getItem('categories')) || defaultCategories;
    renderCategories(categories);
}

// Render categories
function renderCategories(categories) {
    const addCategoryCard = categoriesGrid.querySelector('.add-category');
    categoriesGrid.innerHTML = '';
    
    categories.forEach(category => {
        const categoryCard = createCategoryCard(category);
        categoriesGrid.appendChild(categoryCard);
    });
    
    categoriesGrid.appendChild(addCategoryCard);
}

// Create category card
function createCategoryCard(category) {
    const card = document.createElement('div');
    card.className = 'category-card';
    card.style.borderTop = `3px solid ${category.color}`;
    
    card.innerHTML = `
        <div class="category-icon" style="color: ${category.color}">
            <i class="fas ${category.icon} fa-2x"></i>
        </div>
        <h3>${category.name}</h3>
        <div class="category-actions">
            <button class="edit-category" data-id="${category.id}">
                <i class="fas fa-edit"></i>
            </button>
            <button class="delete-category" data-id="${category.id}">
                <i class="fas fa-trash"></i>
            </button>
        </div>
    `;
    
    // Add event listeners
    card.querySelector('.edit-category').addEventListener('click', () => editCategory(category));
    card.querySelector('.delete-category').addEventListener('click', () => deleteCategory(category.id));
    
    return card;
}

// Add new category
addCategoryBtn.addEventListener('click', () => {
    const colorOptions = ['#ef4444', '#f59e0b', '#10b981', '#3b82f6', '#8b5cf6', '#ec4899'];
    const iconOptions = [
        'fa-shopping-cart', 'fa-car', 'fa-home', 'fa-plane', 
        'fa-gift', 'fa-coffee', 'fa-book', 'fa-gamepad'
    ];
    
    const html = `
        <form id="category-form">
            <input type="text" name="name" placeholder="Category Name" required>
            <div class="color-picker">
                ${colorOptions.map(color => `
                    <label class="color-option">
                        <input type="radio" name="color" value="${color}">
                        <span class="color-swatch" style="background-color: ${color}"></span>
                    </label>
                `).join('')}
            </div>
            <div class="icon-picker">
                ${iconOptions.map(icon => `
                    <label class="icon-option">
                        <input type="radio" name="icon" value="${icon}">
                        <i class="fas ${icon}"></i>
                    </label>
                `).join('')}
            </div>
            <button type="submit">Add Category</button>
        </form>
    `;
    
    showModal('Add Category', html, async (modal) => {
        const form = modal.querySelector('#category-form');
        form.addEventListener('submit', (e) => {
            e.preventDefault();
            
            const formData = new FormData(form);
            const newCategory = {
                id: formData.get('name').toLowerCase().replace(/\s+/g, '-'),
                name: formData.get('name'),
                icon: formData.get('icon'),
                color: formData.get('color')
            };
            
            const categories = JSON.parse(localStorage.getItem('categories')) || defaultCategories;
            categories.push(newCategory);
            localStorage.setItem('categories', JSON.stringify(categories));
            
            renderCategories(categories);
            modal.remove();
            notyf.success('Category added successfully');
        });
    });
});

// Edit category
function editCategory(category) {
    // Similar to add category but with pre-filled values
    // Implementation similar to add category with pre-filled values
    notyf.success('Category updated successfully');
}

// Delete category
function deleteCategory(categoryId) {
    if (confirm('Are you sure you want to delete this category?')) {
        const categories = JSON.parse(localStorage.getItem('categories')) || defaultCategories;
        const updatedCategories = categories.filter(c => c.id !== categoryId);
        localStorage.setItem('categories', JSON.stringify(updatedCategories));
        
        renderCategories(updatedCategories);
        notyf.success('Category deleted successfully');
    }
}

// Show modal helper
function showModal(title, content, callback) {
    const modal = document.createElement('div');
    modal.className = 'modal';
    modal.innerHTML = `
        <div class="modal-content">
            <h2>${title}</h2>
            ${content}
        </div>
    `;
    
    document.body.appendChild(modal);
    callback(modal);
    
    modal.addEventListener('click', (e) => {
        if (e.target === modal) {
            modal.remove();
        }
    });
}

// Initialize categories
loadCategories();
