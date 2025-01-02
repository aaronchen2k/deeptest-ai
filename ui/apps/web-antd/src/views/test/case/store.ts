import type { AntTreeNodeDropEvent } from 'ant-design-vue/lib/tree';

import { ref, watch } from 'vue';

import { filterTreeNodes, genNodeMap, isInArray } from '@vben/utils';

import { message } from 'ant-design-vue';
import { debounce } from 'lodash';
import { acceptHMRUpdate, defineStore } from 'pinia';

import { getCaseApi, loadCaseApi } from '#/api/test/case';
import { useGlobalStore } from '#/store/global';
import {
  getExpandedKeysCache,
  getSelectedKeyCache,
  setExpandedKeysCache,
  setSelectedKeyCache,
} from '#/utils/cache';
import { confirmToDelete } from '#/utils/confirm';

export const useCaseStore = defineStore('case', () => {
  const globalStore = useGlobalStore();

  const keywords = ref<string>('');
  const expandedKeys = ref<number[]>([]);
  const autoExpandParent = ref<boolean>(false);
  const loading = ref<boolean>(false);

  const treeData = ref([]);
  const treeDataMap = ref({} as any);

  const selectedKeys = ref([] as any[]);
  const selectedNode = ref(null as any);

  const activeTabKey = ref(0);
  const caseTabs = ref([] as any[]);
  const caseModel = ref({} as any);

  function selectNode(keys: any[], e: any) {
    window.console.log('selectNode', keys, e?.node?.dataRef);

    if (keys.length === 0 && e) {
      // un-select, keep the old one
      selectedKeys.value = [e.node.dataRef.id];
      return;
    } else {
      selectedKeys.value = keys;
    }

    setSelectedKeyCache(
      'case_tree',
      globalStore.currProject.id,
      selectedKeys.value[0],
    );

    selectedNode.value = treeDataMap.value[selectedKeys.value[0]];

    if (selectedNode.value.type === 'leaf') {
      openCaseTab(selectedNode.value.id);
    } else {
      // TODO: select a dir
    }
  }

  const selectStoredKey = debounce(async () => {
    window.console.log('selectStoredKey');
    const key = await getSelectedKeyCache(
      'case_tree',
      globalStore.currProject.id,
    );
    selectNode(key ? [key] : [], null);
  }, 300);

  function fetchTreeData() {
    if (!globalStore.currProject.id) return;

    loadCaseApi().then((result) => {
      window.console.log('loadCaseApi', result);
      treeData.value = result;

      getExpandedKeysCache('case', globalStore.currProject.id).then((keys) => {
        expandedKeys.value = keys;
      });

      treeDataMap.value = genNodeMap(treeData.value[0]);
    });
  }

  function onExpand(keys: any[], args: any) {
    window.console.log('onExpand', keys, args);
    expandedKeys.value = keys;
    autoExpandParent.value = false;

    setExpandedKeysCache(
      'case',
      globalStore.currProject.id,
      expandedKeys.value,
    );
  }
  function expandAll() {
    const keys: any = [];
    const data = treeData.value;

    function fn(arr: any) {
      if (!Array.isArray(arr)) {
        return;
      }
      arr.forEach((item) => {
        keys.push(item.id);
        if (Array.isArray(item.children)) {
          fn(item.children);
        }
      });
    }
    fn(data);
    expandedKeys.value = keys;
  }

  function createNode(parentId: number, type: string) {
    window.console.log('createNode', parentId, type);
    selectedNode.value = { parentId, type };
  }
  function editNode(node: any) {
    window.console.log('editNode', node.data);
    selectedNode.value = node.data;
  }
  async function deleteNode(node: any) {
    window.console.log('deleteNode', node.data);

    const title =
      node.type === 'dir' ? '将级联删除目录下的所有子目录、快捷调试' : '';
    const context = '删除后无法恢复，请确认是否删除？';

    confirmToDelete(title, context, () => {
      // TODO: delete
    });
  }

  async function onDrop(info: AntTreeNodeDropEvent) {
    if (info.node?.dataRef?.type === 'interface') {
      message.error('仅可移动到目录下');
      return;
    }
    const dropKey = info.node.eventKey;
    const dragKey = info.dragNode.eventKey;

    if (!info.node.pos) return;

    const pos = info.node.pos.split('-');
    const dropPosition = info.dropPosition - Number(pos[pos.length - 1]);

    const res = moveNode({
      dragKey,
      dropKey,
      dropPos: dropPosition, // 0 表示移动到目标节点的子节点，-1 表示移动到目标节点的前面， 1表示移动到目标节点的后面
    });
    if (res) {
      // 移动到目标节点的子节点，则需要展开目标节点
      if (
        dropKey &&
        dropPosition === 0 &&
        !isInArray(dropKey, expandedKeys.value)
      ) {
        expandedKeys.value.push(+dropKey);
      }
      message.success('移动成功');
    } else {
      message.warn('移动失败');
    }
  }

  function moveNode(data: any) {
    window.console.log(moveNode, data);
    return true;
  }

  function openCaseTab(id: number) {
    window.console.log('openCaseTab', id);
    if (id <= 0) {
      caseModel.value = null;
      return;
    }

    const found = caseTabs.value.find((item) => {
      return item.id === id;
    });

    getCaseApi(id).then((result) => {
      window.console.log(result);
      caseModel.value = result;

      if (!found) {
        caseTabs.value.push(result);
      }
    });
  }

  function changeTab(key: any) {
    window.console.log('changeTab', key);
    activeTabKey.value = key;

    openCaseTab(activeTabKey.value);
  }
  function onTabMenuClicked(key: any, action: any): void {
    window.console.log('onTabMenuClicked', key, action);
    if (action === 'remove') {
      removeCaseTab(key);
    }
  }

  async function removeCaseTab(id: number) {
    window.console.log('removeInterfaceTab', id);

    const needReload = id === caseModel.value.id;

    caseTabs.value = caseTabs.value.filter((tab: any) => tab.id !== id);
    window.console.log('after remove', caseTabs.value);

    let newOpenTabId = 0;
    // goto first one, if close curr tab
    if (caseTabs.value.length > 0 && caseModel.value.id === id) {
      newOpenTabId = caseTabs.value[0].id;
    }

    if (needReload) {
      openCaseTab(newOpenTabId);
    }
  }
  async function removeCaseTabs(id: number) {
    const removeTabIds = [] as number[];
    treeDataMap.value[id].children?.forEach((item: any) => {
      removeTabIds.push(item.id);
    });

    const needReload = removeTabIds.includes(caseModel.value.id);

    caseTabs.value = caseTabs.value.filter(
      (tab: any) => !removeTabIds.includes(tab.id),
    );

    let newOpenTabId = 0;
    // goto first one, if close curr tab
    if (
      caseTabs.value.length > 0 &&
      removeTabIds.includes(caseModel.value.id)
    ) {
      newOpenTabId = caseTabs.value[0].id;
    }

    if (needReload) {
      openCaseTab(newOpenTabId);
    }
  }

  function setCaseTabs(val: any) {
    caseTabs.value = val;
  }
  function setCaseModel(val: any) {
    caseModel.value = val;
  }

  watch(
    () => globalStore.currProject.id,
    (val: any) => {
      window.console.log('watch currProject in CaseStore', val);
      if (!val) return;

      selectStoredKey();

      loading.value = true;
      keywords.value = '';

      fetchTreeData();

      setTimeout(() => {
        loading.value = false;
      }, 100);
    },
  );
  watch(
    () => keywords,
    (val: any) => {
      expandedKeys.value = filterTreeNodes(treeData.value, val);
      autoExpandParent.value = true;
    },
  );

  function $reset() {
    keywords.value = '';
    expandedKeys.value = [];
    autoExpandParent.value = false;
    loading.value = false;

    treeData.value = [];
    treeDataMap.value = {};

    selectedKeys.value = [];
    selectedNode.value = null;

    activeTabKey.value = 0;
    caseTabs.value = [];
    caseModel.value = {};
  }

  return {
    keywords,
    expandedKeys,
    autoExpandParent,
    loading,
    treeData,

    caseTabs,
    caseModel,
    activeTabKey,

    fetchTreeData,
    onExpand,
    expandAll,

    selectedKeys,
    selectNode,
    createNode,
    editNode,
    deleteNode,
    changeTab,
    onTabMenuClicked,
    removeCaseTab,
    removeCaseTabs,
    openCaseTab,

    setCaseTabs,
    setCaseModel,

    onDrop,

    $reset,
  };
});

// 解决热更新问题
const hot = import.meta.hot;
if (hot) {
  hot.accept(acceptHMRUpdate(useCaseStore, hot));
}
