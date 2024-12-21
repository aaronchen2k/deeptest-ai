<script lang="ts" setup>
import '@vben/styles/global';

import { Dropdown, Menu, MenuItem } from 'ant-design-vue';

function selectProject() {
  window.console.log('selectProject');
}
</script>

<template>
  <Dropdown>
    <template #overlay>
      <Menu @click="selectProject">
        <MenuItem key="1">
          <span
            class="icon-[ant-design--project-outlined]"
            style="mask-position: 0 3px"
          ></span>
          &nbsp;项目A
        </MenuItem>
        <MenuItem key="2">
          <span
            class="icon-[ant-design--project-outlined]"
            style="mask-position: 0 3px"
          ></span>
          &nbsp;项目B
        </MenuItem>
        <MenuItem key="3">
          <span
            class="icon-[ant-design--project-outlined]"
            style="mask-position: 0 3px"
          ></span>
          &nbsp;项目C
        </MenuItem>
      </Menu>
    </template>

    <a-button class="dp-link">
      <span>项目A</span>&nbsp;
      <span
        class="icon-[ant-design--down-outlined]"
        style="mask-position: 2px 3px; font-weight: bold"
      ></span>
    </a-button>
  </Dropdown>
</template>
