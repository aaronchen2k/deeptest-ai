import { acceptHMRUpdate, defineStore } from 'pinia';

import { listProjectApi } from '#/api';
import { useGlobalStore } from '#/store/global';

export interface ProjectState {
  queryResult: any;
}

export const useProjectStore = defineStore('project', {
  state: (): ProjectState => ({
    queryResult: [],
  }),
  actions: {
    async queryProjects(data: any) {
      await listProjectApi(data).then((result) => {
        window.console.log('listProjectApi', result);
        this.setQueryResult(result);
      });
    },
    setQueryResult(val: any[]) {
      this.queryResult = val;
    },
  },
  getters: {
    projectState(state): string {
      window.console.log('=== gets: projectState');
      const globalStore = useGlobalStore();
      const msg = `${state.queryResult.items?.length}, ${globalStore.currProject.name}, ${new Date()}`;
      return msg;
    },
  },
});

// 解决热更新问题
const hot = import.meta.hot;
if (hot) {
  hot.accept(acceptHMRUpdate(useProjectStore, hot));
}
