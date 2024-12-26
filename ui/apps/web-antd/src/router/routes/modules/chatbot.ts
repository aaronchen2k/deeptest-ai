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
      title: $t('chatbot.title'),
      hideChildrenInMenu: true,
    },
    name: 'Chatbot',
    path: '/chatbot',
    redirect: '/chatbot/index',
    children: [
      {
        name: 'ChatbotIndex',
        path: '/chatbot/index',
        component: () => import('#/views/chatbot/index.vue'),
        meta: {
          title: $t('chatbot.title'),
          icon: 'lucide:copyright',
        },
      },
    ],
  },
];

export default routes;
