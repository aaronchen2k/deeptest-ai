import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      icon: 'ic:baseline-view-in-ar',
      keepAlive: true,
      order: 5000,
      title: $t('admin.title'),
    },
    name: 'Administration',
    path: '/admin',
    redirect: '/home/index',
    children: [
      {
        name: 'ProjectManagement',
        path: 'project',
        component: () => import('#/views/empty/index.vue'),
        meta: {
          title: $t('admin.project'),
          icon: 'lucide:copyright',
        },
      },
    ],
  },
];

export default routes;
