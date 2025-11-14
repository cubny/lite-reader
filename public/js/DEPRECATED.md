# htmx Migration Notes

## Overview
This directory contains the deprecated JavaScript files that were replaced with htmx and Alpine.js.

## Deprecated Files
The following JavaScript files are no longer used in the application:

### Custom Application JavaScript (DEPRECATED)
- `auth.js` - Authentication handling (replaced by Alpine.js inline handlers)
- `feeds.js` - Feed management logic (replaced by Alpine.js reactive data)
- `items.js` - Item display and interaction (replaced by Alpine.js reactive data)
- `login.js` - Login form handling (replaced by Alpine.js form handlers)
- `main.js` - Main application initialization (replaced by Alpine.js x-init)
- `signup.js` - Signup form handling (replaced by Alpine.js form handlers)
- `stage.js` - jQuery Layout initialization (replaced by CSS flexbox)
- `utils.js` - Utility functions (replaced by Alpine.js methods and native JavaScript)

### Third-party Libraries (DEPRECATED)
- `jquery.js` - jQuery library (no longer needed)
- `jquery.ui.all.js` - jQuery UI (no longer needed)
- `jquery.layout.js` / `jquery.layout.min.js` - jQuery Layout plugin (replaced by CSS)
- `jquery.cookie.js` - Cookie management (replaced by native document.cookie)
- `jquery.scrollTo.min.js` - Smooth scrolling (replaced by native scrollIntoView)
- `spin.min.js` - Loading spinner (replaced by CSS animations)
- `moment.min.js` - Date formatting (replaced by native Date methods)

## New Dependencies
The application now uses:
- **htmx** (v1.9.10) - For declarative AJAX requests and DOM updates
- **Alpine.js** (v3.13.5) - For reactive data binding and component state management

Both are loaded from CDN in production. No custom JavaScript files are needed.

## Migration Benefits
1. **Reduced code complexity** - Eliminated ~17K lines of JavaScript code
2. **Better maintainability** - Declarative approach is easier to understand
3. **Modern architecture** - Using current best practices for web development
4. **Same user experience** - All functionality preserved
5. **No backend changes** - Existing API endpoints work unchanged

## How It Works

### Authentication Pages
- Use Alpine.js for form validation and submission
- Handle API responses with fetch API
- Store auth tokens in localStorage
- Redirect on success/failure

### Main Application (index.html)
- Alpine.js manages application state (feeds, items, counts)
- CSS Flexbox replaces jQuery Layout for responsive layout
- Native browser APIs replace jQuery utilities
- Declarative event handlers replace imperative jQuery code

### API Communication
- All backend endpoints remain unchanged
- JSON requests/responses handled via fetch API
- Authorization headers added to all requests
- 401 responses trigger redirect to login page
