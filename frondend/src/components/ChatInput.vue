<template>
  <div class="chat-input-wrapper">
    <div class="chat-input">
      <!-- 语音按钮 -->
      <button class="voice-btn" @click="toggleRecording">
        {{ recording ? '🎤 录音中...' : '🎤' }}
      </button>

      <!-- 输入框 -->
      <input
        v-model="inputText"
        @keyup.enter="sendMessage"
        placeholder="请输入消息或点击语音按钮录制..."
      />

      <!-- 发送按钮 -->
      <button class="send-btn" @click="sendMessage">发送</button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useChatStore } from '../stores/chatStore'

const chatStore = useChatStore()
const inputText = ref('')
const recording = ref(false)

let mediaRecorder
let audioChunks = []
let recognition
let currentStream // 保存当前流

async function sendMessage() {
  if (!inputText.value.trim()) return

  chatStore.addMessage({
    role: 'user',
    selectedRole: chatStore.selectedRole,
    text: inputText.value
  })

  const res = await fetch('/chat/new', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ message: { content: inputText.value } })
  })
  const data = await res.json()

  chatStore.addMessage({
    role: data.role,
    selectedRole: chatStore.selectedRole,
    text: data.content
  })

  inputText.value = ''
}

function toggleRecording() {
  if (!navigator.mediaDevices || !navigator.mediaDevices.getUserMedia) {
    alert('浏览器不支持语音录制')
    return
  }

  // 如果在录音则停止
  if (mediaRecorder && mediaRecorder.state === 'recording') {
    mediaRecorder.stop()
    return
  }

  // 开始录音
  navigator.mediaDevices.getUserMedia({ audio: true }).then(stream => {
    currentStream = stream
    mediaRecorder = new MediaRecorder(stream)
    audioChunks = []
    recording.value = true

    mediaRecorder.ondataavailable = e => audioChunks.push(e.data)
    mediaRecorder.onstop = async () => {
      // 停止占用麦克风
      if (currentStream) {
        currentStream.getTracks().forEach(track => track.stop())
        currentStream = null
      }
      recording.value = false

      const blob = new Blob(audioChunks, { type: 'audio/webm' })
      const formData = new FormData()
      formData.append('file', blob, 'voice.webm')

      try {
        const res = await fetch('/speech-to-text', {
          method: 'POST',
          body: formData
        })
        const result = await res.json()

        // 只把识别出的文字放到输入框
        if (result && result.text) {
          inputText.value = result.text
        } else {
          alert('语音识别失败或返回为空')
        }
      } catch (err) {
        console.error(err)
        alert('语音识别接口出错')
      }
    }

    mediaRecorder.start()
  })
}
</script>


<style scoped>
.chat-input-wrapper {
  width: 50%;
  margin: 0 auto;
  padding: 8px 0;
}

.chat-input {
  display: flex;
  align-items: center;
  gap: 8px;
}

.chat-input input {
  flex: 1;
  padding: 10px 16px;
  border-radius: 24px;
  border: 1px solid #ccc;
  outline: none;
  font-size: 14px;
  transition: all 0.2s;
}

.chat-input input:focus {
  border-color: #007bff;
  box-shadow: 0 0 6px rgba(0,123,255,0.3);
}

.voice-btn {
  padding: 10px;
  border-radius: 50%;
  border: none;
  background: #f0f0f0;
  cursor: pointer;
  transition: all 0.2s;
}

.voice-btn:hover {
  background: #d0d0d0;
}

.send-btn {
  padding: 10px 16px;
  border-radius: 24px;
  border: none;
  background: #007bff;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
}

.send-btn:hover {
  background: #0056b3;
}
</style>
