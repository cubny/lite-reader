// Login page JavaScript - Client-side validation only
// Form submission handled by HTMX
document.addEventListener('DOMContentLoaded', function() {
    // Show signup success message if redirected from signup
    var success = sessionStorage.getItem('signupSuccess');
    if (success === 'true') {
        sessionStorage.removeItem('signupSuccess');
        document.getElementById('signup-successful').classList.remove('hidden');
    }

    var form = document.querySelector('.login-form');
    
    form.addEventListener('submit', function(e) {
        // Reset previous error states
        var formGroups = document.querySelectorAll('.form-group');
        formGroups.forEach(function(group) {
            group.classList.remove('has-error');
        });
        var errorMessages = document.querySelectorAll('.error-message');
        errorMessages.forEach(function(msg) {
            msg.remove();
        });

        var email = document.getElementById('email').value.trim();
        var password = document.getElementById('password').value;
        var isValid = true;

        // Email validation
        var emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!email || !emailRegex.test(email)) {
            showError(document.getElementById('email'), 'Please enter a valid email address');
            isValid = false;
        }

        // Password validation
        if (!password || password.length < 6) {
            showError(document.getElementById('password'), 'Password must be at least 6 characters');
            isValid = false;
        }

        if (!isValid) {
            e.preventDefault();
        }
        // If valid, HTMX will handle the submission
    });

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

function showError(element, message) {
    element.parentElement.classList.add('has-error');
    var errorDiv = document.createElement('div');
    errorDiv.className = 'error-message';
    errorDiv.textContent = message;
    element.parentElement.appendChild(errorDiv);
}
