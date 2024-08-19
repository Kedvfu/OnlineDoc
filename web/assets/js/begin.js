$(document).ready(function() {
    $.ajaxSetup({
        xhrFields: {
            withCredentials: true

        }
    })

})
// var pageData = {
//     "user_id" : getCookie("user_id"),
//
//     "document_id": "",
//     "title": ""
//
// }
function setCookie(name, value, days) {
    var expires = "";
    if (days) {
        var date = new Date();
        date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
        expires = "; expires=" + date.toUTCString();
    }
    document.cookie = name + "=" + (value || "") + expires + "; path=/";
}
function getCookie(name) {
    var nameEQ = name + "=";
    var ca = document.cookie.split(';');
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) === ' ') c = c.substring(1, c.length);
        if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
}


if (location.href.indexOf("/document/") >= 0) {
    pageData.document_id = location.href.substring(location.href.indexOf("/document/")).split("/")[2];
}
pageData.title = document.title;