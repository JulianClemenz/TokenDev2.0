
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

/**
 * Renderiza la lista de ejercicios en el cuerpo de la tabla.
 * @param {Array} exercises - La lista de ejercicios a mostrar.
 * @param {HTMLElement} tableBody 
 */
function renderExercises(exercises, tableBody) {
  tableBody.innerHTML = ''; // Limpiar "cargando" o resultados anteriores

  if (exercises && exercises.length > 0) {
    exercises.forEach(exercise => {
      // Asegúrate de que tu DTO 'ExcerciseResponseDTO' incluya 'id'.
      // (Lo añadimos en la respuesta anterior).
      const exerciseId = exercise.id;

      const row = document.createElement('tr');
      row.innerHTML = `
        <td class="d-flex gap-2">
          <a href="/admin-excercise-edit?id=${exerciseId}" class="btn btn-outline-primary btn-sm">Editar</a>
          <button type="button" class="btn btn-outline-danger btn-sm btn-delete-exercise" data-id="${exerciseId}">
            Eliminar
          </button>
        </td>
        <td>${exercise.Name || ''}</td>
        <td>${exercise.Description || ''}</td>
        <td>${exercise.Category || ''}</td>
        <td>${exercise.MainMuscleGroup || ''}</td>
        <td>${exercise.DifficultLevel || ''}</td>
        <td><a href="${exercise.Example || '#'}" target="_blank" rel="noopener">URL</a></td>
      `;
      tableBody.appendChild(row);
    });
  } else {
    tableBody.innerHTML = '<tr><td colspan="7">No se encontraron ejercicios con esos filtros.</td></tr>';
  }
}

/**
 * Carga los ejercicios desde la API (con o sin filtros) y los muestra en la tabla.
 */
async function loadExercises() {
  const tableBody = document.querySelector('.table tbody');
  tableBody.innerHTML = '<tr><td colspan="7">Cargando ejercicios...</td></tr>';

  // Leer valores de los filtros
  const name = document.getElementById('filter_name').value.trim();
  const category = document.getElementById('filter_category').value;
  const muscleGroup = document.getElementById('filter_muscle_group').value.trim();

  let endpoint = '';
  const params = new URLSearchParams();

  // Construir la URL del endpoint
  if (name) params.append('name', name);
  if (category) params.append('category', category);
  if (muscleGroup) params.append('muscle_group', muscleGroup);

  const queryString = params.toString();


  if (queryString) {
    // Usamos el endpoint de filtros
    endpoint = `/api/exercises/filter?${queryString}`;
  } else {
    // Usamos el endpoint que trae todos
    endpoint = '/api/exercises';
  }

  try {
    const response = await fetchApi(endpoint);

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || `Error ${response.status}: No se pudieron cargar los ejercicios.`);
    }

    const exercises = await response.json();
    renderExercises(exercises, tableBody); // Usar la función de renderizado

  } catch (error) {
    console.error('Error al cargar ejercicios:', error);
    // Manejar el error de "debe ingresar al menos un filtro"
    if (error.message.includes("al menos un filtro")) {
      tableBody.innerHTML = '<tr><td colspan="7">No hay ejercicios para mostrar. Limpie los filtros para ver todos.</td></tr>';
    } else {
      tableBody.innerHTML = `<tr><td colspan="7" class="text-danger">Error: ${error.message}</td></tr>`;
    }
  }
}

/**
 * Maneja el clic en el botón de eliminar ejercicio.
 */
async function handleDeleteExercise(exerciseId) {
  if (!confirm('¿Estás seguro de que deseas eliminar este ejercicio? Esta acción no se puede deshacer.')) {
    return;
  }

  try {
    const response = await fetchApi(`/api/exercises/${exerciseId}`, {
      method: 'DELETE'
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Error al eliminar el ejercicio.');
    }

    alert('Ejercicio eliminado correctamente.');
    loadExercises(); // Recargar la tabla para mostrar los cambios

  } catch (error) {
    console.error('Error al eliminar ejercicio:', error);
    alert(`Error: ${error.message}`);
  }
}

document.addEventListener('DOMContentLoaded', () => {
  // Carga la lista inicial (sin filtros)
  loadExercises();

  // Asigna evento al botón de Filtrar
  document.getElementById('btn_filter').addEventListener('click', loadExercises);

  // Asigna evento al botón de Limpiar
  document.getElementById('btn_clear_filters').addEventListener('click', () => {
    document.getElementById('filter_name').value = '';
    document.getElementById('filter_category').value = '';
    document.getElementById('filter_muscle_group').value = '';
    loadExercises(); // Recargar la lista completa
  });

  const tableBody = document.querySelector('.table tbody');
  tableBody.addEventListener('click', (event) => {
    const deleteButton = event.target.closest('.btn-delete-exercise');

    if (deleteButton) {
      const exerciseId = deleteButton.dataset.id;
      handleDeleteExercise(exerciseId);
    }
  });
});