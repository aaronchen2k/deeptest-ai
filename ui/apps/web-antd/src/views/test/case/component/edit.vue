<script lang="ts" setup>
import { ref, watch } from 'vue';

import { Page } from '@vben/common-ui';

import { useVbenForm } from '#/adapter/form';

interface Props {
  data: any;
}
const props = withDefaults(defineProps<Props>(), {
  data: {},
});

const model = ref(null as any);

const [Form, formApi] = useVbenForm({
  commonConfig: {
    componentProps: {
      class: 'w-full',
    },
  },
  // 提交函数
  handleSubmit: submit,
  layout: 'horizontal',
  schema: [
    {
      component: 'Input',
      componentProps: {
        placeholder: '',
      },
      fieldName: 'name',
      label: '名称',
      rules: 'required',
    },
    {
      component: 'Textarea',
      componentProps: {
        placeholder: '',
        autoSize: { minRows: 3, maxRows: 8 },
      },
      defaultValue: '',
      fieldName: 'desc',
      label: '描述',
    },
  ],
  wrapperClass: 'grid-cols-1',
});

watch(
  [() => props.data],
  async ([val]) => {
    window.console.log('watch', val);
    if (val?.id) {
      // getProjectApi(props.data.id).then((result) => {
      //   window.console.log(result);
      //   model.value = result;
      //
      //   formApi.setValues(model.value);
      // });
    }
  },
  { immediate: true, deep: true },
);

function submit(values: Record<string, any>) {
  window.console.log('submit', values, model.value);

  const data = Object.assign({}, model.value, values);
}
</script>

<template>
  <Page auto-content-height class="case-main">
    <span class="icon-[ant-design--folder-outlined]"></span>
    <span class="icon-[ant-design--caret-down-outlined]"></span>
  </Page>
</template>
