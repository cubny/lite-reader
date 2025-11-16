// Login page - HTMX handles form submission
// This script only handles: signup success message and auth token storage
document.addEventListener('DOMContentLoaded', function() {
    // Show signup success message if redirected from signup
    var success = sessionStorage.getItem('signupSuccess');
    if (success === 'true') {
        sessionStorage.removeItem('signupSuccess');
        document.getElementById('signup-successful').classList.remove('hidden');
    }

    // Listen for HTMX events to handle auth token
    document.body.addEventListener('htmx:beforeSwap', function(event) {
        // Check if this is the login form response
        if (event.detail.target && event.detail.target.id === 'form-messages') {
            var xhr = event.detail.xhr;
            if (xhr && xhr.status === 200) {
                // Check for auth token in cookie
                var cookies = document.cookie.split(';');
                for (var i = 0; i < cookies.length; i++) {
                    var cookie = cookies[i].trim();
                    console.log('Checking cookie:', cookie);
                    if (cookie.startsWith('authToken=')) {
                        var token = cookie.substring('authToken='.length);
                        if (token) {
                            // Store token in localStorage
                            setAuthToken(token);
                            // Clear the cookie
                            document.cookie = 'authToken=; path=/; max-age=0';
                        }
                    }
                }
            }
        }
    });
});
