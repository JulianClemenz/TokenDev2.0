
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

// --- Lógica de la Página ---

/**
 * Carga los datos de la rutina en el formulario.
 */
async function loadRoutineData(id) {
  const routineNameInput = document.getElementById('routine_name');
  const errorElement = document.getElementById('error_msg');

  try {
    const response = await fetchApi(`/api/routines/${id}`);
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'No se pudo cargar la rutina.');
    }

    const routine = await response.json();
    routineNameInput.value = routine.Name; // Rellenar el input con el nombre actual

  } catch (error) {
    console.error('Error al cargar la rutina:', error);
    errorElement.textContent = `Error al cargar datos: ${error.message}`;
    document.getElementById('btn_save_routine').disabled = true; // Deshabilitar guardado si no se carga
  }
}

/**
 * Maneja el clic en el botón de guardar cambios.
 */
async function handleUpdateRoutine(id) {
  const routineNameInput = document.getElementById('routine_name');
  const errorElement = document.getElementById('error_msg');
  const routineName = routineNameInput.value.trim();

  errorElement.textContent = ''; // Limpiar errores

  // 1. Validar
  if (!routineName) {
    errorElement.textContent = 'El nombre de la rutina no puede estar vacío.';
    routineNameInput.focus();
    return;
  }

  const payload = {
    name: routineName
  };

  try {
    // Enviar al endpoint (PUT /api/routines/:id)
    const response = await fetchApi(`/api/routines/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload)
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'No se pudo actualizar la rutina.');
    }

    // Éxito
    alert('¡Rutina actualizada exitosamente!');
    window.location.href = '/user-routines';

  } catch (error) {
    console.error('Error al actualizar rutina:', error);
    errorElement.textContent = error.message;
  }
}

// --- Inicialización ---
document.addEventListener('DOMContentLoaded', () => {
  //Obtener el ID de la rutina de la URL
  const urlParams = new URLSearchParams(window.location.search);
  const routineId = urlParams.get('id');
  const errorElement = document.getElementById('error_msg');

  if (!routineId) {
    errorElement.textContent = 'Error: No se especificó una rutina para editar.';
    return;
  }

  // Cargar los datos de esa rutina
  loadRoutineData(routineId);

  //Asignar el evento al botón de guardar
  const saveButton = document.getElementById('btn_save_routine');
  if (saveButton) {
    saveButton.addEventListener('click', () => handleUpdateRoutine(routineId));
  }
});