
export function useSpeechRecognition(onResult) {
  const recognition = new (window.SpeechRecognition || window.webkitSpeechRecognition)()
  recognition.lang = 'zh-CN'
  recognition.interimResults = false

  recognition.onresult = (event) => {
    const text = event.results[0][0].transcript
    onResult(text)
  }

  const start = () => recognition.start()
  const stop = () => recognition.stop()

  return { start, stop }
}

export function speak(text) {
  const utter = new SpeechSynthesisUtterance(text)
  utter.lang = 'zh-CN'
  window.speechSynthesis.speak(utter)
}
