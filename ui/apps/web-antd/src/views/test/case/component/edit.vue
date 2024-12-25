<script lang="ts" setup>
import { ref } from 'vue';

import { Page } from '@vben/common-ui';

import { Dropdown, Empty, Menu, MenuItem, TabPane, Tabs } from 'ant-design-vue';

import { useCaseStore } from '#/views/test/case/store';

const caseStore = useCaseStore();

/**
 * tabs 右键菜单操作
 */
const visible = ref({} as any);
const dropdownVisible = ref(false); // 三点菜单
const cancelVisible = () => {
  window.console.log('cancelVisible');
  caseStore.caseTabs.forEach((e: any) => {
    visible.value[e.id] = false;
  });
};

const openDropdown = (item: any) => {
  caseStore.caseTabs.forEach((e: any) => {
    visible.value[e.id] = false;
  });
  visible.value[item.id] = true;
};

const onContextMenuClick = (e: any, record?: any) => {
  dropdownVisible.value = false;
  switch (e.key) {
    case 'close_all': {
      caseStore.setCaseTabs([]);
      caseStore.setCaseModel(null);
      break;
    }
    case 'close_cur': {
      if (!record || record?.id === caseStore.caseModel.id) {
        caseStore.removeCaseTab(caseStore.caseModel.id);
      } else {
        caseStore.removeCaseTab(record.id);
      }
      break;
    }
    case 'close_other': {
      if (!record || record?.id === caseStore.caseModel.id) {
        const caseTabs = caseStore.caseTabs.filter(
          (e) => e.id === caseStore.caseModel.id,
        );
        caseStore.setCaseTabs(caseTabs);
      } else {
        caseStore.setCaseTabs([record]);
        caseStore.openCaseTab(record);
      }
      break;
    }
    default: {
      break;
    }
  }
};

const getTitle = (title: string) => {
  const len = title.length;
  if (len <= 12) return title;

  // eslint-disable-next-line unicorn/prefer-string-slice
  return `${title.slice(0, 16)}...${title.substring(len - 6, len)}`;
};

const tabsContextMenu = [
  {
    key: 'close_cur',
    label: '关闭当前标签页',
  },
  {
    key: 'close_other',
    label: '关闭其他标签页',
  },
  {
    key: 'close_all',
    label: '关闭所有标签页',
  },
];
</script>

<template>
  <Page auto-content-height class="case-main">
    <Tabs
      v-if="caseStore.caseTabs?.length"
      :active-key="caseStore.caseModel.id"
      :closable="true"
      class="dp-tabs-full-height"
      type="editable-card"
      @change="caseStore.changeTab"
      @edit="caseStore.onTabMenuClicked"
    >
      <TabPane
        v-for="tab in caseStore.caseTabs"
        :key="tab.id"
        :title="tab.title"
        class="dp-relative"
      >
        <template #tab>
          <Dropdown :trigger="['contextmenu']" :visible="visible[tab.id]">
            <div
              v-on-click-outside="cancelVisible"
              @contextmenu="openDropdown(tab)"
            >
              <span :title="tab.title">{{ getTitle(tab.title) }}</span>
            </div>
            <template #overlay>
              <Menu @click="(e: any) => onContextMenuClick(e, tab)">
                <MenuItem v-for="item in tabsContextMenu" :key="item.key">
                  {{ item.label }}
                </MenuItem>
              </Menu>
            </template>
          </Dropdown>
        </template>

        <div class="interface-tabs-content">FORM</div>
      </TabPane>

      <template #addIcon>
        <div
          :class="[dropdownVisible ? 'visible' : '']"
          class="extra-menu"
          @mouseenter="dropdownVisible = true"
          @mouseleave="dropdownVisible = false"
        >
          <span style="cursor: pointer">
            <span class="icon-[ant-design--ellipsis-outlined]"></span>
          </span>
          <Menu :selected-keys="[]" @click="(e: any) => onContextMenuClick(e)">
            <MenuItem v-for="item in tabsContextMenu" :key="item.key">
              {{ item.label }}
            </MenuItem>
          </Menu>
        </div>
      </template>
    </Tabs>

    <div v-else style="margin-top: 36px">
      <Empty description="请在左侧选择用例。" />
    </div>
  </Page>
</template>
