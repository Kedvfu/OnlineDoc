const registerForm = $('#registerForm');

$('#registerBtn').on('click', function(e) {
    e.preventDefault();
    const url = registerForm.attr('action');
    const data = registerForm.serialize();

    $.post(url, data, function(response) {
        $("#error-message").text(response.message)
        if (response.success === 1) {
            window.location.href = "/login"
        }
    });
})
