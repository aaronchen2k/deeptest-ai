<script lang="ts" setup>
import { ref, watch } from 'vue';

import {Page, useVbenModal} from '@vben/common-ui';
import { useVbenForm, z } from '#/adapter/form';
import {getProjectApi, saveProjectApi} from '#/api';

interface Props {
  data: any;
}
const props = withDefaults(defineProps<Props>(), {
  data: {},
});

const model = ref(null as any);

const [Modal, modalApi] = useVbenModal({
  onCancel() {
    window.console.log('onCancel');
    modalApi.close();
  },
  onConfirm() {
    window.console.log('onConfirm');
    modalApi.close();
  },

  onOpenChange(isOpen) {
    if (isOpen) {
      window.console.log(props.data);
      // handleUpdate(10);
    }
  },
  footer: false,
});

const [Form, formApi] = useVbenForm({
  // 所有表单项共用，可单独在表单内覆盖
  commonConfig: {
    // 所有表单项
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
      getProjectApi(props.data.id).then((result) => {
        window.console.log(result);
        model.value = result;

        formApi.setValues(model.value);
      });
    }
  },
  { immediate: true, deep: true },
);

const emits = defineEmits<{
  finish: [event: any];
}>();
function submit(values: Record<string, any>) {
  window.console.log('submit', values, model.value);

  const data = Object.assign({}, model.value, values);

  saveProjectApi(data).then((result) => {
    emits('finish', null);
    modalApi.close();
  });
}

// function handleUpdate(len: number) {
//   modalApi.setState({ confirmDisabled: true, loading: true });
//   setTimeout(() => {
//     list.value = Array.from({ length: len }, (_v, k) => k + 1);
//     modalApi.setState({ confirmDisabled: false, loading: false });
//   }, 2000);
// }
</script>

<template>
  <Modal title="编辑项目">
    <Form />
  </Modal>
</template>
