
/**
 * Obtiene el token de autenticación desde sessionStorage.
 */
function getToken() {
    return sessionStorage.getItem('access_token');
}

function getCurrentUser() {
    const userStr = sessionStorage.getItem('user');
    if (!userStr) {
        logout(); // Redirigir al login si no hay usuario
        return null;
    }
    return JSON.parse(userStr);
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
        logout();
        throw new Error('No autorizado');
    }

    return response;
}

/**
 * Función de logout (usada por helpers)
 */
function logout() {
    sessionStorage.removeItem('access_token');
    sessionStorage.removeItem('refresh_token');
    sessionStorage.removeItem('user');
    window.location.href = '/';
}

// --- Lógica de la Página de Perfil ---

/**
 * Carga los datos del perfil del usuario.
 */
async function loadProfile() {
    const errorElement = document.getElementById('error_msg');
    const currentUser = getCurrentUser();

    if (!currentUser || !currentUser.id) {
        errorElement.textContent = 'No se pudo identificar al usuario. Inicia sesión de nuevo.';
        return;
    }

    const userId = currentUser.id;

    try {
        //Llamar a la API (GET /api/users/:id)
        const response = await fetchApi(`/api/users/${userId}`);
        if (!response.ok) {
            const err = await response.json();
            throw new Error(err.error || 'No se pudieron cargar los datos del perfil.');
        }

        const user = await response.json();

        document.getElementById('profile_username').textContent = user.UserName;
        document.getElementById('profile_name').textContent = user.Name;
        document.getElementById('profile_lastname').textContent = user.LastName;
        document.getElementById('profile_email').textContent = user.Email;

        // Formatear la fecha
        const birthDate = new Date(user.BirthDate);
        document.getElementById('profile_birthdate').textContent = birthDate.toLocaleDateString('es-ES', {
            day: '2-digit', month: '2-digit', year: 'numeric'
        });

        document.getElementById('profile_height').textContent = `${user.Height} cm`;
        document.getElementById('profile_weight').textContent = `${user.Weight} kg`;
        document.getElementById('profile_experience').textContent = user.Experience;
        document.getElementById('profile_objective').textContent = user.Objetive;

    } catch (error) {
        console.error('Error al cargar perfil:', error);
        errorElement.textContent = error.message;
    }
}

// --- Inicialización ---
document.addEventListener('DOMContentLoaded', loadProfile);