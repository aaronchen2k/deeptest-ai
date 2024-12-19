import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      icon: 'ic:baseline-view-in-ar',
      keepAlive: true,
      order: 3000,
      title: $t('plan.title'),
      hideInMenu: true,
    },
    name: 'Project',
    path: '/project',
    redirect: '/project/index',
    children: [
      {
        name: 'ProjectIndex',
        path: 'index',
        component: () => import('#/views/project/index.vue'),
        meta: {
          title: $t('project.title'),
          icon: 'lucide:copyright',
        },
      },
    ],
  },
];

export default routes;
