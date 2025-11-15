// Signup page JavaScript - Client-side validation only
// Form submission handled by HTMX
document.addEventListener('DOMContentLoaded', function() {
    var form = document.querySelector('.login-form');
    
    form.addEventListener('submit', function(e) {
        // Reset previous errors
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
        var confirmPassword = document.getElementById('confirm-password').value;
        var isValid = true;

        // Validate email
        var emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(email)) {
            showError(document.getElementById('email'), 'Please enter a valid email address');
            isValid = false;
        }

        // Validate password
        if (!password || password.length < 6) {
            showError(document.getElementById('password'), 'Password must be at least 6 characters long');
            isValid = false;
        }

        // Validate password confirmation
        if (password !== confirmPassword) {
            showError(document.getElementById('confirm-password'), 'Passwords do not match');
            isValid = false;
        }

        if (!isValid) {
            e.preventDefault();
        }
        // If valid, HTMX will handle the submission
    });

    // Listen for successful signup (redirect via HX-Redirect header)
    document.body.addEventListener('htmx:beforeSwap', function(event) {
        if (event.detail.target.id === 'form-messages' && event.detail.xhr.status === 201) {
            // Signup successful, set flag for login page
            sessionStorage.setItem('signupSuccess', 'true');
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