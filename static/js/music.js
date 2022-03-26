$(document).ready(function () {
  pauseResumeBtn = $(".discord-bot-music-play button.play-pause")
  pauseResumeBtn.click(function () {
    order = "pause"
    if (pauseResumeBtn.hasClass("paused")) {
      order = "resume"
    }
    fetch(window.location + "/" + order)
      .then(response => {
        if (!response.ok) {
          console.log("errorPauseResume")
        } else {
          console.log("okPauseResume")
        }
      })
    pauseResumeBtn.toggleClass("paused");
  });
  $(".discord-bot-music-append .right button").click(function () {
    order = "play"
    textBar = $(".discord-bot-music-append .append-queue-text")[0]
    musicToAppend = textBar.value
    textBar.value = ""
    fetch(window.location + "/" + order, {
      method: 'POST',
      headers: new Headers({
        'Content-Type': 'application/x-www-form-urlencoded',
      }),
      body: "param1=" + musicToAppend // <-- Post parameters
    })
      .then(response => {
        if (!response.ok) {
          console.log("errorPlay")
        } else {
          console.log("okPlay")
        }
      })
  });

  var loc = window.location, new_uri;
  if (loc.protocol === "https:") {
    new_uri = "wss:";
  } else {
    new_uri = "ws:";
  }
  new_uri += "//" + loc.host;
  new_uri += loc.pathname + "/ws";
  ws = new WebSocket(new_uri);
  ws.onopen = function (evt) {
    console.log("OPEN");
  }
  ws.onclose = function (evt) {
    console.log("CLOSE");
    ws = null;
  }
  ws.onmessage = function (evt) {
    console.log("RESPONSE: " + evt.data);
    obj = JSON.parse(evt.data)
    if (obj.PlayStatus != pauseResumeBtn.hasClass("paused")) {
      pauseResumeBtn.toggleClass("paused");
    }
    $(".currentTime").html(obj.CurrentTime);
    $(".currentTitle").html(obj.CurrentTitle);
  }
  ws.onerror = function (evt) {
    console.log("ERROR: " + evt.data);
  }
});
