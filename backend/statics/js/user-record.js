
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

// --- Lógica de la Página de Historial ---

/**
 * Carga el historial de entrenamientos (workouts) del usuario.
 */
async function loadRecords() {
  const tableBody = document.getElementById('record-table-body');
  const errorElement = document.getElementById('error_msg');
  errorElement.textContent = '';
  tableBody.innerHTML = '<tr><td colspan="3">Cargando historial...</td></tr>';

  try {
    // 1. Llama al endpoint GetWorkouts
    const response = await fetchApi('/api/workouts');

    if (!response.ok) {
      if (response.status === 404) {
        tableBody.innerHTML = '<tr><td colspan="3">Aún no has registrado ningún entrenamiento.</td></tr>';
        return;
      }
      const err = await response.json();
      throw new Error(err.error || 'No se pudo cargar el historial.');
    }

    const records = await response.json(); // Array de WorkoutResponseDTO

    // 2. Renderizar la tabla
    tableBody.innerHTML = '';
    if (records && records.length > 0) {
      records.forEach(record => {
        const row = document.createElement('tr');

        const date = new Date(record.DoneAt);
        const formattedDate = date.toLocaleString('es-ES', {
          day: '2-digit',
          month: '2-digit',
          year: 'numeric',
          hour: '2-digit',
          minute: '2-digit'
        });

        row.innerHTML = `
          <td>${record.RoutineName || 'Rutina eliminada'}</td>
          <td>${formattedDate}</td>
          <td>
            <button 
              type="button" 
              class="btn btn-outline-danger btn-sm btn-delete-record" 
              data-id="${record.id}">
              Eliminar
            </button>
          </td>
        `;
        tableBody.appendChild(row);
      });
    } else {
      tableBody.innerHTML = '<tr><td colspan="3">No hay registros de entrenamientos.</td></tr>';
    }

  } catch (error) {
    console.error('Error al cargar historial:', error);
    errorElement.textContent = error.message;
    tableBody.innerHTML = `<tr><td colspan="3" class="text-danger">Error al cargar.</td></tr>`;
  }
}

/**
 * Maneja el clic en el botón de eliminar un registro de workout.
 * @param {string} workoutId - El ID del workout a eliminar.
 */
async function handleDeleteRecord(workoutId) {
  if (!confirm('¿Estás seguro de que deseas eliminar este registro del historial?')) {
    return;
  }

  const errorElement = document.getElementById('error_msg');
  errorElement.textContent = '';

  try {
    // Llama al endpoint DeleteWorkout
    const response = await fetchApi(`/api/workouts/${workoutId}`, {
      method: 'DELETE'
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Error al eliminar el registro.');
    }

    alert('Registro eliminado correctamente.');
    loadRecords(); // Recargar la tabla

  } catch (error) {
    console.error('Error al eliminar registro:', error);
    errorElement.textContent = `Error: ${error.message}`;
  }
}

// --- Inicialización ---
document.addEventListener('DOMContentLoaded', () => {
  //Cargar los registros al iniciar
  loadRecords();

  //Usar delegación de eventos para los botones de eliminar
  const tableBody = document.getElementById('record-table-body');
  tableBody.addEventListener('click', (event) => {
    const deleteButton = event.target.closest('.btn-delete-record');

    if (deleteButton) {
      const workoutId = deleteButton.dataset.id;
      if (workoutId) {
        handleDeleteRecord(workoutId);
      } else {
        document.getElementById('error_msg').textContent = 'Error: No se pudo encontrar el ID del registro.';
      }
    }
  });
});