<template>
  <div class="chat-window" ref="chatWindow">
    <div v-for="msg in messages" :key="msg.id" :class="['msg-wrapper', msg.role]">
  <div class="bubble">
    <!-- 只在 AI 消息显示角色名 -->
    <template v-if="msg.role==='ai'">
      <span class="role">{{ msg.selectedRole || 'AI助手' }}</span>:
    </template>
    {{ msg.text }}
    <button v-if="msg.role==='ai'" @click="playVoice(msg)" :disabled="msg.playing">
      {{ msg.playing ? '🔊 播放中...' : '🔊' }}
    </button>
  </div>
</div>
  </div>
</template>

<script setup>
import { ref, nextTick, watch } from 'vue'
import { useChatStore } from '../stores/chatStore'

const chatStore = useChatStore()
const messages = chatStore.messages
const chatWindow = ref(null)

// 页面滚动到底部
watch(messages, () => {
  nextTick(() => {
    if (chatWindow.value) chatWindow.value.scrollTop = chatWindow.value.scrollHeight
  })
})

// 预加载语音
const voices = speechSynthesis.getVoices()
if (!voices.length) {
  speechSynthesis.onvoiceschanged = () => {
    voices.splice(0, voices.length, ...speechSynthesis.getVoices())
  }
}

// 播放 AI 消息语音
function playVoice(msg) {
  if (!msg.text) return

  // 防止多条同时播放
  speechSynthesis.cancel()
  
  msg.playing = true

  // 选择角色语音
  const roleVoice = getVoiceForRole(msg.selectedRole || msg.role)

  // 拆分长文本，每100字符为一段
  const chunks = msg.text.match(/.{1,100}/g) || [msg.text]

  let index = 0

  function speakNext() {
    if (index >= chunks.length) {
      msg.playing = false
      return
    }
    const utterance = new SpeechSynthesisUtterance(chunks[index])
    utterance.voice = roleVoice
    utterance.onend = () => {
      index++
      speakNext()
    }
    speechSynthesis.speak(utterance)
  }

  speakNext()
}

// 根据角色返回语音对象
function getVoiceForRole(roleName) {
  // 这里可以按实际需求自定义角色语音
  // 示例：AI助手 → 英文女声, 小助手 → 英文男声
  if (!voices.length) return null
  if (roleName.includes('女')) return voices.find(v => v.lang.includes('en') && v.name.includes('Female')) || voices[0]
  if (roleName.includes('男')) return voices.find(v => v.lang.includes('en') && v.name.includes('Male')) || voices[0]
  return voices[0]
}
</script>

<style scoped>
.chat-window {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow-y: auto;
  padding: 16px;
}

.msg-wrapper {
  display: flex;
  margin-bottom: 12px;
}

.msg-wrapper.user {
  justify-content: flex-end; /* 用户靠右 */
}

.msg-wrapper.ai {
  justify-content: flex-start; /* AI靠左 */
}

.bubble {
  max-width: 60%;
  word-break: break-word;
  padding: 8px 12px;
  border-radius: 16px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.msg-wrapper.user .bubble {
  background: #d0f0ff;
}

.msg-wrapper.ai .bubble {
  background: #f1f0f0;
}


.bubble button {
  border: none;
  background: none;
  cursor: pointer;
  font-size: 16px;
}
</style>
