// Login page JavaScript - Shows signup success message
// Form submission and validation handled by HTML5 + HTMX
document.addEventListener('DOMContentLoaded', function() {
    // Show signup success message if redirected from signup
    var success = sessionStorage.getItem('signupSuccess');
    if (success === 'true') {
        sessionStorage.removeItem('signupSuccess');
        document.getElementById('signup-successful').classList.remove('hidden');
    }

    // HTMX event listener for successful login
    document.body.addEventListener('htmx:beforeSwap', function(event) {
        // Check if this is the login form
        if (event.detail.target.id === 'form-messages') {
            // Check for auth token cookie set by server
            var cookies = document.cookie.split(';');
            for (var i = 0; i < cookies.length; i++) {
                var cookie = cookies[i].trim();
                if (cookie.startsWith('authToken=')) {
                    var token = cookie.substring('authToken='.length);
                    if (token) {
                        // Store token in localStorage
                        setAuthToken(token);
                        // Clear the cookie (we only used it for transfer)
                        document.cookie = 'authToken=; path=/; max-age=0';
                    }
                }
            }
        }
    });
});
