<script lang="ts" setup>
import { computed } from 'vue';

import { Page } from '@vben/common-ui';

import { useCaseStore } from '#/views/test/case/store';

const caseStore = useCaseStore();
const model = computed(() => caseStore.designModel);
</script>

<template>
  <Page class="case-edit-main">
    {{ model }}
  </Page>
</template>

<style scoped lang="less">
.case-edit-main {
  height: 100%;
}
</style>
