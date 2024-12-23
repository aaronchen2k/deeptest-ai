import { ref, watch } from 'vue';

import { acceptHMRUpdate, defineStore } from 'pinia';

import { loadCaseApi } from '#/api/test/case';
import { useGlobalStore } from '#/store/global';

export const useCaseStore = defineStore('case', () => {
  const globalStore = useGlobalStore();

  const treeData = ref([]);

  async function fetchTreeData() {
    await loadCaseApi().then((result) => {
      window.console.log('loadCaseApi', result);
      treeData.value = result;
    });
  }

  watch(
    () => globalStore.currProject,
    (val: any) => {
      window.console.log('watch currProject in CaseStore', val.id);
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
