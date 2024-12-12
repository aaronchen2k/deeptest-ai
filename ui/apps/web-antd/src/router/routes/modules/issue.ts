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
      title: $t('issue.title'),
      hideChildrenInMenu: true,
    },
    name: 'Demos',
    path: '/issue',
    redirect: '/issue/index',
    children: [
      {
        name: 'IssueIndex',
        path: 'index',
        component: () => import('#/views/empty/index.vue'),
        meta: {
          title: $t('issue.title'),
          icon: 'lucide:copyright',
        },
      },
    ],
  },
];

export default routes;
