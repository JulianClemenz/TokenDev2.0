
/**
 * Obtiene el token de autenticación desde sessionStorage.
 */
function getToken() {
  return sessionStorage.getItem('access_token');
}

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

// --- Lógica de la Página de Estadísticas ---

/**
 * Carga el recuento total de usuarios y procesa el gráfico de edades.
 * Llama a (AdminHandler.GetLogs)
 */
async function loadUserCountAndAges() {
  const countElement = document.getElementById('user_count_number');
  if (!countElement) return;
  countElement.textContent = '...';

  try {
    const response = await fetchApi('/api/admin/stats/users');

    if (response.status === 204) {
      countElement.textContent = '0';
      return;
    }

    if (!response.ok) {
      throw new Error('Error al cargar el conteo de usuarios');
    }

    const data = await response.json();

    countElement.textContent = data.total || 0;

    // Tarjeta 3 (Gráfico de Edades)
    if (data.users && data.users.length > 0) {
      processAgeChart(data.users); // Llamar a la función del gráfico
    }

  } catch (error) {
    console.error('Error fetching user count:', error);
    countElement.textContent = 'Error';
  }
}

/**
 * Carga el ranking de los 3 ejercicios más usados en la segunda tarjeta.
 * Llama a(AdminHandler.GetGlobalStats)
 */
async function loadTopExercises() {
  const tableBody = document.getElementById('top_exercises_table_body');
  if (!tableBody) return;
  tableBody.innerHTML = '<tr><td colspan="4">Cargando ranking...</td></tr>';

  try {
    const response = await fetchApi('/api/admin/stats/exercises');

    if (response.status === 204) { // 204 No Content
      tableBody.innerHTML = '<tr><td colspan="4">No hay datos de ejercicios.</td></tr>';
      return;
    }

    if (!response.ok) {
      throw new Error('Error al cargar el ranking de ejercicios');
    }

    const exercises = await response.json();
    tableBody.innerHTML = '';

    if (exercises && exercises.length > 0) {
      const top3 = exercises.slice(0, 3);

      top3.forEach((exercise, index) => {
        const row = document.createElement('tr');
        row.innerHTML = `
          <th scope="row">${index + 1}</th>
          <td>${exercise.ExcerciseName || 'N/D'}</td>
          <td>${exercise.ExcerciseID || 'N/D'}</td>
          <td>${exercise.Count || 0}</td>
        `;
        tableBody.appendChild(row);
      });
    } else {
      tableBody.innerHTML = '<tr><td colspan="4">No hay datos de ejercicios.</td></tr>';
    }

  } catch (error) {
    console.error('Error fetching top exercises:', error);
    tableBody.innerHTML = `<tr><td colspan="4" class="text-danger">Error: ${error.message}</td></tr>`;
  }
}

/**
 * Calcula la edad a partir de una fecha de nacimiento en string ISO ("YYYY-MM-DD...").
 */
function calculateAge(birthDateString) {
  if (!birthDateString) return null;
  try {
    const birthDate = new Date(birthDateString);
    const today = new Date();
    let age = today.getFullYear() - birthDate.getFullYear();
    const monthDifference = today.getMonth() - birthDate.getMonth();
    if (monthDifference < 0 || (monthDifference === 0 && today.getDate() < birthDate.getDate())) {
      age--;
    }
    return age;
  } catch (e) {
    return null;
  }
}

/**
 * Procesa la lista de usuarios para agruparlos por rango de edad.
 */
function processAgeChart(users) {
  // Define los rangos de edad
  const ageGroups = {
    'Menos de 20': 0,
    '20-29': 0,
    '30-39': 0,
    '40-49': 0,
    '50 o más': 0,
    'N/D': 0, // Para usuarios sin fecha de nacimiento
  };

  users.forEach(user => {
    const age = calculateAge(user.BirthDate); // user.BirthDate viene de GetLogs

    if (age === null || isNaN(age)) {
      ageGroups['N/D']++;
    } else if (age < 20) {
      ageGroups['Menos de 20']++;
    } else if (age <= 29) {
      ageGroups['20-29']++;
    } else if (age <= 39) {
      ageGroups['30-39']++;
    } else if (age <= 49) {
      ageGroups['40-49']++;
    } else {
      ageGroups['50 o más']++;
    }
  });

  const labels = Object.keys(ageGroups);
  const data = Object.values(ageGroups);
  renderAgeChart(labels, data);
}

/**
 * Renderiza el gráfico de barras de edades usando Chart.js en el <canvas>.
 */
function renderAgeChart(labels, data) {
  const ctx = document.getElementById('ageChart');
  if (!ctx) return;

  new Chart(ctx, {
    type: 'bar',
    data: {
      labels: labels,
      datasets: [{
        label: 'Cantidad de Usuarios',
        data: data,
        backgroundColor: 'rgba(54, 162, 235, 0.6)',
        borderColor: 'rgba(54, 162, 235, 1)',
        borderWidth: 1
      }]
    },
    options: {
      responsive: true,
      maintainAspectRatio: true, // Puedes ajustarlo si es necesario
      scales: {
        y: {
          beginAtZero: true,
          ticks: {
            // Forzar que el eje Y solo muestre números enteros
            stepSize: 1
          }
        }
      },
      plugins: {
        legend: {
          display: false // Ocultar la leyenda "Cantidad de Usuarios"
        }
      }
    }
  });
}

// --- Inicialización ---
document.addEventListener('DOMContentLoaded', () => {
  loadUserCountAndAges();
  loadTopExercises();
});