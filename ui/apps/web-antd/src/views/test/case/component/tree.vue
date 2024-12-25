<script setup lang="ts">
import { computed, unref } from 'vue';

import { filterByKeyword } from '@vben/utils';

import { InputSearch, Spin, Tree } from 'ant-design-vue';
import cloneDeep from 'lodash/cloneDeep';

import { DropdownActionMenu } from '#/component/DropDownMenu';
import { type MenuItem } from '#/component/DropDownMenu/type';
import { useCaseStore } from '#/views/test/case/store';

const caseStore = useCaseStore();
const treeData = computed(() => caseStore.treeData);

const keywords = computed(() => caseStore.keywords);

caseStore.fetchTreeData();

const treeDataComputed = computed(() => {
  const data = cloneDeep(unref(treeData));
  if (data?.length > 0) {
    return [...filterByKeyword(data, keywords.value, 'title')];
  }
  return [];
});

const DropdownMenuList = [
  {
    label: '新建目录',
    ifShow: (nodeProps) => nodeProps.data.type !== 'leaf',
    action: (nodeProps) => caseStore.createNode(nodeProps.data?.id, 'dir'),
  },
  {
    label: '新建用例',
    ifShow: (nodeProps) => nodeProps.data?.type !== 'leaf',
    action: (nodeProps) => caseStore.createNode(nodeProps.data?.id, 'case'),
  },
  {
    label: (nodeProps: any) => {
      return `编辑${nodeProps.data.type === 'leaf' ? '用例' : '目录'}`;
    },
    ifShow: (nodeProps: any) => nodeProps.data.type !== 'root',
    action: (nodeProps: any) => caseStore.editNode(nodeProps),
  },
  {
    label: (nodeProps: any) => {
      return `删除${nodeProps.data.type === 'leaf' ? '用例' : '目录'}`;
    },
    auth: 'p-api-debug-del',
    ifShow: (nodeProps) => nodeProps.data.type !== 'root',
    action: (nodeProps) => caseStore.deleteNode(nodeProps),
  },
] as MenuItem[];

// 根据搜索关键词搜索结果为空展示
const showKeywordsTip = computed(() => {
  return keywords.value && treeDataComputed.value.length === 0;
});
</script>

<template>
  <div class="case-tree-main">
    <div class="dp-tree-container">
      <div class="tree-filter">
        <InputSearch
          v-model:value="keywords"
          class="search-input"
          placeholder="输入关键字过滤"
        />
        <div class="add-btn" @click="caseStore.createNode(0, 'dir')">
          <span
            class="icon-[ant-design--plus-outlined]"
            style="font-size: 16px"
          ></span>
        </div>
      </div>

      <div class="tree-content">
        <Tree
          :auto-expand-parent="caseStore.autoExpandParent"
          :expanded-keys="caseStore.expandedKeys"
          :field-names="{ key: 'id' }"
          :selected-keys="caseStore.selectedKeys"
          :show-icon="true"
          :tree-data="treeDataComputed"
          block-node
          class="dp-tree"
          draggable
          @drop="caseStore.onDrop"
          @expand="caseStore.onExpand"
          @select="caseStore.selectNode"
        >
          <template #switcherIcon="nodeProps">
            <span
              v-if="!nodeProps.expanded"
              class="icon-[ant-design--caret-right-outlined]"
            ></span>
            <span v-else class="icon-[ant-design--caret-down-outlined]"></span>
          </template>

          <template #title="nodeProps">
            <div
              :draggable="nodeProps.data.type !== 'root'"
              :title="nodeProps.data.title"
              class="tree-title"
            >
              <span class="tree-icon">
                <span
                  v-if="nodeProps.data.type !== 'leaf' && !nodeProps.expanded"
                  class="icon-[ant-design--folder-outlined] dp-antdv-icon"
                >
                </span>
                <span
                  v-if="nodeProps.data.type !== 'leaf' && nodeProps.expanded"
                  class="icon-[ant-design--folder-open-outlined] dp-antdv-icon"
                >
                </span>
              </span>

              <span
                v-if="nodeProps.data.title.indexOf(keywords) > -1"
                class="tree-title-text"
              >
                <span>{{
                  nodeProps.data.title.substring(
                    0,
                    nodeProps.data.title.indexOf(keywords),
                  )
                }}</span>
                <span style="color: #f50">{{ keywords }}</span>
                <span>{{
                  nodeProps.data.title.substring(
                    nodeProps.data.title.indexOf(keywords) + keywords.length,
                  )
                }}</span>
              </span>
              <span v-else class="tree-title-text">{{
                nodeProps.data.title
              }}</span>

              <span v-if="nodeProps.data.id > 0" class="more-icon">
                <DropdownActionMenu
                  :dropdown-list="DropdownMenuList"
                  :record="nodeProps"
                />
              </span>
            </div>
          </template>
        </Tree>

        <div v-if="treeDataComputed.length === 0" class="loading-container">
          <div v-if="showKeywordsTip" class="nodata-tip">搜索结果为空 ~</div>
          <div
            v-else-if="!caseStore.loading && !showKeywordsTip"
            class="nodata-tip"
          >
            请点击上方按钮添加目录 ~
          </div>
          <Spin v-else style="margin-top: 20px" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">
.case-tree-main {
  height: 100%;

  .loading-container {
    display: flex;
    align-items: center;
    justify-content: center;
  }

  :deep(.ant-tree-node-content-wrapper) {
    width: 100%;
    display: inline-flex;
    align-items: center;

    .ant-tree-title {
      width: 100%;
      display: inline-flex;
      align-items: center;
    }
  }

  .tree-title {
    display: inline-flex;
    width: 100%;
    align-items: center;

    .tree-icon {
      margin-right: 4px;
    }
  }

  .nodata-tip {
    margin-left: 0 !important;
  }
}
</style>
