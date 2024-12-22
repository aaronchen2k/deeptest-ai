<script setup lang="ts">
import type { SupportedLanguagesType } from '@vben/locales';
import type { ThemeModeType } from '@vben/types';
import type { SegmentedItem } from '@vben-core/shadcn-ui';

import { computed, ref } from 'vue';

import { $t } from '@vben/locales';
import {
  clearPreferencesCache,
  resetPreferences,
  usePreferences,
} from '@vben/preferences';
import { useVbenDrawer } from '@vben-core/popup-ui';
import { VbenButton, VbenSegmented } from '@vben-core/shadcn-ui';

import { Block, General, Theme } from './blocks';

const emit = defineEmits<{ clearPreferencesAndLogout: [] }>();

const appLocale = defineModel<SupportedLanguagesType>('appLocale');
const appDynamicTitle = defineModel<boolean>('appDynamicTitle');
const appWatermark = defineModel<boolean>('appWatermark');
const appEnableCheckUpdates = defineModel<boolean>('appEnableCheckUpdates');

const themeMode = defineModel<ThemeModeType>('themeMode');
const themeSemiDarkSidebar = defineModel<boolean>('themeSemiDarkSidebar');
const themeSemiDarkHeader = defineModel<boolean>('themeSemiDarkHeader');

const { diffPreference } = usePreferences();

const [Drawer] = useVbenDrawer();

const activeTab = ref('appearance');

const tabs = computed((): SegmentedItem[] => {
  return [
    {
      label: $t('preferences.appearance'),
      value: 'appearance',
    },
  ];
});

async function handleClearCache() {
  window.console.log('handleClearCache');
  resetPreferences();
  clearPreferencesCache();
  emit('clearPreferencesAndLogout');
}
</script>

<template>
  <div>
    <Drawer
      :description="$t('preferences.subtitle')"
      :title="$t('preferences.title')"
      class="sm:max-w-sm"
    >
      <template #extra>
        <div class="flex items-center"></div>
      </template>

      <div class="p-1">
        <VbenSegmented v-model="activeTab" :tabs="tabs">
          <template #appearance>
            <Block :title="$t('preferences.theme.title')">
              <Theme
                v-model="themeMode"
                v-model:theme-semi-dark-header="themeSemiDarkHeader"
                v-model:theme-semi-dark-sidebar="themeSemiDarkSidebar"
              />
            </Block>
            <Block :title="$t('preferences.general')">
              <General
                v-model:app-dynamic-title="appDynamicTitle"
                v-model:app-enable-check-updates="appEnableCheckUpdates"
                v-model:app-locale="appLocale"
                v-model:app-watermark="appWatermark"
              />
            </Block>
          </template>
        </VbenSegmented>
      </div>

      <template #footer>
        <VbenButton
          :disabled="!diffPreference"

          class="mr-4 w-full"
          size="sm"
          variant="ghost"
          @click="handleClearCache"
        >
          <span v-if="diffPreference">
            {{ $t('preferences.clearAndLogout') }}
          </span>
          <span v-else>
            {{ $t('preferences.clearAndLogoutDisabled') }}
          </span>
        </VbenButton>
      </template>
    </Drawer>
  </div>
</template>
