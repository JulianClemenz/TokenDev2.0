
/**
 * Obtiene el token de autenticación desde sessionStorage.
 */
function getToken() {
  return sessionStorage.getItem('access_token');
}

/**
 * Realiza una solicitud fetch autenticada.
 */
async function fetchApi(url, options = {}) {
  const token = getToken();

  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
    ...options.headers,
  };

  const response = await fetch(url, { ...options, headers });

  if (response.status === 401) {
    alert('Tu sesión ha expirado. Por favor, inicia sesión de nuevo.');
    window.location.href = '/login';
    throw new Error('No autorizado');
  }

  return response;
}

// --- Lógica de la Página de Progreso ---

/**
 * Carga las estadísticas del usuario desde /api/workouts/stats
 */
async function loadStats() {
  const errorElement = document.getElementById('error_msg');
  const totalEl = document.getElementById('stats_total_workouts');
  const freqEl = document.getElementById('stats_weekly_frequency');
  const routinesBody = document.getElementById('stats_top_routines_body');

  try {
    errorElement.textContent = '';
    totalEl.textContent = 'Cargando...';
    freqEl.textContent = 'Cargando...';
    routinesBody.innerHTML = '<tr><td colspan="3">Cargando...</td></tr>';

    const response = await fetchApi('/api/workouts/stats');

    if (!response.ok) {
      if (response.status === 404) {
        throw new Error('Aún no tienes suficientes datos para mostrar estadísticas.');
      }
      const err = await response.json();
      throw new Error(err.error || 'No se pudieron cargar las estadísticas');
    }

    const stats = await response.json(); // Esto es tu WorkoutStatsDTO

    // Renderizar Tarjetas Simples
    totalEl.textContent = stats.TotalWorkouts || 0;
    freqEl.textContent = (stats.WeeklyFrequency || 0).toFixed(1);

    //Renderizar Tabla de Rutinas Más Usadas
    renderTopRoutines(stats.MostUsedRoutines, routinesBody);

    // Renderizar Gráfico de Progreso
    renderProgressChart(stats.ProgressOverTime);

  } catch (error) {
    console.error('Error al cargar estadísticas:', error);
    errorElement.textContent = error.message;
    totalEl.textContent = 'Error';
    freqEl.textContent = 'Error';
    routinesBody.innerHTML = `<tr><td colspan="3" class="text-danger">Error</td></tr>`;
  }
}

/**
 * Renderiza la tabla de rutinas más usadas.
 * @param {Array} routines - La lista de RoutineUsageDTO
 * @param {HTMLElement} tableBody - El <tbody> de la tabla
 */
function renderTopRoutines(routines, tableBody) {
  tableBody.innerHTML = '';
  if (routines && routines.length > 0) {
    routines.forEach((routine, index) => {
      const row = document.createElement('tr');
      row.innerHTML = `
        <th scope="row">${index + 1}</th>
        <td>${routine.RoutineName || 'N/D'}</td>
        <td>${routine.Count || 0}</td>
      `;
      tableBody.appendChild(row);
    });
  } else {
    tableBody.innerHTML = '<tr><td colspan="3">Sin datos de rutinas.</td></tr>';
  }
}

/**
 * Renderiza el gráfico de líneas de progreso.
 * @param {Array} progressData - La lista de ProgressPointDTO
 */
function renderProgressChart(progressData) {
  const ctx = document.getElementById('progressChart');
  if (!ctx || !progressData) return;

  const labels = progressData.map(point => point.Date);
  const data = progressData.map(point => point.Count);

  new Chart(ctx, {
    type: 'line',
    data: {
      labels: labels,
      datasets: [{
        label: 'Entrenamientos',
        data: data,
        fill: false,
        borderColor: 'rgb(75, 192, 192)',
        tension: 0.1
      }]
    },
    options: {
      responsive: true,
      scales: {
        y: {
          beginAtZero: true,
          ticks: {
            stepSize: 1
          }
        }
      },
      plugins: {
        legend: {
          display: false
        }
      }
    }
  });
}


// --- Inicialización ---
document.addEventListener('DOMContentLoaded', loadStats);