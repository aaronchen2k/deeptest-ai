import { createVNode } from 'vue';

import { Modal } from 'ant-design-vue';

export function confirmToDo(
  title: string,
  content: string,
  callback: () => void,
  confirmText?: string,
  cancelText?: string,
) {
  Modal.confirm({
    okType: 'danger',
    title,
    icon: createVNode('span', {
      class: 'icon-[ant-design--exclamation-circle-outlined]',
    }),
    content,
    okText: () => confirmText || '确定',
    cancelText: () => cancelText || '取消',
    onOk: async () => {
      if (callback) await callback();
    },
    onCancel() {
      window.console.log('Cancel');
    },
  });
}

export function confirmToDelete(
  title: string,
  content: string,
  callback: () => void,
  confirmText?: string,
  cancelText?: string,
) {
  Modal.confirm({
    okType: 'danger',
    title,
    icon: createVNode('span', {
      class: 'icon-[ant-design--exclamation-circle-outlined]',
    }),
    content,
    okText: () => confirmText || '确定',
    cancelText: () => cancelText || '取消',
    onOk: async () => {
      if (callback) callback();
    },
    onCancel() {
      window.console.log('Cancel');
    },
  });
}

export function confirmToSave(
  callback: () => void,
  title?: string,
  content?: string,
  confirmText?: string,
  cancelText?: string,
) {
  Modal.confirm({
    okType: 'danger',
    title: title || '有修改内容未保存',
    content: content || '是否放弃未保存的修改？',
    okText: () => confirmText || '确定',
    cancelText: () => cancelText || '取消',
    onOk: async () => {
      if (callback) callback();
    },
    onCancel() {
      window.console.log('Cancel');
    },
  });
}
