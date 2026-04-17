
/**
 * Inicia sesión de un usuario.
 * Se conecta al endpoint POST /login
 * @param {string} email
 * @param {string} password
 */
async function login(email, password) {
    // El backend espera un JSON con "email" y "password"
    const response = await fetch('/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email: email, password: password }),
    });

    const data = await response.json();

    if (!response.ok) {
        // Si la API devuelve un error (401, 500, etc.), lo lanzamos
        throw new Error(data.error || 'Error al iniciar sesión');
    }

    console.log('Login exitoso:', data);

    // Guardamos los tokens y los datos del usuario en sessionStorage
    sessionStorage.setItem('access_token', data.access_token);
    sessionStorage.setItem('refresh_token', data.refresh_token);
    sessionStorage.setItem('user', JSON.stringify(data.user));

    return data;
}

/**
 * Registra un nuevo usuario.
 * Se conecta al endpoint POST /register
 * @param {object} payload 
 */
async function register(payload) {

    // 1. "Traducimos" los nombres del JS a los que espera el UserRegisterDTO de Go
    const apiPayload = {
        name: payload.name,
        last_name: payload.surname,
        user_name: payload.username,
        email: payload.email,
        password: payload.password,
        birth_date: payload.birthdate,
        weight: payload.weight_kg,
        height: payload.height_cm,
        experience: payload.experience,
        objetive: payload.goal,
        role: 'client'
    };

    // 2. Llamamos al endpoint del backend
    const response = await fetch('/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(apiPayload),
    });

    const data = await response.json();

    if (!response.ok) {
        // Si la API da error (ej: "email ya existe"), lo mostramos
        //
        throw new Error(data.error || 'Error al registrar usuario');
    }

    // 3. Éxito
    console.log('Usuario registrado:', data);
    return data;
}