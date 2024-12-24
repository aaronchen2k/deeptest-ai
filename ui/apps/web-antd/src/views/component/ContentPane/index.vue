<!--
  适用于 左侧目录树 + 右边区域表格筛选，且左侧目录树可伸缩
-->
<script setup lang="ts">
import { defineProps, ref } from 'vue';

import { SvgCollapseIcon, SvgExpandIcon } from '@vben/icons';

import { Multipane, MultipaneResizer } from '#/component/Resize/index';
import settings from '#/config/settings';

defineProps(['containerStyle', 'showExpand']);
const isFold = ref(true);
const paneLeft = ref();

const toggle = async () => {
  isFold.value = !isFold.value;
};

const handlePaneResize = () => {
  isFold.value = true;
  bus.emit(settings.paneResizeTop);
};
</script>

<template>
  <div :style="containerStyle || {}" class="container h-full">
    <div :class="[showExpand && 'expand-content']" class="content h-full">
      <Multipane
        class="vertical-panes h-full w-full"
        layout="vertical"
        @pane-resize="handlePaneResize"
      >
        <div ref="paneLeft" :class="[!isFold && 'unfold']" class="pane left">
          <slot name="left"></slot>
        </div>

        <MultipaneResizer />

        <div :class="[!isFold && 'unfold']" class="pane right">
          <slot name="right"></slot>
          <div v-if="showExpand" class="expand-icon" @click="toggle">
            <SvgExpandIcon v-if="!isFold" class="expand-icon-svg" />
            <SvgCollapseIcon v-else class="expand-icon-svg" />
          </div>
        </div>
      </Multipane>
    </div>
  </div>
</template>

<style lang="less" scoped>
.container {
  :deep(.ant-pagination) {
    margin-right: 24px;
  }
  .content {
    display: flex;
    width: 100%;
    position: relative;

    &.expand-content {
      .right {
        overflow: unset;
      }
    }

    .left {
      position: relative;
      min-width: 150px;
      width: 250px;
      max-width: 600px;

      &.unfold {
        width: 0 !important;
        min-width: 0 !important;
        overflow: hidden;
      }
    }

    .right {
      flex: 1;
      overflow: scroll;
      position: relative;
      z-index: 2;
      padding-left: 1px;

      &.unfold {
        overflow: unset;
      }

      &:has(.expand-icon:hover) {
        overflow: unset;
        z-index: 2;
      }

      .expand-icon {
        position: absolute;
        top: 50%;
        transform: translateY(-50%);
        left: -3px;
        cursor: pointer;

        .expand-icon-svg {
          fill: #f3f3f3;
          transition: fill 0.3s ease;
          box-shadow: 1px 1px 4px #ccc;
        }
      }
    }
  }
}
</style>
