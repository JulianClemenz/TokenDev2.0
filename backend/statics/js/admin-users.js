
function getToken() {
  return sessionStorage.getItem('access_token');
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
    window.location.href = '/login'; 
    throw new Error('No autorizado');
  }

  return response;
}

async function loadUsers() {
  const tableBody = document.querySelector('.table tbody');
  tableBody.innerHTML = '<tr><td colspan="12">Cargando usuarios...</td></tr>';

  try {
    const response = await fetchApi('/api/admin/stats/users');

    if (!response.ok) {
      if (response.status === 204) {
        tableBody.innerHTML = '<tr><td colspan="12">No hay usuarios registrados.</td></tr>';
        return;
      }
      throw new Error(`Error ${response.status}: No se pudieron cargar los usuarios.`);
    }

    const data = await response.json();
    tableBody.innerHTML = '';

    if (data.users && data.users.length > 0) {
      data.users.forEach((user, index) => {
        const row = document.createElement('tr');

        const birthDate = user.BirthDate ? user.BirthDate.split('T')[0] : 'N/D';

        const roleBadge = user.Role === 'admin' 
          ? '<span class="badge bg-success">Admin</span>' 
          : '<span class="badge bg-secondary">Client</span>';


        const userDataString = encodeURIComponent(JSON.stringify(user));

        row.innerHTML = `
          <th scope="row">${index + 1}</th>
          <td>${user.UserName || ''}</td>
          <td>${user.Name || ''}</td>
          <td>${user.LastName || ''}</td>
          <td>${user.Email || ''}</td>
          <td>${birthDate}</td>
          <td>${user.Height || 0} cm</td>
          <td>${user.Weight || 0} kg</td>
          <td>${user.Experience || ''}</td>
          <td>${user.Objetive || ''}</td>
          <td>${roleBadge}</td>
          <td>
            <button 
              type="button" 
              class="btn btn-outline-warning btn-sm btn-promote" 
              data-user-id="${user.id}"
              data-user-data="${userDataString}"
              ${user.Role === 'admin' ? 'disabled' : ''}>
              Hacer Admin
            </button>
          </td>
        `;
        tableBody.appendChild(row);
      });
    } else {
      tableBody.innerHTML = '<tr><td colspan="12">No hay usuarios registrados.</td></tr>';
    }
  } catch (error) {
    console.error('Error al cargar usuarios:', error);
    tableBody.innerHTML = `<tr class="text-center"><td colspan="12" class="text-danger">Error: ${error.message}</td></tr>`;
  }
}


async function handlePromoteUser(event) {
  const button = event.target;
  const userId = button.dataset.userId;
  const errorElement = document.getElementById('error_msg');
  errorElement.textContent = '';

  if (!confirm('¿Estás seguro de que deseas ascender a este usuario a Administrador?')) {
    return;
  }

  try {
    const userData = JSON.parse(decodeURIComponent(button.dataset.userData));


    const payload = {
      user_name: userData.UserName,
      email: userData.Email,
      role: 'admin',
      weight: userData.Weight,
      height: userData.Height,
      experience: userData.Experience,
      objetive: userData.Objetive
    };

    const response = await fetchApi(`/api/users/${userId}`, {
      method: 'PUT',
      body: JSON.stringify(payload)
    });

    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || 'No se pudo actualizar el rol del usuario.');
    }

    alert('¡Usuario ascendido a Administrador exitosamente!');
    loadUsers(); 

  } catch (error) {
    console.error('Error al ascender a admin:', error);
    errorElement.textContent = error.message;
  }
}


document.addEventListener('DOMContentLoaded', () => {
  loadUsers();


  const tableBody = document.querySelector('.table tbody');
  tableBody.addEventListener('click', (event) => {
    if (event.target.classList.contains('btn-promote')) {
      handlePromoteUser(event);
    }
  });
});