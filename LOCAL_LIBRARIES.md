# Local JavaScript Libraries

## Overview
This document describes the local JavaScript libraries used in Lite Reader after the htmx migration.

## Library Files

### Location
All JavaScript libraries are stored in `public/js/vendor/`:

```
public/js/vendor/
├── alpine.min.js    (44KB) - Alpine.js v3.13.5
├── htmx.min.js      (51KB) - htmx v2.0.8
└── json-enc.js      (619B) - htmx JSON encoding extension
```

## Why Local Files?

1. **No External Dependencies** - Application works without internet access
2. **Testing** - Enables full functional testing in sandbox environments
3. **Reliability** - No CDN downtime or network issues
4. **Performance** - Potentially faster load times (no DNS lookup, no SSL handshake to external domain)
5. **Security** - No third-party CDN access, full control over code
6. **Privacy** - No external requests that could track users

## Installation

Libraries are installed via npm during development:

```bash
npm install htmx.org alpinejs
cp node_modules/htmx.org/dist/htmx.min.js public/js/vendor/
cp node_modules/alpinejs/dist/cdn.min.js public/js/vendor/alpine.min.js
cp node_modules/htmx.org/dist/ext/json-enc.js public/js/vendor/
```

The vendor files are committed to the repository, so end users don't need npm.

## Usage

### In HTML Files

**Login and Signup Pages:**
```html
<script src="js/vendor/htmx.min.js"></script>
<script defer src="js/vendor/alpine.min.js"></script>
```

**Main Application (index.html):**
```html
<script src="js/vendor/htmx.min.js"></script>
<script src="js/vendor/json-enc.js"></script>
<script defer src="js/vendor/alpine.min.js"></script>
```

## Library Details

### htmx v2.0.8
- **Purpose:** Declarative AJAX, DOM updates
- **Size:** 51KB minified
- **License:** BSD 2-Clause
- **Website:** https://htmx.org
- **Used for:** Minimal HTTP interaction patterns (though Alpine.js does most of the work)

### Alpine.js v3.13.5
- **Purpose:** Reactive UI framework
- **Size:** 44KB minified
- **License:** MIT
- **Website:** https://alpinejs.dev
- **Used for:** 
  - Form validation
  - Reactive state management (feeds, items, counts)
  - DOM binding and updates
  - Event handling

### json-enc.js
- **Purpose:** htmx extension for JSON encoding
- **Size:** 619 bytes
- **License:** BSD 2-Clause (part of htmx)
- **Used for:** JSON request encoding (minimal usage)

## Version Updates

To update libraries:

1. Update npm packages:
   ```bash
   npm update htmx.org alpinejs
   ```

2. Copy updated files:
   ```bash
   cp node_modules/htmx.org/dist/htmx.min.js public/js/vendor/
   cp node_modules/alpinejs/dist/cdn.min.js public/js/vendor/alpine.min.js
   cp node_modules/htmx.org/dist/ext/json-enc.js public/js/vendor/
   ```

3. Test thoroughly before committing

## Browser Compatibility

Both libraries support modern browsers:
- Chrome/Edge 90+
- Firefox 88+
- Safari 14+
- Mobile browsers (iOS Safari 14+, Chrome Android 90+)

## Notes

- The htmx json-enc extension shows a minor version warning (htmx 1 extension with htmx 2) but is non-critical
- Alpine.js must be loaded with `defer` attribute to ensure DOM is ready
- htmx should be loaded before Alpine.js
- Total JavaScript bundle size: ~96KB minified (~30KB gzipped)

## Comparison with Previous Setup

### Before (jQuery-based):
- jQuery: 120KB
- jQuery UI: 304KB
- jQuery Layout: 83KB
- Custom JS: 17KB
- **Total: ~524KB minified (~150KB+ gzipped)**

### After (htmx + Alpine.js):
- htmx: 51KB
- Alpine.js: 44KB
- json-enc: 0.6KB
- **Total: ~96KB minified (~30KB gzipped)**

**Reduction: 82% smaller bundle size**
