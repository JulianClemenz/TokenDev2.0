
function getToken() {
    return sessionStorage.getItem('access_token');
}

function getCurrentUser() {
    const userStr = sessionStorage.getItem('user');
    if (!userStr) {
        logout();
        return null;
    }
    return JSON.parse(userStr);
}

async function fetchApi(url, options = {}) {
    const token = getToken();
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
        ...options.headers,
    };
    const response = await fetch(url, { ...options, headers });
    if (response.status === 401 && !url.endsWith('/password')) {
        // No cerramos sesión si el error 401 es por "contraseña incorrecta"
        alert('Tu sesión ha expirado. Por favor, inicia sesión de nuevo.');
        logout();
        throw new Error('No autorizado');
    }
    return response;
}

function logout() {
    sessionStorage.removeItem('access_token');
    sessionStorage.removeItem('refresh_token');
    sessionStorage.removeItem('user');
    window.location.href = '/';
}


// --- Lógica de la Página de Cambio de Contraseña ---

/**
 * Envía la solicitud de cambio de contraseña.
 */
async function handleChangePassword(userId) {
    const errorElement = document.getElementById('error_msg');
    const successElement = document.getElementById('success_msg');
    errorElement.textContent = '';
    successElement.textContent = '';

    const currentInput = document.getElementById('pass_current');
    const newInput = document.getElementById('pass_new');
    const confirmInput = document.getElementById('pass_confirm');

    try {
        // 1. Recolectar datos (Payload debe coincidir con PasswordChange DTO)
        const payload = {
            current_password: currentInput.value,
            new_password: newInput.value,
            confirm_password: confirmInput.value
        };

        // 2. Validación simple de cliente
        if (!payload.current_password || !payload.new_password || !payload.confirm_password) {
            throw new Error('Todos los campos son obligatorios.');
        }
        if (payload.new_password !== payload.confirm_password) {
            throw new Error('La nueva contraseña y su confirmación no coinciden.');
        }
        if (payload.new_password.length < 7) {
            throw new Error('La nueva contraseña debe tener al menos 7 caracteres.');
        }

        // 3. Enviar a la API (POST /api/users/:id/password)
        const response = await fetchApi(`/api/users/${userId}/password`, {
            method: 'POST',
            body: JSON.stringify(payload)
        });

        const data = await response.json();

        if (!response.ok) {
            // El backend responde con 401, 400, 409, 500
            throw new Error(data.error || 'Ocurrió un error al cambiar la contraseña.');
        }

        // 4. Éxito
        successElement.textContent = data.mensaje;
        currentInput.value = '';
        newInput.value = '';
        confirmInput.value = '';

    } catch (error) {
        errorElement.textContent = error.message;
        console.error('Error al cambiar contraseña:', error);
    }
}


document.addEventListener('DOMContentLoaded', () => {
    const currentUser = getCurrentUser();
    if (!currentUser || !currentUser.id) {
        document.getElementById('error_msg').textContent = 'Error de autenticación.';
        return;
    }

    const userId = currentUser.id;

    // Asignar evento al botón de guardar
    const saveButton = document.getElementById('btn_save_password');
    saveButton.addEventListener('click', () => handleChangePassword(userId));
});
