import { acceptHMRUpdate, defineStore } from 'pinia';

interface GlobalState {
  currProject: any;
  projects: any[];
}

export const useGlobalStore = defineStore('global', {
  actions: {
    setCurrProject(val: any) {
      this.currProject = val;
    },
    setProjects(val: any[]) {
      this.projects = val;
    },
  },
  state: (): GlobalState => ({
    currProject: {},
    projects: [],
  }),
});

// 解决热更新问题
const hot = import.meta.hot;
if (hot) {
  hot.accept(acceptHMRUpdate(useGlobalStore, hot));
}
