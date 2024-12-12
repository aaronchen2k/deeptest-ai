import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      badgeType: 'dot',
      icon: 'ic:baseline-view-in-ar',
      order: 1000,
      title: $t('case.title'),
      hideChildrenInMenu: true,
    },
    name: 'Case',
    path: '/case',
    redirect: '/case/index',
    children: [
      {
        name: 'CaseIndex',
        path: '/case/index',
        component: () => import('#/views/empty/index.vue'),
        meta: {
          title: $t('case.title'),
          icon: 'lucide:copyright',
        },
      },
    ],
  },
];

export default routes;
