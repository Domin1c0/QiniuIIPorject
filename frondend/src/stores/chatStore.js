import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useChatStore = defineStore('chat', () => {
  const messages = ref([])
  const selectedRole = ref('AI助手') // 默认角色

  function addMessage(msg) {
    messages.value.push({ ...msg, id: Date.now() + Math.random() })
  }

  function setRole(role) {
    selectedRole.value = role
  }

  return { messages, selectedRole, addMessage, setRole }
})
