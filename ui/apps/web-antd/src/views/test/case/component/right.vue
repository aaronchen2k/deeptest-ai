<script lang="ts" setup>
import { ref } from 'vue';

import { Page } from '@vben/common-ui';

import { vOnClickOutside } from '@vueuse/components';
import { Dropdown, Empty, Menu, MenuItem, TabPane, Tabs } from 'ant-design-vue';

import { useCaseStore } from '#/views/test/case/store';

import Design from './design.vue';

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
      caseStore.setDesignModel(null);
      break;
    }
    case 'close_cur': {
      if (!record || record?.id === caseStore.designModel.id) {
        caseStore.removeCaseTab(caseStore.designModel.id);
      } else {
        caseStore.removeCaseTab(record.id);
      }
      break;
    }
    case 'close_other': {
      if (!record || record?.id === caseStore.designModel.id) {
        const caseTabs = caseStore.caseTabs.filter(
          (e) => e.id === caseStore.designModel.id,
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
  window.console.log('getTitle');

  const len = title?.length;
  if (!title || len <= 12) return title;

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
  <Page class="case-main">
    <Tabs
      v-if="caseStore.caseTabs?.length > 0"
      :active-key="caseStore.designModel.id"
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
          <Dropdown :open="visible[tab.id]" :trigger="['contextmenu']">
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

        <div class="interface-tabs-content">
          <Design />
        </div>
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

<style lang="less">
.case-main {
  .ant-tabs .ant-tabs-tab.ant-tabs-tab-active .ant-tabs-tab-btn {
    color: rgba(242, 242, 242, 0.85);
  }

  .ant-tabs-nav-wrap {
    overflow: visible !important;
    scrollbar-width: none;
    z-index: 999;
    width: 100%;

    &.ant-tabs-nav-wrap-ping-right,
    &.ant-tabs-nav-wrap-ping-left {
      overflow: hidden;
    }
  }

  .ant-tabs-nav-more {
    display: none;
  }

  .extra-menu {
    position: relative;

    &.visible {
      .ant-menu {
        height: max-content;
        transition: all 0.3s ease-in-out;
        opacity: 1;
      }
    }

    .ant-menu {
      position: absolute;
      width: 126px;
      right: 16px;
      top: 12px;
      height: 0;
      border-radius: 8px;
      box-shadow:
        0 6px 16px 0 rgba(0, 0, 0, 0.08),
        0 3px 6px -4px rgba(0, 0, 0, 0.12),
        0 9px 28px 8px rgba(0, 0, 0, 0.05);
      overflow: hidden;
      transition: all 0.3s ease-in-out;
      opacity: 0;
      z-index: 9999;
      .ant-menu-item {
        padding: 0 0 0 3px;
        margin: 0;
        width: 100% !important;
        line-height: 32px;
        height: 32px !important;
        text-align: center;

        &:first-child {
          margin-top: 4px;
        }

        &:last-child {
          margin-bottom: 4px;
        }
      }
    }
  }
}
</style>

<style scoped lang="less">
.case-main {
  height: 100%;
}
</style>
