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
    //console.log("OPEN");
  }
  ws.onclose = function (evt) {
    //console.log("CLOSE");
    ws = null;
  }
  ws.onmessage = function (evt) {
    //console.log("RESPONSE: " + evt.data);
    obj = JSON.parse(evt.data)
    if (obj.DataType == 1) {
      updateQueueDisplay($(".discord-bot-music-play .body")[0], obj.Queue)
    }
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

function updateQueueDisplay(body, queueData) {
  //console.log(queueData.length)
  let i = 0
  let j = 1
  let divs = body.children
  while (j < queueData.length || i < divs.length) {
    if (i >= divs.length) {
      body.appendChild(newQueueElt(queueData[j].Title))
      ++i
      ++j
    } else if (j >= queueData.length || divs[i].children[0].firstChild.nodeValue != queueData[j].Title) {
      //if (j < queueData.length) {
      //  console.log("Queue Changed  : " + divs[i].children[0].firstChild.nodeValue + " != " + queueData[j].Title)
      //}
      body.removeChild(divs[i])
    } else {
      ++i
      ++j
    }
  }

  /*
  body.empty()
  for(let i = 1; i < queueData.length; ++i) {
    body.appendChild(newQueueElt([i].Title))
  }
  */
}

function newQueueElt(title) {
  let e = document.createElement("div")
  e.innerHTML = "<span>" + title + "</span>"
  return e
}