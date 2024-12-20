<script lang="ts" setup>
import type { VbenFormProps } from '#/adapter/form';
import type { VxeGridProps } from '#/adapter/vxe-table';

import { ref } from 'vue';

import { Page, useVbenModal } from '@vben/common-ui';

import { Button } from 'ant-design-vue';

import { useVbenVxeGrid } from '#/adapter/vxe-table';
import { listProjectApi } from '#/api';

import EditModalComp from './edit.vue';

interface RowType {
  category: string;
  color: string;
  id: string;
  price: string;
  productName: string;
  releaseDate: string;
}

const currPage = ref(0);

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
    { title: '', align: 'left', type: 'checkbox', width: 100 },
    { title: '名称', field: 'name' },
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
        currPage.value = page.currentPage;
        return await query({
          page: page.currentPage,
          pageSize: page.pageSize,
          ...formValues,
        });
      },
    },
  },
};

const [Grid, gridApi] = useVbenVxeGrid({ formOptions, gridOptions });

async function query(data: any) {
  window.console.log(
    '=====',
    currPage.value,
    gridApi.grid.pagerConfig?.pageSize,
  );
  return listProjectApi(data);
}

const [EditModal, editModalApi] = useVbenModal({
  connectedComponent: EditModalComp,
});

const model = ref<any>(null as any);
function edit(item: any) {
  window.console.log(item);
  model.value = item;
  editModalApi.open();
}
</script>

<template>
  <Page auto-content-height>
    <Grid>
      <template #action="{ row }">
        <Button type="link" @click="edit(row)">编辑</Button>
      </template>
    </Grid>

    <EditModal :data="model" />
  </Page>
</template>
