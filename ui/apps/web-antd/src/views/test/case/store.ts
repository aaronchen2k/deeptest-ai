import type { AntTreeNodeDropEvent } from 'ant-design-vue/lib/tree';

import { ref, watch } from 'vue';

import {
  filterTreeNodes,
  genNodeMap,
  isInArray,
} from '@vben/utils';

import { message } from 'ant-design-vue';
import { debounce } from 'lodash';
import { acceptHMRUpdate, defineStore } from 'pinia';

import { loadCaseApi } from '#/api/test/case';
import { useGlobalStore } from '#/store/global';
import {
  getExpandedKeysCache,
  getSelectedKeyCache,
  setExpandedKeysCache,
  setSelectedKeyCache
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

  function selectNode(keys: any[], e: any) {
    window.console.log('selectNode', keys, e?.node?.dataRef);
    if (!e?.node) {
      return;
    }

    selectedKeys.value = [e.node.dataRef.id];

    setSelectedKeyCache(
      'case_tree',
      globalStore.currProject.id,
      selectedKeys.value[0],
    );

    selectedNode.value = treeDataMap.value[selectedKeys.value[0]];
    openCaseTab(selectedNode.value);
  }

  const selectStoredKeyCall = debounce(async () => {
    window.console.log('selectStoredKeyCall');
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
    selectedNode.value = node;
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

    const res = await moveNode({
      dragKey, // 移动谁
      dropKey, // 移动那儿
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

  function openCaseTab(item: any) {
    window.console.log('openCaseTab', item);
  }

  watch(
    () => globalStore.currProject.id,
    (val: any) => {
      window.console.log('watch currProject in CaseStore', val);

      selectStoredKeyCall();

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

  return {
    keywords,
    expandedKeys,
    autoExpandParent,
    loading,
    treeData,

    fetchTreeData,
    onExpand,
    expandAll,

    selectedKeys,
    selectNode,
    createNode,
    editNode,
    deleteNode,

    onDrop,
  };
});

// 解决热更新问题
const hot = import.meta.hot;
if (hot) {
  hot.accept(acceptHMRUpdate(useCaseStore, hot));
}
