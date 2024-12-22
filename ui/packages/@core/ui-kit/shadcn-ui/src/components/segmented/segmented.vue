<script setup lang="ts">
import type { SegmentedItem } from './types';

import { computed } from 'vue';

import { Tabs, TabsContent } from '../../ui';

interface Props {
  defaultValue?: string;
  tabs: SegmentedItem[];
}

const props = withDefaults(defineProps<Props>(), {
  defaultValue: '',
  tabs: () => [],
});

const activeTab = defineModel<string>();

const getDefaultValue = computed(() => {
  return props.defaultValue || props.tabs[0]?.value;
});
</script>

<template>
  <Tabs v-model="activeTab" :default-value="getDefaultValue">
    <template v-for="tab in tabs" :key="tab.value">
      <TabsContent :value="tab.value">
        <slot :name="tab.value"></slot>
      </TabsContent>
    </template>
  </Tabs>
</template>
