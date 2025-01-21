// Initialize charts
let trendChart;
let categoryChart;

// Update charts with new data
async function updateCharts() {
    const token = localStorage.getItem('token');
    if (!token) return;

    try {
        // Fetch data for trend chart
        const trendResponse = await fetch(`${API_URL}/expenses/trend`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        const trendData = await trendResponse.json();
        updateTrendChart(trendData);

        // Fetch data for category distribution
        const categoryResponse = await fetch(`${API_URL}/expenses/by-category`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });
        const categoryData = await categoryResponse.json();
        updateCategoryChart(categoryData);
    } catch (error) {
        console.error('Failed to update charts:', error);
    }
}

// Update expense trend chart
function updateTrendChart(data) {
    const ctx = document.getElementById('trend-chart').getContext('2d');
    
    // Destroy existing chart if it exists
    if (trendChart) {
        trendChart.destroy();
    }

    trendChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: data.labels,
            datasets: [{
                label: 'Daily Expenses',
                data: data.values,
                borderColor: '#2563eb',
                backgroundColor: 'rgba(37, 99, 235, 0.1)',
                tension: 0.4,
                fill: true
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    display: false
                }
            },
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        callback: function(value) {
                            return '$' + value;
                        }
                    }
                }
            }
        }
    });
}

// Update category distribution chart
function updateCategoryChart(data) {
    const ctx = document.getElementById('category-chart').getContext('2d');
    
    // Destroy existing chart if it exists
    if (categoryChart) {
        categoryChart.destroy();
    }

    // Color palette for categories
    const colors = [
        '#2563eb', // Primary blue
        '#7c3aed', // Purple
        '#059669', // Green
        '#dc2626', // Red
        '#d97706', // Orange
        '#2563eb', // Blue
        '#db2777', // Pink
    ];

    categoryChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: data.labels,
            datasets: [{
                data: data.values,
                backgroundColor: colors.slice(0, data.labels.length),
                borderWidth: 0
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'right',
                    labels: {
                        usePointStyle: true,
                        padding: 20
                    }
                }
            },
            cutout: '70%'
        }
    });
}

// Initialize charts when page loads
window.addEventListener('load', () => {
    if (document.getElementById('trend-chart') && document.getElementById('category-chart')) {
        updateCharts();
    }
});
