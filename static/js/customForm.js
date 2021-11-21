$(document).ready(function () {
    $(".customForm").find("button.confirm").click(function (e) { 
        e.preventDefault()
        window.location.href = $(this)[0].baseURI + '/' + $(this).closest(".customForm").find(".customFormValue").val()
    });
});