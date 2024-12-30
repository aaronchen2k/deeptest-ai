<script lang="ts" setup>
import { ref } from 'vue';

import Chatbot from './component/Chat/index.vue';

const chatbotServerUrl = ref(import.meta.env.VITE_GLOB_API_SERVER);
</script>

<template>
  <div class="chatbot-main">
    <Chatbot :server-url="chatbotServerUrl">
      <template #actions>
        <span class="icon-container" title="帮助">
          <span class="help icon dp-link"></span>
        </span>
        <span class="icon-container" title="评论">
          <span class="comments icon dp-link"></span>
        </span>
        <span class="icon-container" title="点赞">
          <span class="up icon dp-link"></span>
        </span>
        <span class="icon-container" title="拍砖">
          <span class="down icon dp-link"></span>
        </span>
      </template>
    </Chatbot>
  </div>
</template>

<style lang="less">
.chatbot-main {
}
</style>

<style lang="less" scoped>
.chatbot-main {
  padding: 16px;
}
</style>
