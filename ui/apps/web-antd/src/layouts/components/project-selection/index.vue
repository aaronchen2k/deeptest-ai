<script lang="ts" setup>
import {computed, ref} from 'vue';

import '@vben/styles/global';

import { Dropdown, Menu, MenuItem } from 'ant-design-vue';

import { loadProjectsApi, updateUserProject } from '#/api';
import { useGlobalStore } from '#/store/global';

const globalStore = useGlobalStore();
const currProject = computed(() => globalStore.currProject);
const projects = computed(() => globalStore.projects);

loadProjectsApi().then((result) => {
  window.console.log('loadProjectsApi', result);
  globalStore.setCurrProject(result.default);
  globalStore.setProjects(result.items);
});

function selectProject(item: any) {
  window.console.log('selectProject', item);
  updateUserProject(item.id).then(() => {
    globalStore.setCurrProject(item);
  });
}
</script>

<template>
  <Dropdown class="dp-dropdown" overlay-class-name="dp-dropdown-overlay">
    <template #overlay>
      <Menu>
        <template v-for="project in projects" :key="project.id">
          <MenuItem
            v-if="currProject.id !== project.id"
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

    <span>
      <span class="dp-inline-block" style="padding-right: 8px">
        {{ currProject.name }}
      </span>
      <span class="dp-dropdown-icon icon-[ant-design--down-outlined]"></span>
    </span>
  </Dropdown>
</template>
