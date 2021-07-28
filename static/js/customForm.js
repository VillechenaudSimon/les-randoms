$(document).ready(function () {
    $(".customForm").find("button.confirm").click(function (e) { 
        e.preventDefault()
        $(this).closest(".customForm").find(".customFormValue").attr("value", $(this).text())
        $(this).parent().submit()
    });
});