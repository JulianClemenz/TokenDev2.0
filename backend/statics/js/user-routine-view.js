// --- Funciones de Ayuda ---

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
 * Carga los datos de UNA rutina y TODOS sus ejercicios.
 */
async function loadRoutineView() {
    // Obtener el ID de la rutina de la URL
    const urlParams = new URLSearchParams(window.location.search);
    const routineId = urlParams.get('id');

    // Elementos del DOM donde mostraremos los datos
    const routineNameEl = document.getElementById('routine-name-display');
    const exerciseListEl = document.getElementById('exercise-list-container');
    const errorEl = document.getElementById('error_msg');

    if (!routineId) {
        errorEl.textContent = 'Error: No se especificó una rutina para ver. Vuelve a la lista.';
        return;
    }

    try {
        // Mostrar "Cargando..."
        routineNameEl.textContent = 'Cargando rutina...';
        exerciseListEl.innerHTML = '';

        // Llama a GetRoutineById en RoutineHandler
        const routineResponse = await fetchApi(`/api/routines/${routineId}`);
        if (!routineResponse.ok) {
            const err = await routineResponse.json();
            throw new Error(err.error || 'No se pudo cargar la rutina.');
        }
        const routine = await routineResponse.json();

        // Mostrar el nombre de la rutina
        routineNameEl.textContent = routine.Name;

        // Revisar si hay ejercicios en la lista
        // (Usamos la propiedad 'ExcerciseList' tal como está en tu DTO de Go)
        if (!routine.ExcerciseList || routine.ExcerciseList.length === 0) {
            exerciseListEl.innerHTML = '<p class="text-muted fs-5">Esta rutina aún no tiene ejercicios.</p>';
            return;
        }

        // Obtener los detalles (nombre) de CADA ejercicio en la lista        
        const exerciseDetailPromises = routine.ExcerciseList.map(async (exerciseItem) => {
            try {
                // Llama a GetExcerciseByID en ExerciseHandler
                const exerciseResponse = await fetchApi(`/api/exercises/${exerciseItem.exercise_id}`);
                if (!exerciseResponse.ok) {
                    console.warn(`No se pudo cargar el ejercicio ID: ${exerciseItem.exercise_id}`);
                    return { ...exerciseItem, Name: 'Ejercicio no encontrado' }; // Fallback
                }
                const exerciseDetails = await exerciseResponse.json();

                // Combinamos los datos (Series, Reps, Peso) con el nombre del ejercicio
                return {
                    ...exerciseItem, // Contiene Repetitions, Series, Weight
                    Name: exerciseDetails.Name // Añade el nombre del ejercicio
                };
            } catch (e) {
                console.error(e);
                return { ...exerciseItem, Name: 'Error al cargar ejercicio' }; // Fallback
            }
        });

        // Esperamos a que todas las llamadas a GetExcerciseByID terminen
        const exercisesWithDetails = await Promise.all(exerciseDetailPromises);

        // Renderizar los ejercicios en tarjetas
        exercisesWithDetails.forEach(ex => {
            const card = document.createElement('div');
            card.className = 'col-md-6 col-lg-4 mb-3';
            card.innerHTML = `
                <div class="card h-100 shadow-sm">
                    <div class="card-body">
                        <h5 class="card-title text-primary">${ex.Name}</h5>
                        <ul class="list-group list-group-flush">
                            <li class="list-group-item"><strong>Series:</strong> ${ex.series}</li>
                            <li class="list-group-item"><strong>Repeticiones:</strong> ${ex.repetitions}</li>
                            <li class="list-group-item"><strong>Peso:</strong> ${ex.weight} kg</li>
                        </ul>
                    </div>
                </div>
            `;
            exerciseListEl.appendChild(card);
        });

    } catch (error) {
        console.error('Error al cargar la vista de rutina:', error);
        errorEl.textContent = error.message;
        routineNameEl.textContent = 'Error al cargar';
    }
}

// --- Inicialización ---
document.addEventListener('DOMContentLoaded', loadRoutineView);

