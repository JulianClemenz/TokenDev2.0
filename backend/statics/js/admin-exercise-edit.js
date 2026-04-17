
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

// --- Lógica de la Página de Edición ---

/**
 * Carga los datos de un ejercicio específico en el formulario.
 */
async function loadExerciseData(id) {
  const errorElement = document.getElementById('error_msg');
  try {
    const response = await fetchApi(`/api/exercises/${id}`);

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'No se pudo cargar el ejercicio.');
    }

    const exercise = await response.json();

    // Rellenar el formulario
    document.getElementById('ex_name').value = exercise.Name;
    document.getElementById('ex_group').value = exercise.MainMuscleGroup;
    document.getElementById('ex_desc').value = exercise.Description;
    document.getElementById('ex_category').value = exercise.Category;
    document.getElementById('ex_difficulty').value = exercise.DifficultLevel;
    document.getElementById('ex_sample').value = exercise.Example;
    document.getElementById('ex_instructions').value = exercise.Instructions;

  } catch (error) {
    console.error('Error al cargar ejercicio:', error);
    errorElement.textContent = `Error al cargar datos: ${error.message}`;
    document.getElementById('btn_update_exercise').disabled = true;
  }
}

/**
 * Maneja el envío del formulario para actualizar un ejercicio.
 */
async function handleUpdateExercise(id) {
  const errorElement = document.getElementById('error_msg');
  errorElement.textContent = ''; // Limpiar errores previos

  // Recolecta los datos del formulario (igual que en 'nuevo')
  const payload = {
    name: document.getElementById('ex_name').value.trim(),
    main_muscle_group: document.getElementById('ex_group').value.trim(),
    description: document.getElementById('ex_desc').value.trim(),
    category: document.getElementById('ex_category').value,
    difficult_level: document.getElementById('ex_difficulty').value,
    example: document.getElementById('ex_sample').value.trim(),
    instructions: document.getElementById('ex_instructions').value.trim(),
  };

  //  Validación 
  if (!payload.name || !payload.main_muscle_group || !payload.description || !payload.category || !payload.difficult_level || !payload.example || !payload.instructions) {
    errorElement.textContent = 'Error: Debes completar todos los campos.';
    return;
  }

  try {
    // Enviar al endpoint (PUT /api/exercises/:id)
    const response = await fetchApi(`/api/exercises/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'Ocurrió un error al actualizar.');
    }

    //Éxito
    alert('¡Ejercicio actualizado exitosamente!');
    window.location.href = '/admin-exercises'; // Redirigir de vuelta al listado

  } catch (error) {
    console.error('Error al actualizar ejercicio:', error);
    errorElement.textContent = error.message;
  }
}

document.addEventListener('DOMContentLoaded', () => {
  //Obtener el ID del ejercicio de la URL
  const urlParams = new URLSearchParams(window.location.search);
  const exerciseId = urlParams.get('id');

  if (!exerciseId) {
    document.getElementById('error_msg').textContent = 'Error: No se ha especificado un ID de ejercicio. Vuelve a la lista e inténtalo de nuevo.';
    return;
  }

  // Cargar los datos del ejercicio
  loadExerciseData(exerciseId);

  // Asignar el evento al botón de guardar
  const saveButton = document.getElementById('btn_update_exercise');
  if (saveButton) {
    saveButton.addEventListener('click', () => handleUpdateExercise(exerciseId));
  }
});