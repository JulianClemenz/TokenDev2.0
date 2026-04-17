
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

async function handleSaveRoutine() {
  const routineNameInput = document.getElementById('routine_name');
  const errorElement = document.getElementById('error_msg');
  const routineName = routineNameInput.value.trim();

  errorElement.textContent = '';

  // 1. Validar que el nombre no esté vacío
  if (!routineName) {
    errorElement.textContent = 'Por favor, introduce un nombre para la rutina.';
    routineNameInput.focus();
    return;
  }

  const payload = {
    name: routineName
  };

  try {
    // 2. Enviar al endpoint (POST /api/routines)
    const response = await fetchApi('/api/routines', {
      method: 'POST',
      body: JSON.stringify(payload)
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'No se pudo crear la rutina.');
    }

    // 3. Éxito
    alert('¡Rutina creada exitosamente!');

    //redirigimos de vuelta al listado.
    window.location.href = '/user-routines';

  } catch (error) {
    console.error('Error al crear rutina:', error);
    errorElement.textContent = error.message;
  }
}

// --- Inicialización ---
document.addEventListener('DOMContentLoaded', () => {
  const saveButton = document.getElementById('btn_save_routine');
  if (saveButton) {
    saveButton.addEventListener('click', handleSaveRoutine);
  }
});