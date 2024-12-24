<script lang="ts" setup>
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { computed, ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';

import { Button } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { useProjectStore } from '#/views/project/store';

import EditModalComp from './edit.vue';

const projectStore = useProjectStore();
const projectState = computed(() => projectStore.projectState);

interface RowType {
  category: string;
  color: string;
  id: string;
  price: string;
  productName: string;
  releaseDate: string;
}

const formOptions: VbenFormProps = {
  // 默认展开
  collapsed: false,
  schema: [
    {
      component: 'Input',
      defaultValue: '',
      fieldName: 'keywords',
      label: '关键字',
    },
    {
      component: 'Select',
      componentProps: {
        allowClear: true,
        options: [
          {
            label: '启用',
            value: 'true',
          },
          {
            label: '禁用',
            value: 'false',
          },
        ],
        placeholder: '请选择',
      },
      fieldName: 'enabled',
      label: '状态',
    },
  ],
  showCollapseButton: true,
  submitOnChange: true,
  submitOnEnter: false,
};

const gridOptions: VxeGridProps<RowType> = {
  checkboxConfig: {
    highlight: true,
    labelField: 'name',
  },
  columns: [
    { title: '序号', type: 'seq', width: 50 },
    { title: '名称', field: 'name', type: 'checkbox', align: 'left' },
    { title: '修改人', field: 'updatedUser' },
    { title: '修改时间', field: 'updatedAt', formatter: 'formatDateTime' },
    {
      title: '操作',
      field: 'action',
      fixed: 'right',
      slots: { default: 'action' },
      width: 120,
    },
  ],
  height: 'auto',
  pagerConfig: {
    enabled: true,
    pageSize: 3,
    // currentPage: 0,
  },
  sortConfig: {
    multiple: true,
  },
  proxyConfig: {
    ajax: {
      query: async ({ page }, formValues) => {
        await projectStore.queryProjects({
          page: page.currentPage,
          pageSize: page.pageSize,
          ...formValues,
        });
        return projectStore.queryResult;
      },
    },
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

const [EditModal, editModalApi] = useVbenModal({
  connectedComponent: EditModalComp,
});

const model = ref<any>(null as any);
function edit(item: any) {
  window.console.log(item);
  model.value = item;
  editModalApi.open();
}

function finish() {
  window.console.log('finish');
  gridApi.query();
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #action="{ row }">
        <Button type="link" @click="edit(row)">编辑</Button>
      </template>
    </Grid>

    <EditModal :data="model" @finish="finish" />
    {{ projectState }}
  </Page>
</template>
