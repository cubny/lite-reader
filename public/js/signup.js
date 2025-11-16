// Signup page - HTMX handles form submission
// This script only handles: setting success flag for login page
document.addEventListener('DOMContentLoaded', function() {
    // Client-side password confirmation validation
    var form = document.querySelector('.login-form');
    var password = document.getElementById('password');
    var confirmPassword = document.getElementById('confirm-password');
    
    // Add custom validation for password confirmation
    confirmPassword.addEventListener('input', function() {
        if (password.value !== confirmPassword.value) {
            confirmPassword.setCustomValidity('Passwords do not match');
        } else {
            confirmPassword.setCustomValidity('');
        }
    });
    
    password.addEventListener('input', function() {
        if (password.value !== confirmPassword.value) {
            confirmPassword.setCustomValidity('Passwords do not match');
        } else {
            confirmPassword.setCustomValidity('');
        }
    });

    // Listen for successful signup
    document.body.addEventListener('htmx:beforeSwap', function(event) {
        if (event.detail.target && event.detail.target.id === 'form-messages') {
            var xhr = event.detail.xhr;
            if (xhr && xhr.status === 201) {
                // Signup successful, set flag for login page
                sessionStorage.setItem('signupSuccess', 'true');
            }
        }
    });
});
