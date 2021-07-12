$(document).ready(function () {
    $(".nav-group").click(function () {
        if ($(this).hasClass("expanded")) {
            $(this).removeClass("expanded")
            $(this).addClass("not-expanded")
            $(this).find(".nav-group-content").hide()
        } else {
            $(this).removeClass("not-expanded")
            $(this).addClass("expanded")
            $(this).find(".nav-group-content").show()
        }
    })
});