/**
 * Signup page HTMX handlers and validation
 * Client-side validation and success handling
 */

// Client-side password validation before form submission
// HTMX doesn't provide built-in validation, so JS is required
document.querySelector('.login-form').addEventListener('submit', function(evt) {
    const password = document.getElementById('password').value;
    const confirmPassword = document.getElementById('confirm-password').value;
    const messageContainer = document.getElementById('message-container');
    
    // Clear previous messages
    messageContainer.innerHTML = '';
    
    // Validate password match
    if (password !== confirmPassword) {
        evt.preventDefault();
        messageContainer.innerHTML = '<div class="error-message">Passwords do not match</div>';
        return false;
    }
    
    // Validate password length
    if (password.length < 6) {
        evt.preventDefault();
        messageContainer.innerHTML = '<div class="error-message">Password must be at least 6 characters long</div>';
        return false;
    }
});

// Handle successful signup - redirect to login page
// HTMX doesn't handle redirects with sessionStorage, so JS is required
document.body.addEventListener('htmx:afterSwap', function(evt) {
    if (evt.detail.target.id === 'message-container') {
        const response = evt.detail.xhr.responseText;
        // Check if signup was successful (look for success indicator)
        if (response.includes('data-success="true"')) {
            sessionStorage.setItem('signupSuccess', 'true');
            window.location.href = '/login.html';
        }
    }
});
