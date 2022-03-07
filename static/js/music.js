$(document).ready(function () {
  var btn = $(".discord-bot-music-play button.play-pause");
  btn.click(function () {
    order = "pause"
    if (btn.hasClass("paused")) {
      order = "resume"
    }
    fetch(window.location + "/" + order)
      .then(response => {
        if(!response.ok) {
          console.log("error")
        } else {
          console.log("ok")
        }
      })
    btn.toggleClass("paused");
    return false;
  });
});
