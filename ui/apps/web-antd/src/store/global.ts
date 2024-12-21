import { acceptHMRUpdate, defineStore } from 'pinia';

import { loadProjectsApi, updateUserProject } from '#/api';

export interface GlobalState {
  currProject: any;
  projects: any[];
}

export const useGlobalStore = defineStore('global', {
  state: (): GlobalState => ({
    currProject: {},
    projects: [],
  }),
  // getters: {
  //   fullName(state): string {
  //     return `${state.name} ${state.email.split('@')[0]}`;
  //   },
  // },
  actions: {
    loadProjectData() {
      loadProjectsApi().then((result) => {
        window.console.log('loadProjectsApi', result);
        this.setCurrProject(result.default);
        this.setProjects(result.items);
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
