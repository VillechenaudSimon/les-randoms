$(document).ready(function () {
    $(".customHoverMenu").parent().css("position", "relative")
    $(".customHoverMenu").parent().hover(function () {
            $(this).find(".customHoverMenu").css("display","flex")
        }, function () {
            $(this).find(".customHoverMenu").hide()
        }
    );
});