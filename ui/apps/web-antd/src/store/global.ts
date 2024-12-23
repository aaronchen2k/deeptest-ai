import { acceptHMRUpdate, defineStore } from 'pinia';

import { listMyProjectApi, getCurrProjectApi, updateUserProject } from '#/api';

export interface GlobalState {
  currProject: any;
  projects: any[];
}

export const useGlobalStore = defineStore('global', {
  state: (): GlobalState => ({
    currProject: {},
    projects: [],
  }),
  actions: {
    listMyProject() {
      listMyProjectApi().then((result) => {
        window.console.log('listMyProjectApi', result);
        this.setProjects(result);
      });
    },
    getCurrProject() {
      getCurrProjectApi().then((result) => {
        window.console.log('getCurrProjectApi', result);
        this.setCurrProject(result);
      });
    },

    updateUserProject(item: any) {
      updateUserProject(item.id).then(() => {
        this.setCurrProject(item);
      });
    },

    setCurrProject(val: any) {
      this.currProject = val;
    },
    setProjects(val: any[]) {
      this.projects = val;
    },
  },
});

// 解决热更新问题
const hot = import.meta.hot;
if (hot) {
  hot.accept(acceptHMRUpdate(useGlobalStore, hot));
}
