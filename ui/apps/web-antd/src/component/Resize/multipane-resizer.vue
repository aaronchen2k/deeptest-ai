<script>
export default {
  name: 'multipane-resizer',
  props: {
    collapsed: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<template>
  <div :class="{'collapsed': collapsed}" class="multipane-resizer">
    <slot></slot>
  </div>
</template>
<style lang="less" scoped>
.dp-multipane-con.layout-v .multipane-resizer {
  position: relative;
  z-index: 2;
  width: 8px;
  display: flex;
  align-items: center;
  cursor: col-resize;
  &:after {
    width: 1px;
    content: '';
    height: 100%;
    position: absolute;
    background-color: #f0f0f0;
    left: 4px;
  }

  &:hover:after {
    cursor: col-resize;
    width: 2px;
    background-color: #1890ff;
  }

  &.collapsed {
    &:after {
      display: none;
    }
  }
}

.dp-multipane-con.layout-h .multipane-resizer {
  position: relative;
  z-index: 2;
  height: 8px;
  width: 100%;
  display: flex;
  align-items: center;
  cursor: row-resize;
  &:after {
    height: 2px;
    content: '';
    width: 100%;
    position: absolute;
    background-color: #f0f0f0;
    left: 0;
    top: 4px;
  }

  &:hover:after {
    cursor: row-resize;
    height: 2px;
    width: 100%;
    background-color: #1890ff;
  }

  &.collapsed {
    &:after {
      display: none;
    }
  }
}
</style>
