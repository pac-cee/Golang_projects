// Fetch and display the current BTC price
async function fetchPrice() {
  try {
    const res = await fetch('http://localhost:8088/price');
    const data = await res.json();
    document.getElementById('price').textContent = `BTC/USDT: $${data.price.toLocaleString(undefined, {maximumFractionDigits: 2})}`;
  } catch (err) {
    document.getElementById('price').textContent = 'Error loading price';
  }
}

// Fetch and display candlestick chart
async function fetchCandles() {
  try {
    const res = await fetch('http://localhost:8088/candles');
    const candles = await res.json();
    renderCandlestickChart(candles);
  } catch (err) {
    const ctx = document.getElementById('candlestickChart').getContext('2d');
    ctx.font = '16px sans-serif';
    ctx.fillText('Error loading chart', 50, 150);
  }
}

function renderCandlestickChart(candles) {
  // Prepare data for Chart.js Financial plugin (or use OHLC as bar chart for simplicity)
  const labels = candles.map(c => new Date(c.open_time).toLocaleTimeString());
  const data = {
    labels,
    datasets: [
      {
        label: 'BTC/USDT (Close)',
        data: candles.map(c => c.close),
        borderColor: '#f7931a',
        backgroundColor: 'rgba(247,147,26,0.1)',
        tension: 0.1,
      }
    ]
  };
  const config = {
    type: 'line',
    data,
    options: {
      responsive: false,
      plugins: {
        legend: { display: false }
      },
      scales: {
        x: { display: false },
        y: { display: true }
      }
    }
  };
  // Destroy previous chart if exists
  if(window.candleChart) window.candleChart.destroy();
  const ctx = document.getElementById('candlestickChart').getContext('2d');
  window.candleChart = new Chart(ctx, config);
}

fetchPrice();
fetchCandles();
setInterval(fetchPrice, 1000); // Fetch price every 1 second
setInterval(fetchCandles, 5000); // Fetch candles every 5 seconds
