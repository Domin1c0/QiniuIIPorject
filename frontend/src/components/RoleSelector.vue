<template>
  <div>
    <h2>选择角色</h2>
    <input v-model="keyword" placeholder="搜索角色"/>
    <button @click="selectRole">确定</button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const keyword = ref('')

function selectRole() {
  alert(`选中角色：${keyword.value}`)
}
</script>
