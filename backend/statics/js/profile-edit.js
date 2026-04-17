// --- Funciones de Ayuda ---

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
    if (response.status === 401) {
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
// --- Lógica de la Página de Edición ---

/**
 * Carga los datos actuales del usuario en los inputs del formulario.
 */
async function loadCurrentData(userId) {
    const errorElement = document.getElementById('error_msg');
    try {
        const response = await fetchApi(`/api/users/${userId}`);
        if (!response.ok) {
            throw new Error('No se pudieron cargar tus datos.');
        }

        const user = await response.json(); // UserResponseDTO

        // Rellenar el formulario
        document.getElementById('edit_name').value = user.Name;
        document.getElementById('edit_lastname').value = user.LastName;
        document.getElementById('edit_email').value = user.Email;

        // (Campos editables)
        document.getElementById('edit_username').value = user.UserName;
        document.getElementById('edit_height').value = user.Height;
        document.getElementById('edit_weight').value = user.Weight;
        document.getElementById('edit_experience').value = user.Experience;
        document.getElementById('edit_objective').value = user.Objetive;

    } catch (error) {
        errorElement.textContent = error.message;
        console.error('Error al cargar datos:', error);
    }
}

/**
 * Envía los cambios del perfil a la API.
 */
async function handleSaveChanges(userId, userRole) {
    const errorElement = document.getElementById('error_msg');
    errorElement.textContent = '';

    try {
        // 1. Recolectar datos del formulario
        const payload = {
            user_name: document.getElementById('edit_username').value.trim(),
            email: document.getElementById('edit_email').value.trim(),
            height: parseFloat(document.getElementById('edit_height').value),
            weight: parseFloat(document.getElementById('edit_weight').value),
            experience: document.getElementById('edit_experience').value,
            objetive: document.getElementById('edit_objective').value,
            role: userRole
        };

        // 2. Validación simple
        if (!payload.user_name || !payload.email || payload.height <= 0 || payload.weight <= 0) {
            throw new Error('Completa todos los campos. El peso y la altura deben ser positivos.');
        }

        // 3. Enviar a la API (PUT /api/users/:id)
        const response = await fetchApi(`/api/users/${userId}`, {
            method: 'PUT',
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            const err = await response.json();
            throw new Error(err.error || 'No se pudo actualizar el perfil.');
        }

        const updatedUser = await response.json(); // UserModifyResponseDTO

        const currentUser = getCurrentUser();
        currentUser.UserName = updatedUser.UserName;
        currentUser.Email = updatedUser.Email;
        sessionStorage.setItem('user', JSON.stringify(currentUser));

        // 5. Éxito
        alert('Perfil actualizado exitosamente.');
        window.location.href = 'profile.html';

    } catch (error) {
        errorElement.textContent = error.message;
        console.error('Error al guardar cambios:', error);
    }
}


document.addEventListener('DOMContentLoaded', () => {
    const currentUser = getCurrentUser();
    if (!currentUser || !currentUser.id) {
        document.getElementById('error_msg').textContent = 'Error de autenticación.';
        return;
    }

    const userId = currentUser.id;
    const userRole = currentUser.role; // Obtenemos el rol de la sesión

    // 1. Cargar los datos en el formulario
    loadCurrentData(userId);

    // 2. Asignar evento al botón de guardar
    const saveButton = document.getElementById('btn_save_changes');
    saveButton.addEventListener('click', () => handleSaveChanges(userId, userRole));
});
