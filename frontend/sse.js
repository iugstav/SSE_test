let eventSource
let message = []

const messagesPreWrapper = document.querySelector("#messages_wrapper")

let messageProxy = new Proxy(message, {
  set: function(target, property, value) {
    target[property] = value

    const messageTag = document.createElement("p")
    messageTag.textContent = value

    messagesPreWrapper.appendChild(messageTag)

    return true
  }
})

window.addEventListener("DOMContentLoaded", () => {
  if (!eventSource) {
    eventSource = new EventSource("http://localhost:8080/logs")
    messageProxy[messageProxy.length - 1] = "streaming started"

    eventSource.onmessage = (ev) => {
      let data = JSON.parse(ev.data)
      console.log(data)
      messageProxy[messageProxy.length - 1] = `[LOG|${data.priority}]: ${data.value}`
    }
  }
})