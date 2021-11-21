$(document).ready(function () {
    $(".subnav").find("button.confirm").click(function (e) {
        e.preventDefault()
        urlParts = window.location.href.split('/')
        window.location.href = urlParts[0] + '/' + urlParts[1] + '/'  + urlParts[2] + '/'  + urlParts[3] + '/' + $(this)[0].innerHTML.replace(" ", '')
    });
});