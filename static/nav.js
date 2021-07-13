$(document).ready(function () {
    $(".nav-group-header").click(function () {
        if ($(this).parent().hasClass("expanded")) {
            $(this).parent().removeClass("expanded")
            $(this).parent().addClass("not-expanded")
            $(this).parent().find(".nav-group-content").hide()
        } else {
            $(this).parent().removeClass("not-expanded")
            $(this).parent().addClass("expanded")
            $(this).parent().find(".nav-group-content").show()
        }
    })
});