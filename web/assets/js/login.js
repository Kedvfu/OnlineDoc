const loginForm = $('#loginForm');

$('#loginBtn').on('click', function(e) {
    e.preventDefault();

    const url = loginForm.attr('action');
    const data = loginForm.serialize();

    $.post(url , data, function(response) {
        //alert(response.message);
        if (response.session_token != null) {
            setCookie('session_token', response.session_token, 7);
            setCookie('user_id', response.user_id, 7);

            window.location.href = '/home';
        }
        $("#error-message").text(response.message)
    });
});
