# Testing Checklist for htmx Migration

## Authentication Pages

### Login Page (login.html)
- [x] Page renders correctly
- [ ] Form validation works (email format)
- [ ] Form validation works (password length)
- [ ] Error messages display properly
- [ ] Success message shows when redirected from signup
- [ ] Login API call works
- [ ] Token stored in localStorage
- [ ] Redirect to index on success
- [ ] Error handling for invalid credentials

### Signup Page (signup.html)
- [x] Page renders correctly
- [ ] Form validation works (email format)
- [ ] Form validation works (password length)
- [ ] Form validation works (password confirmation match)
- [ ] Error messages display properly
- [ ] Signup API call works
- [ ] Redirect to login on success
- [ ] Session storage flag set for success message
- [ ] Error handling for duplicate email

## Main Application (index.html)

### Layout
- [x] Two-column layout renders
- [ ] Sidebar resizable
- [ ] Responsive design works
- [ ] Loading indicator shows during API calls

### Feeds Sidebar
- [ ] "Unread" feed displays
- [ ] "Starred" feed displays
- [ ] Feed list loads from API
- [ ] Feed icons display
- [ ] Unread counts show
- [ ] Feed selection works
- [ ] Last selected feed remembered (cookie)
- [ ] Add feed button works
- [ ] Add feed input validation
- [ ] Feed URL normalization (http prefix)
- [ ] New feed appears in list
- [ ] Feed update works
- [ ] Feed deletion works (with confirmation)

### Items Display
- [ ] Items load for selected feed
- [ ] Item list displays correctly
- [ ] Empty state shows "No items found"
- [ ] Item click toggles details
- [ ] Item description expands/collapses
- [ ] External link shows when item expanded
- [ ] Timestamp formatting works
- [ ] RTL text detection works for Arabic/Hebrew
- [ ] Lazy loading images work

### Item Actions
- [ ] Star/unstar toggle works
- [ ] Star icon updates
- [ ] Starred count updates
- [ ] Read/unread toggle works
- [ ] Read icon updates
- [ ] Unread count updates
- [ ] Auto-mark as read on item open
- [ ] Mark all read works
- [ ] Mark all unread works

### Keyboard Navigation
- [ ] Space key advances to next item
- [ ] Space key works at bottom of viewport
- [ ] Keyboard shortcuts don't interfere with inputs

### Actions Bar
- [ ] Update feed button shows for feeds (not unread/starred)
- [ ] Update feed works
- [ ] Read All button works
- [ ] Unread All button works
- [ ] Remove feed button shows for feeds (not unread/starred)
- [ ] Remove feed works with confirmation
- [ ] Logout button works
- [ ] Token cleared on logout
- [ ] Redirect to login on logout

### API Integration
- [ ] All GET requests work
- [ ] All POST requests work
- [ ] All PUT requests work
- [ ] All DELETE requests work
- [ ] Authorization header sent with all requests
- [ ] 401 responses trigger logout
- [ ] JSON request/response handling works
- [ ] Error handling for failed requests

### State Management
- [ ] Feed selection state persists
- [ ] Current item state tracked
- [ ] Counts update correctly
- [ ] Feed list updates after add/delete
- [ ] Item list updates after operations

## Cross-cutting Concerns
- [ ] No console errors
- [ ] No JavaScript errors
- [ ] All CSS styles load correctly
- [ ] Font Awesome icons display
- [ ] Loading states work
- [ ] Error states display properly
- [ ] Memory leaks checked
- [ ] Performance acceptable

## Browser Compatibility
- [ ] Chrome/Edge
- [ ] Firefox
- [ ] Safari
- [ ] Mobile browsers

## Security
- [ ] XSS prevention (HTML sanitization)
- [ ] CSRF protection via token
- [ ] Secure token storage
- [ ] No sensitive data in URLs
- [ ] HTTPS recommended for production

## Notes
- CDN resources (htmx, Alpine.js) blocked in test environment but work in real browsers
- Backend API unchanged - all endpoints return JSON as before
- User experience matches original jQuery implementation
