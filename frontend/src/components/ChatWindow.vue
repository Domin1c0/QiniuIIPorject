<template>
  <div>
    <h2>聊天窗口</h2>
    <ul>
      <li v-for="(msg, idx) in messages" :key="idx">{{ msg }}</li>
    </ul>
    <input v-model="newMsg" placeholder="输入消息…" />
    <button @click="send">发送</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const messages = ref<string[]>([])
const newMsg = ref('')

function send() {
  if (newMsg.value.trim() === '') return
  messages.value.push(newMsg.value)
  newMsg.value = ''
}
</script>

<style scoped>
ul {
  list-style: none;
  padding: 0;
}
li {
  background: #f3f3f3;
  margin-bottom: 4px;
  padding: 4px 8px;
  border-radius: 4px;
}
</style>
