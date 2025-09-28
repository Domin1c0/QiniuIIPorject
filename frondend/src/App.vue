<template>
  <div class="app-container">
    <RoleSelector />
    <div class="chat-area">
      <ChatWindow />
      <ChatInput />
    </div>
  </div>
</template>

<script setup>
import RoleSelector from './components/RoleSelector.vue'
import ChatWindow from './components/ChatWindow.vue'
import ChatInput from './components/ChatInput.vue'
</script>

<style scoped>
.app-container {
  display: flex;
  height: 100vh;
  background: #f0f2f5;
}

.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 16px;
}
</style>
