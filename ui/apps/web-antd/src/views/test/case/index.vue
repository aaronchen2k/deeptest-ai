<script lang="ts" setup>
import { computed } from 'vue';

import { Page } from '@vben/common-ui';

import { useCaseStore } from '#/views/test/case/store';

const caseStore = useCaseStore();
const treeData = computed(() => caseStore.treeData);

caseStore.fetchTreeData();
</script>

<template>
  <Page auto-content-height>
    {{ treeData }}
  </Page>
</template>
