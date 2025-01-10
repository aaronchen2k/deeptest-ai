<script lang="ts" setup>
import type { UnwrapRef } from 'vue';
import { defineEmits, defineProps, reactive, ref, toRaw, watch } from 'vue';

import { Form, FormItem, Input, Modal } from 'ant-design-vue';

const props = defineProps({
  nodeInfo: {
    required: true,
    type: Object,
  },
});

const emit = defineEmits(['ok', 'abandon']);

const formState = ref({
  id: 0,
  title: '',
  desc: '',
  type: '',
  parentId: 0,
});
watch(
  props.nodeInfo,
  () => {
    window.console.log('watch props.nodeInfo', props?.nodeInfo?.type);
    formState.value = {
      id: props?.nodeInfo?.id,
      title: props?.nodeInfo?.title,
      desc: props?.nodeInfo?.desc,
      type: props?.nodeInfo?.type,
      parentId: props?.nodeInfo?.parentId,
    };
  },
  { immediate: true, deep: true },
);

const rules = ref<any>({
  title: [{ required: true, message: '请输入名称', trigger: 'change' }],
  desc: [{ required: false }],
});

const useForm = Form.useForm;
const { resetFields, validate, validateInfos } = useForm(formState, rules, {
  onValidate: (...args) => window.console.log(...args),
});

function ok() {
  window.console.log('ok', toRaw(formState));

  validate()
    .then(() => {
      emit('ok', {
        ...formState.value,
      });
      resetFields();
    })
    .catch((error: any) => {
      window.console.log('error', error);
    });
}

function abandon() {
  emit('abandon');
  resetFields();
}
</script>

<template>
  <Modal
    :open="!!nodeInfo"
    :title="
      (!nodeInfo.id ? '新建' : '编辑') +
      (formState.type === 'leaf' ? '用例' : '目录')
    "
    width="600px"
    @cancel="abandon"
    @ok="ok"
  >
    <Form :wrapper-col="{ span: 14 }" class="custom-center-form">
      <FormItem
        :label="`${formState.type === 'leaf' ? '用例' : '目录'}名称`"
        v-bind="validateInfos.title"
      >
        <Input v-model:value="formState.title" placeholder="请输入名称" />
      </FormItem>

      <!-- <FormItem :label="(formState.type === 'leaf' ? '接口' : '目录') + '备注'" name="desc">
        <Input placeholder="请输入备注" v-model:value="formState.desc"/>
      </FormItem> -->
    </Form>
  </Modal>
</template>

<style lang="less" scoped>
.modal-btns {
  display: flex;
  justify-content: flex-end;
}
</style>
