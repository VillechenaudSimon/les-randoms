$(document).ready(function () {
    $(".subnav").find("button.confirm").click(function (e) {
        e.preventDefault()
        urlParts = window.location.href.split('/')
        newUrl = urlParts[0]
        for(i = 1; /*i < urlParts.length-1*/isNotSubNavElt($(this).parent(), urlParts[i]); i++) {
            newUrl += '/' + urlParts[i]
        }
        newUrl += '/' + escapeSubNavItemToURL($(this)[0].innerHTML)
        window.location.href = newUrl
    });
});

function isNotSubNavElt(subNavForm, toTest) {
    buttons = subNavForm.children()
    for(j = 0; j < buttons.length; j++) {
        if(escapeSubNavItemToURL(buttons[j].innerHTML) == toTest) {
            return false
        }
    }
    return true
}

function escapeSubNavItemToURL(s) {
    res = s.replace(/ /g, '')
    idx = res.indexOf("(")
    if (idx > 0) {
        res = res.substring(0, idx)
    }
    return res
}