$(document).ready(function () {
  var btn = $(".discord-bot-music-play button.play-pause");
  btn.click(function () {
    order = "pause"
    if (btn.hasClass("paused")) {
      order = "resume"
    }
    fetch(window.location + "/" + order)
      .then(response => {
        if (!response.ok) {
          console.log("error")
        } else {
          console.log("ok")
        }
      })
    btn.toggleClass("paused");
    return false;
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
    //console.log("OPEN");
  }
  ws.onclose = function (evt) {
    //console.log("CLOSE");
    ws = null;
  }
  ws.onmessage = function (evt) {
    //console.log("RESPONSE: " + evt.data);
    obj = JSON.parse(evt.data)
    if (obj.PlayStatus != btn.hasClass("paused")) {
      btn.toggleClass("paused");
    }
    $(".currentTime").html(obj.CurrentTime);
  }
  ws.onerror = function (evt) {
    console.log("ERROR: " + evt.data);
  }
});
