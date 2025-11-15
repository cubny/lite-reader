// Signup page JavaScript - Listen for successful signup
// Form submission and validation handled by HTML5 + HTMX
document.addEventListener('DOMContentLoaded', function() {
    // Listen for successful signup (redirect via HX-Redirect header)
    document.body.addEventListener('htmx:beforeSwap', function(event) {
        if (event.detail.target.id === 'form-messages' && event.detail.xhr.status === 201) {
            // Signup successful, set flag for login page
            sessionStorage.setItem('signupSuccess', 'true');
        }
    });
});