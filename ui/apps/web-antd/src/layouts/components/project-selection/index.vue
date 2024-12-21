<script lang="ts" setup>
import '@vben/styles/global';

import { Dropdown, Menu, MenuItem, Button } from 'ant-design-vue';

function selectProject(item) {
  window.console.log('selectProject', item);
}
</script>

<template>
  <Dropdown class="dp-dropdown" overlay-class-name="dp-dropdown-overlay">
    <template #overlay>
      <Menu @click="selectProject">
        <MenuItem key="1" @click="selectProject(1)">
          <span class="dp-dropdown-icon icon-[ant-design--project-outlined]" />
          &nbsp;项目B
        </MenuItem>
      </Menu>
    </template>

    <span>
      <span>项目A</span>
      <span class="dp-dropdown-icon icon-[ant-design--down-outlined]"></span>
    </span>
  </Dropdown>
</template>
