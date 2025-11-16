// Token management functions
function setAuthToken(token) {
    localStorage.setItem('authToken', token);
}

function getAuthToken() {
    let token = localStorage.getItem('authToken');
    if (!token) {
        console.log('Auth token not found in localStorage, checking cookies');
        // read from cookie as fallback
        const cookies = document.cookie.split(';');
        for (let i = 0; i < cookies.length; i++) {
            let cookie = cookies[i].trim();
            if (cookie.startsWith('authToken=')) {
                token = cookie.substring('authToken='.length);
                setAuthToken(token);
                document.cookie = 'authToken=; path=/; max-age=0';
                return token;
            }
        }
    }
    return token;
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