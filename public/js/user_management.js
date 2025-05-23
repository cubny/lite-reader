$(document).ready(function() {
    // Initialize the dialog
    var addUserDialog = $("#addUserModal").dialog({
        autoOpen: false,
        height: 400, // Adjusted height for better spacing
        width: 350,
        modal: true,
        buttons: {
            "Create User": function() {
                $("#addUserForm").submit();
            },
            Cancel: function() {
                addUserDialog.dialog("close");
            }
        },
        close: function() {
            $("#addUserForm")[0].reset();
            $("#newUserEmail").removeClass("ui-state-error");
            $("#newUserPassword").removeClass("ui-state-error");
            $("#confirmNewUserPassword").removeClass("ui-state-error");
            $("#addUserMessage").empty().removeClass("success-message error-message"); // Clear classes too
        }
    });

    // Handle "Add User" button click
    $("#adduser a.add").on("click", function(e) {
        e.preventDefault();
        addUserDialog.dialog("open");
    });

    // Handle form submission
    $("#addUserForm").on("submit", function(e) {
        e.preventDefault();
        var email = $("#newUserEmail").val().trim();
        var password = $("#newUserPassword").val();
        var confirmPassword = $("#confirmNewUserPassword").val();
        var isValid = true;

        // Basic validation
        $("#newUserEmail, #newUserPassword, #confirmNewUserPassword").removeClass("ui-state-error");
        $("#addUserMessage").empty().removeClass("success-message error-message");

        if (email === "" || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) { // Simple email regex
            $("#newUserEmail").addClass("ui-state-error");
            isValid = false;
            $("#addUserMessage").html("Please enter a valid email.").addClass("error-message");
            return;
        }
        if (password.length < 6) { // As per prompt for client-side validation
            $("#newUserPassword").addClass("ui-state-error");
            isValid = false;
            $("#addUserMessage").html("Password must be at least 6 characters.").addClass("error-message");
            return;
        }
        if (password !== confirmPassword) {
            $("#confirmNewUserPassword").addClass("ui-state-error");
            isValid = false;
            $("#addUserMessage").html("Passwords do not match.").addClass("error-message");
            return;
        }

        if (isValid) {
            var authToken = getCookie("access_token"); // Assumes getCookie is defined elsewhere (e.g. utils.js or main.js)
            
            $("#addUserMessage").html("Creating user...").removeClass("error-message success-message");

            $.ajax({
                url: '/api/users',
                type: 'POST',
                contentType: 'application/json',
                data: JSON.stringify({ email: email, password: password }),
                beforeSend: function(xhr) {
                    if (authToken) {
                        xhr.setRequestHeader('Authorization', 'Bearer ' + authToken);
                    } else {
                        // Handle case where auth token is not found, though ideally UI should prevent this state
                        $("#addUserMessage").html("Authentication token not found. Please log in again.").addClass("error-message");
                        isValid = false; // Prevent further processing
                        return false; // Cancel AJAX request
                    }
                },
                success: function(response) {
                    $("#addUserMessage").html("User created successfully!").addClass("success-message").removeClass("error-message");
                    setTimeout(function(){ addUserDialog.dialog("close"); }, 2000);
                },
                error: function(xhr) {
                    var errorMsg = "Error creating user.";
                    if (xhr.responseJSON && xhr.responseJSON.error && xhr.responseJSON.error.details) { // Adjusted to new error structure
                        errorMsg = xhr.responseJSON.error.details;
                    } else if (xhr.responseJSON && xhr.responseJSON.message) { // Fallback for older error structure
                        errorMsg = xhr.responseJSON.message;
                    } else if (xhr.responseText) {
                        try {
                            var resp = JSON.parse(xhr.responseText);
                            if (resp.error && resp.error.details) {
                                errorMsg = resp.error.details;
                            } else if (resp.message) {
                                errorMsg = resp.message;
                            }
                        } catch (e) { /* ignore if responseText is not JSON */ }
                    }
                    $("#addUserMessage").html(errorMsg).addClass("error-message").removeClass("success-message");
                }
            });
        }
    });

    // Adding some basic styling for messages, can be moved to CSS
    $("<style type='text/css'> .error-message { color: red; } .success-message { color: green; } </style>").appendTo("head");
});
