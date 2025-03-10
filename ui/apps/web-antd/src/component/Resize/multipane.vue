<script>
const LAYOUT_HORIZONTAL = 'horizontal';
const LAYOUT_VERTICAL = 'vertical';

export default {
  name: 'Multipane',
  props: {
    layout: {
      type: String,
      default: LAYOUT_VERTICAL,
    },
  },
  emits: ['paneResizeStart', 'paneResizeStop', 'paneResize'],

  data() {
    return {
      isResizing: false,
    };
  },

  computed: {
    classnames() {
      return [
        'multipane',
        `layout-${this.layout.slice(0, 1)}`,
        this.isResizing ? 'is-resizing' : '',
      ];
    },
    cursor() {
      if (this.isResizing)
        return this.layout === LAYOUT_VERTICAL ? 'col-resize' : 'row-resize';
      else return '';
    },
    userSelect() {
      return this.isResizing ? 'none' : '';
    },
  },

  methods: {
    onMouseDown({ target: resizer, pageX: initialPageX, pageY: initialPageY }) {
      if (
        resizer?.className &&
        resizer?.className?.match &&
        resizer?.className?.match('multipane-resizer')
      ) {
        // eslint-disable-next-line unicorn/no-this-assignment
        const self = this;
        const { $el: container, layout } = self;

        const topOrLeftPane = resizer.previousElementSibling;
        const {
          offsetWidth: initialPaneWidth,
          offsetHeight: initialPaneHeight,
        } = topOrLeftPane;

        const usePercentage = `${topOrLeftPane.style.width}`.match('%');

        const { addEventListener, removeEventListener } = window;

        const resize = (initialSize, offset = 0) => {
          // if(this.collapsed){
          //   return pane.style.width = '0px'
          // }

          if (layout === LAYOUT_VERTICAL) {
            const containerWidth = container.clientWidth;
            const paneWidth = initialSize + offset;

            return (topOrLeftPane.style.width = usePercentage
              ? `${(paneWidth / containerWidth) * 100}%`
              : `${paneWidth}px`);
          }

          if (layout === LAYOUT_HORIZONTAL) {
            const containerHeight = container.clientHeight;
            const paneHeight = initialSize + offset;

            return (topOrLeftPane.style.height = usePercentage
              ? `${(paneHeight / containerHeight) * 100}%`
              : `${paneHeight}px`);
          }
        };

        // This adds is-resizing class to container
        self.isResizing = true;

        // Resize once to get current computed size
        let size = resize(
          layout === LAYOUT_VERTICAL ? initialPaneWidth : initialPaneHeight,
          0,
        );

        // Trigger paneResizeStart event
        self.$emit('paneResizeStart', topOrLeftPane, resizer, size);

        const onMouseMove = function ({ pageX, pageY }) {
          window.console.log('onMouseMove');
          size =
            layout === LAYOUT_VERTICAL
              ? resize(initialPaneWidth, pageX - initialPageX)
              : resize(initialPaneHeight, pageY - initialPageY);
          topOrLeftPane.classList.add('left-drag');
          self.$emit('paneResize', topOrLeftPane, resizer, size);
        };

        const onMouseUp = function () {
          // Run resize one more time to set computed width/height.
          size =
            layout === LAYOUT_VERTICAL
              ? resize(topOrLeftPane.clientWidth)
              : resize(topOrLeftPane.clientHeight);

          // This removes is-resizing class to container
          self.isResizing = false;

          removeEventListener('mousemove', onMouseMove);
          removeEventListener('mouseup', onMouseUp);

          self.$emit('paneResizeStop', topOrLeftPane, resizer, size);
        };

        addEventListener('mousemove', onMouseMove);
        addEventListener('mouseup', onMouseUp);
      }
    },
  },
};
</script>

<template>
  <div
    :class="classnames"
    :style="{ cursor, userSelect }"
    class="dp-multipane-con"
    @mousedown.stop="onMouseDown"
  >
    <slot></slot>
  </div>
</template>

<style lang="less">
.multipane {
  display: flex;

  &.layout-h {
    flex-direction: column;
  }

  &.layout-v {
    flex-direction: row;
  }
}

.multipane > div {
  position: relative;
  z-index: 1;
}
</style>
