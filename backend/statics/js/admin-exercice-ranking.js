
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
    // Token inválido o expirado, redirigir al login
    alert('Tu sesión ha expirado. Por favor, inicia sesión de nuevo.');
    window.location.href = '/login';
    throw new Error('No autorizado');
  }

  return response;
}

// --- Lógica de la Página de Ranking ---

/**
 * Carga el ranking COMPLETO de ejercicios.
 * Llama a(AdminHandler.GetGlobalStats)
 */
async function loadFullRanking() {
  const tableBody = document.getElementById('ranking_table_body');
  if (!tableBody) return;
  tableBody.innerHTML = '<tr><td colspan="4">Cargando ranking completo...</td></tr>';

  try {
    const response = await fetchApi('/api/admin/stats/exercises');

    if (response.status === 204) { // 204 No Content
      tableBody.innerHTML = '<tr><td colspan="4">No hay datos de ejercicios.</td></tr>';
      return;
    }

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Error al cargar el ranking de ejercicios');
    }

    const exercises = await response.json();
    tableBody.innerHTML = '';

    if (exercises && exercises.length > 0) {
      exercises.forEach((exercise, index) => {
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
    console.error('Error fetching full exercise ranking:', error);
    tableBody.innerHTML = `<tr><td colspan="4" class="text-danger">Error: ${error.message}</td></tr>`;
  }
}

document.addEventListener('DOMContentLoaded', () => {
  loadFullRanking();
});