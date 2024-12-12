import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      icon: 'ic:baseline-view-in-ar',
      keepAlive: true,
      order: 4000,
      title: $t('exec.title'),
      hideChildrenInMenu: true,
    },
    name: 'Execution',
    path: '/exec',
    redirect: '/exec/index',
    children: [
      {
        name: 'ExecutionIndex',
        path: 'index',
        component: () => import('#/views/empty/index.vue'),
        meta: {
          title: $t('exec.title'),
          icon: 'lucide:copyright',
        },
      },
    ],
  },
];

export default routes;
