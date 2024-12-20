import { acceptHMRUpdate, defineStore } from 'pinia';

interface GlobalState {
  currProjectId: number;
}

export const useGlobalStore = defineStore('global', {
  actions: {
    setCurrProjectId(id: number) {
      this.currProjectId = id;
    },
  },
  state: (): GlobalState => ({
    currProjectId: -1,
  }),
});

// 解决热更新问题
const hot = import.meta.hot;
if (hot) {
  hot.accept(acceptHMRUpdate(useGlobalStore, hot));
}
