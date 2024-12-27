<script setup lang="ts">
import { ref } from 'vue';

import { useAppConfig } from '@vben/hooks';

import {
  Button,
  message,
  Upload,
  type UploadChangeParam,
  type UploadProps,
} from 'ant-design-vue';

const { apiURL } = useAppConfig(import.meta.env, import.meta.env.PROD);

const uploadUrl = `${apiURL}/file/upload`;

const handleChange = (info: UploadChangeParam) => {
  if (info.file.status !== 'uploading') {
    window.console.log(info.file, info.fileList);
  }
  if (info.file.status === 'done') {
    message.success(`${info.file.name} file uploaded successfully`);
  } else if (info.file.status === 'error') {
    message.error(`${info.file.name} file upload failed.`);
  }
};

const fileList = ref([]);
const progress: UploadProps['progress'] = {
  strokeColor: {
    '0%': '#108ee9',
    '100%': '#87d068',
  },
  strokeWidth: 3,
  format: (numb = 0) => {
    return `${Number.parseFloat(numb.toFixed(2))}%`;
  },
  class: 'test',
};
const headers = { authorization: 'authorization-text' };

</script>

<template>
  <div class="uploader-main">
    <Upload
      v-model:file-list="fileList"
      :action="uploadUrl"
      :headers="headers"
      :progress="progress"
      name="file"
      @change="handleChange"
    >
      <Button class="upload-btn">
        <span
          class="dp-dropdown-icon icon-[ant-design--upload-outlined]"
        ></span>
        <!-- <span class="upload-text">Upload</span>-->
      </Button>
    </Upload>
  </div>
</template>

<style lang="less" scoped>
.uploader-main {
  width: 44px;

  .upload-btn {
    background-color: transparent;
    color: rgba(0, 0, 0, 0.88);
    border: 0;
  }

  .upload-text {
    display: inline-block;
    padding-left: 5px;
  }
}
</style>
