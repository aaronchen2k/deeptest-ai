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
    },
    name: 'Case',
    path: '/case',
    redirect: '/case/index',
    children: [
      {
        name: 'TestCase',
        path: '/case/index',
        component: () => import('#/views/_core/about/index.vue'),
        meta: {
          icon: 'lucide:copyright',
          title: $t('case.case'),
        },
      },
      {
        name: 'TestSuite',
        path: '/case/suite',
        component: () => import('#/views/_core/about/index.vue'),
        meta: {
          icon: 'lucide:copyright',
          title: $t('case.suite'),
        },
      },
    ],
  },
];

export default routes;
