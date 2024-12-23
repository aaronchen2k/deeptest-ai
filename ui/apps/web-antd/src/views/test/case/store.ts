import { ref, watch } from 'vue';

import { acceptHMRUpdate, defineStore } from 'pinia';

import { loadCaseApi } from '#/api/test/case';
import { useGlobalStore } from '#/store/global';

export const useCaseStore = defineStore('case', () => {
  const globalStore = useGlobalStore();

  const treeData = ref([]);

  function fetchTreeData() {
    if (!globalStore.currProject.id) return;

    loadCaseApi().then((result) => {
      window.console.log('loadCaseApi', result);
      treeData.value = result;
    });
  }

  watch(
    () => globalStore.currProject.id,
    (val: any) => {
      window.console.log('watch currProject in CaseStore', val);
      fetchTreeData();
    },
  );

  return {
    treeData,
    fetchTreeData,
  };
});

// 解决热更新问题
const hot = import.meta.hot;
if (hot) {
  hot.accept(acceptHMRUpdate(useCaseStore, hot));
}
