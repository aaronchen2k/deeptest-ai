import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      icon: 'ic:baseline-view-in-ar',
      keepAlive: true,
      order: 2000,
      title: $t('set.title'),
      hideChildrenInMenu: true,
    },
    name: 'Set',
    path: '/set',
    redirect: '/set/index',
    children: [
      {
        name: 'SetIndex',
        path: 'index',
        component: () => import('#/views/empty/index.vue'),
        meta: {
          title: $t('set.title'),
          icon: 'lucide:copyright',
        },
      },
    ],
  },
];

export default routes;
