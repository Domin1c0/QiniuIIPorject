<template>
  <div class="role-selector">
    <h3>选择角色</h3>
    <ul>
      <li v-for="role in roles" :key="role" 
          :class="{ active: role === currentRole }" 
          @click="select(role)">
        {{ role }}
      </li>
    </ul>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useChatStore } from '../stores/chatStore'

const chatStore = useChatStore()
const roles = ['otto', '塔菲', '洞穴里安']
const currentRole = computed(() => chatStore.selectedRole)

function select(role) {
  chatStore.setRole(role)
}
</script>

<style scoped>
.role-selector { width: 200px; background:#fff; border-right:1px solid #ddd; padding:16px; }
.role-selector ul { list-style:none; padding:0; }
.role-selector li { padding:8px 12px; margin-bottom:8px; cursor:pointer; border-radius:6px; transition: background 0.2s; }
.role-selector li.active, li:hover { background:#007bff; color:white; }
</style>
