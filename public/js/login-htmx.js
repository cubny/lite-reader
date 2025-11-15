/**
 * Login page HTMX handlers
 * Handles token storage and redirect after successful login
 */

// Check for signup success message from sessionStorage
if (sessionStorage.getItem('signupSuccess') === 'true') {
    sessionStorage.removeItem('signupSuccess');
    document.getElementById('signup-successful').classList.remove('hidden');
}

// Store token and redirect on successful login
// HTMX cannot access response data to extract token, so JS is required
document.body.addEventListener('htmx:afterSwap', function(evt) {
    if (evt.detail.target.id === 'message-container') {
        const response = evt.detail.xhr.responseText;
        // Check if we got a token (successful login)
        if (response.includes('data-token')) {
            const parser = new DOMParser();
            const doc = parser.parseFromString(response, 'text/html');
            const tokenElement = doc.querySelector('[data-token]');
            if (tokenElement) {
                const token = tokenElement.getAttribute('data-token');
                localStorage.setItem('authToken', token);
                window.location.href = '/';
            }
        }
    }
});
