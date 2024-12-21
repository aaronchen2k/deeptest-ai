import type { RouteRecordRaw } from 'vue-router';

import { BasicLayout } from '#/layouts';
import { $t } from '#/locales';

const routes: RouteRecordRaw[] = [
  {
    component: BasicLayout,
    meta: {
      icon: 'lucide:settings',
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
        redirect: '/project/index',
        meta: {
          title: $t('admin.project'),
          icon: 'lucide:copyright',
        },
      },
    ],
  },
];

export default routes;
