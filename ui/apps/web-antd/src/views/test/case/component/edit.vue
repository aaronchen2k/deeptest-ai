<script lang="ts" setup>
import { computed } from 'vue';

import { Page } from '@vben/common-ui';

import { useCaseStore } from '#/views/test/case/store';

const caseStore = useCaseStore();
const model = computed(() => caseStore.caseModel);
</script>

<template>
  <Page auto-content-height class="case-main">
    {{ model }}
  </Page>
</template>

<style lang="less">
.case-edit {
}
</style>

<style scoped lang="less">
.case-edit {
}
</style>
