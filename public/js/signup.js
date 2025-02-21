$(document).ready(function () {
    // Signup form validation and submission handling
    var signup = {
        form: $('.login-form'),
        
        init: function() {
            this.form.submit(function(e) {
                e.preventDefault();
                signup.validate();
            });
        },

        validate: function() {
            // Reset previous errors
            $('.form-group').removeClass('has-error');
            $('.error-message').remove();

            var email = $('#email').val().trim();
            var password = $('#password').val();
            var confirmPassword = $('#confirm-password').val();
            var isValid = true;

            // Validate email
            if (!this.isValidEmail(email)) {
                this.showError($('#email'), 'Please enter a valid email address');
                isValid = false;
            }

            // Validate password
            if (!this.isValidPassword(password)) {
                this.showError($('#password'), 'Password must be at least 6 characters long');
                isValid = false;
            }

            // Validate password confirmation
            if (password !== confirmPassword) {
                this.showError($('#confirm-password'), 'Passwords do not match');
                isValid = false;
            }

            if (isValid) {
                this.submitForm(email, password);
            }

            return isValid;
        },

        isValidEmail: function(email) {
            var emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            return emailRegex.test(email);
        },

        isValidPassword: function(password) {
            return password && password.length >= 6;
        },
        
        showError: function(element, message) {
            element.parent('.form-group').addClass('has-error');
            $('<div class="error-message">' + message + '</div>')
                .insertAfter(element);
        },

        submitForm: function(email, password) {
            $.ajax({
                url: '/signup',
                type: 'POST',
                data: JSON.stringify({
                    email: email,
                    password: password
                }),
                contentType: 'application/json',
                success: function(response) {
                    sessionStorage.setItem('signupSuccess', 'true');
                    window.location.href = '/login.html';

                },
                error: function(xhr) {
                    signup.showError($('#email'), 'Error creating account. Email may already be registered.');
                }
            });
        }
    };

    signup.init();
});