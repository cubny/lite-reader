// Token management functions
function setAuthToken(token) {
    localStorage.setItem('authToken', token);
}

function getAuthToken() {
    return localStorage.getItem('authToken');
}

function clearAuthToken() {
    localStorage.removeItem('authToken');
}

// Auth-related AJAX setup
$.ajaxSetup({
    beforeSend: function(xhr) {
        const token = getAuthToken();
        if (token) {
            xhr.setRequestHeader('Authorization', 'Bearer ' + token);
        }
    },
    error: function(jqXHR) {
        if (jqXHR.status === 401) {
            clearAuthToken();
            window.location.href = '/login.html';
            return false;
        }
    },
    success: function(response) {
        if (response.redirect) {
            window.location.href = response.redirect;
            return false;
        }
    }
});


function logout() {
    clearAuthToken();
    window.location.href = '/login.html';
} 