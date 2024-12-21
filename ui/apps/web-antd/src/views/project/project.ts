import { acceptHMRUpdate, defineStore } from 'pinia';

import { listProjectApi } from '#/api';

export interface ProjectState {
  queryResult: any[];
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
});

// 解决热更新问题
const hot = import.meta.hot;
if (hot) {
  hot.accept(acceptHMRUpdate(useProjectStore, hot));
}
