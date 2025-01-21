import { categories as categoriesApi } from './api.js';
import { notify, requireAuth } from './utils.js';

// State management
let categories = [];
let currentCategory = null;

// Default categories (used only for new users)
const defaultCategories = [
    { name: 'Food & Dining', icon: 'fa-utensils', color: '#ef4444', budget: 500 },
    { name: 'Transportation', icon: 'fa-car', color: '#3b82f6', budget: 300 },
    { name: 'Utilities', icon: 'fa-bolt', color: '#f59e0b', budget: 200 },
    { name: 'Entertainment', icon: 'fa-film', color: '#8b5cf6', budget: 150 },
    { name: 'Shopping', icon: 'fa-shopping-bag', color: '#ec4899', budget: 400 },
    { name: 'Healthcare', icon: 'fa-heart', color: '#10b981', budget: 200 }
];

// Initialize when document is ready
document.addEventListener('DOMContentLoaded', () => {
    if (!requireAuth()) return;
    initializeCategories();
});

async function initializeCategories() {
    initializeCategoryForm();
    await loadCategories();
}

function initializeCategoryForm() {
    const addCategoryBtn = document.getElementById('addCategoryBtn');
    const categoryModal = document.getElementById('categoryModal');
    const categoryForm = document.getElementById('categoryForm');
    const cancelCategoryBtn = document.getElementById('cancelCategory');
    const colorPicker = document.getElementById('categoryColor');
    const iconPicker = document.getElementById('categoryIcon');

    // Initialize color picker
    if (colorPicker) {
        const colors = ['#ef4444', '#f59e0b', '#10b981', '#3b82f6', '#8b5cf6', '#ec4899'];
        colorPicker.innerHTML = colors.map(color => `
            <div class="color-option" style="background-color: ${color}" 
                 data-color="${color}" onclick="selectColor('${color}')"></div>
        `).join('');
    }

    // Initialize icon picker
    if (iconPicker) {
        const icons = ['utensils', 'car', 'bolt', 'film', 'shopping-bag', 'heart', 
                      'home', 'plane', 'gift', 'book', 'graduation-cap', 'briefcase'];
        iconPicker.innerHTML = icons.map(icon => `
            <div class="icon-option" data-icon="fa-${icon}" onclick="selectIcon('fa-${icon}')">
                <i class="fas fa-${icon}"></i>
            </div>
        `).join('');
    }

    addCategoryBtn?.addEventListener('click', () => {
        currentCategory = null;
        categoryForm.reset();
        categoryModal.classList.remove('hidden');
    });

    cancelCategoryBtn?.addEventListener('click', () => {
        categoryModal.classList.add('hidden');
    });

    categoryForm?.addEventListener('submit', handleCategorySubmit);
}

async function loadCategories() {
    try {
        const response = await categoriesApi.getAll();
        categories = response.categories;
        
        if (categories.length === 0) {
            // For new users, create default categories
            await Promise.all(defaultCategories.map(category => 
                categoriesApi.create(category)
            ));
            categories = (await categoriesApi.getAll()).categories;
        }
        
        renderCategories();
    } catch (error) {
        notify.error('Failed to load categories');
        console.error(error);
    }
}

function renderCategories() {
    const grid = document.getElementById('categoriesGrid');
    if (!grid) return;

    grid.innerHTML = categories.map(category => `
        <div class="category-card" style="border-color: ${category.color}">
            <div class="category-header" style="background-color: ${category.color}">
                <i class="fas ${category.icon}"></i>
                <h3>${category.name}</h3>
            </div>
            <div class="category-body">
                <div class="budget-info">
                    <p>Monthly Budget</p>
                    <h4>$${category.budget.toFixed(2)}</h4>
                </div>
                <div class="progress-bar">
                    <div class="progress" style="width: ${(category.spent / category.budget * 100) || 0}%; 
                         background-color: ${category.color}"></div>
                </div>
                <p class="spent-info">
                    $${(category.spent || 0).toFixed(2)} spent of $${category.budget.toFixed(2)}
                </p>
            </div>
            <div class="category-actions">
                <button class="icon-btn" onclick="editCategory('${category._id}')">
                    <i class="fas fa-edit"></i>
                </button>
                <button class="icon-btn" onclick="deleteCategory('${category._id}')">
                    <i class="fas fa-trash"></i>
                </button>
            </div>
        </div>
    `).join('');
}

async function handleCategorySubmit(e) {
    e.preventDefault();
    const form = e.target;
    const formData = new FormData(form);

    const categoryData = {
        name: formData.get('name'),
        icon: formData.get('icon'),
        color: formData.get('color'),
        budget: parseFloat(formData.get('budget'))
    };

    try {
        if (currentCategory) {
            await categoriesApi.update(currentCategory._id, categoryData);
            notify.success('Category updated successfully');
        } else {
            await categoriesApi.create(categoryData);
            notify.success('Category added successfully');
        }

        document.getElementById('categoryModal').classList.add('hidden');
        await loadCategories();
    } catch (error) {
        notify.error(error.message || 'Failed to save category');
    }
}

async function editCategory(id) {
    const category = categories.find(c => c._id === id);
    if (!category) return;

    currentCategory = category;
    const form = document.getElementById('categoryForm');
    if (!form) return;

    form.name.value = category.name;
    form.budget.value = category.budget;
    selectColor(category.color);
    selectIcon(category.icon);

    document.getElementById('categoryModal').classList.remove('hidden');
}

async function deleteCategory(id) {
    if (!confirm('Are you sure you want to delete this category? All expenses in this category will be set to "Uncategorized".')) 
        return;

    try {
        await categoriesApi.delete(id);
        notify.success('Category deleted successfully');
        await loadCategories();
    } catch (error) {
        notify.error('Failed to delete category');
    }
}

// Color and Icon selection helpers
function selectColor(color) {
    document.querySelectorAll('.color-option').forEach(opt => 
        opt.classList.toggle('selected', opt.dataset.color === color)
    );
    document.querySelector('input[name="color"]').value = color;
}

function selectIcon(icon) {
    document.querySelectorAll('.icon-option').forEach(opt => 
        opt.classList.toggle('selected', opt.dataset.icon === icon)
    );
    document.querySelector('input[name="icon"]').value = icon;
}

// Export functions that need to be accessed globally
window.editCategory = editCategory;
window.deleteCategory = deleteCategory;
window.selectColor = selectColor;
window.selectIcon = selectIcon;
