import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      icon: 'lucide:layout-dashboard',
      order: -1,
      title: $t('home.title'),
    },
    name: 'Dashboard',
    path: '/home',
    children: [
      {
        name: 'Index',
        path: 'index',
        component: () => import('#/views/dashboard/analytics/index.vue'),
        meta: {
          affixTab: true,
          icon: 'lucide:area-chart',
          title: $t('home.title'),
          hideInMenu: true,
        },
      },
    ],
  },
];

export default routes;
