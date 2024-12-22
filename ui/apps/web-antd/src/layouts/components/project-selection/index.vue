<script lang="ts" setup>
import {computed, onMounted, onUnmounted, ref} from 'vue';

import '@vben/styles/global';

import { DropdownButton, Menu, MenuItem } from 'ant-design-vue';

import { useGlobalStore } from '#/store/global';

const globalStore = useGlobalStore();
const currProject = computed(() => globalStore.currProject);
const projects = computed(() => globalStore.projects);

globalStore.loadUserProjects();

function selectProject(item: any) {
  window.console.log('selectProject', item);
  globalStore.updateUserProject(item);
}

onMounted(() => {
  window.console.log('onMounted');
});
onUnmounted(() => {
  window.console.log('onUnmounted');
});
</script>

<template>
  <div>
    <span class="dp-inline-block" style="padding-right: 20px">项目</span>

    <DropdownButton overlay-class-name="dp-dropdown-overlay">
      <span>
        {{ currProject.name }}
      </span>

      <template #overlay>
        <Menu>
          <template v-for="project in projects" :key="project.id">
            <MenuItem
              :disabled="currProject.id === project.id"
              @click="selectProject(project)"
            >
              <span
                class="dp-dropdown-icon icon-[ant-design--project-outlined]"
              ></span>
              <span class="dp-inline-block" style="padding-left: 10px">
                {{ project.name }}
              </span>
            </MenuItem>
          </template>
        </Menu>
      </template>
    </DropdownButton>
  </div>
</template>
