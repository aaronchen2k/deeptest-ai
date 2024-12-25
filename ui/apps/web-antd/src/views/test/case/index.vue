<script lang="ts" setup>
import { Page } from '@vben/common-ui';

import ContentPane from '#/views/component/ContentPane/index.vue';

import Edit from './component/edit.vue';
import Tree from './component/tree.vue';
</script>

<template>
  <Page auto-content-height class="case-main">
    <ContentPane
      :container-style="{ padding: 0, margin: 0 }"
      :show-expand="true"
    >
      <template #left>
        <Tree />
      </template>
      <template #right>
        <Edit />
      </template>
    </ContentPane>
  </Page>
</template>
