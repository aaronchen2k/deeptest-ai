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
          <span
            class="dp-dropdown-icon icon-[ant-design--project-outlined]"
          ></span>
          <span class="dp-inline-block" style="padding-left: 10px">项目B</span>
        </MenuItem>
      </Menu>
    </template>

    <span>
      <span class="dp-inline-block" style="padding-right: 8px">项目A</span>
      <span class="dp-dropdown-icon icon-[ant-design--down-outlined]"></span>
    </span>
  </Dropdown>
</template>
