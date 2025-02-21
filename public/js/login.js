$(document).ready(function () {
    // Login form validation and submission handling
    var login = {
        form: $('.login-form'),
        
        init: function() {
            console.log(this.form);
            var success = sessionStorage.getItem('signupSuccess');
            if (success === 'true') {
                sessionStorage.removeItem('signupSuccess');
                $('#signup-successful').removeClass('hidden');
            }
            this.form.submit(function(e) {
                e.preventDefault();
                login.validate();
            });
        },

        validate: function() {
            var email = $('#email').val().trim();
            var password = $('#password').val();
            var isValid = true;

            // Reset previous error states
            $('.form-group').removeClass('has-error');
            $('.error-message').remove();

            // Email validation
            if (!email || !this.isValidEmail(email)) {
                this.showError($('#email'), 'Please enter a valid email address');
                isValid = false;
            }

            // Password validation
            if (!password || password.length < 6) {
                this.showError($('#password'), 'Password must be at least 6 characters');
                isValid = false;
            }

            if (isValid) {
                this.submitForm();
            }
        },

        isValidEmail: function(email) {
            var emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            return emailRegex.test(email);
        },

        showError: function(element, message) {
            element.parent('.form-group').addClass('has-error');
            $('<div class="error-message">' + message + '</div>')
                .insertAfter(element);
        },

        submitForm: function() {
            var data = {
                email: $('#email').val().trim(),
                password: $('#password').val()
            };
            $.ajax({
                url: '/login',
                type: 'POST',
                data: JSON.stringify(data),     
                dataType: 'json',
                contentType: 'application/json',
                success: function(response) {
                    setAuthToken(response.access_token);
                    window.location.href = '/';
                },
                error: function(xhr) {
                    console.log(xhr);
                    login.showError($('#email'), 'Invalid email or password');
                }
            });
        }
    };

    // Initialize login functionality
    login.init();
});
