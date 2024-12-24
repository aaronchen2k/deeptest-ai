import { defineComponent, inject, type PropType, ref, toRefs } from 'vue';

import { type MenuItem, type Recordable } from './type';

import './index.less';

const DropdownMenuProps = {
  dropdownList: {
    type: Array as PropType<MenuItem[]>,
    default: [],
  }, // 下拉菜单
  actionList: {
    type: Array as PropType<MenuItem[]>,
    default: [],
  }, // 无下拉的菜单
  record: {
    type: Object,
    default: {},
  }, // 当前操作项
  selectedKey: {
    type: [Number, String],
    default: '',
  },
};

const RenderMenuItem = ({
  item,
  record,
}: {
  item: MenuItem;
  record: Recordable;
}) => {
  const handleClick = (_e?: any) => {
    if (item.disabled) {
      _e.preventDefault();
      return;
    }
    if (_e?.target?.children[0]?.disabled) {
      _e.preventDefault();
      _e.stopPropagation();
      return;
    }
    item.action?.(record);
  };

  return !item.children || item.children?.length === 0 ? (
    <a-menu-item key={item.key} onClick={(e: any) => handleClick(e)}>
      <span class="drop-down-menu-text">{item.label}</span>
    </a-menu-item>
  ) : (
    <a-sub-menu
      class={{ 'dp-action-submenu': true }}
      key={item.key}
      title={item.label}
    >
      {item.children?.map((e: any) => RenderMenuItem({ item: e, record }))}
    </a-sub-menu>
  );
};

const ActionList = (opts: { list: MenuItem[]; record: Recordable }) => {
  const { list, record } = opts;
  const customRenderLoadingLabel = (item: any) => {
    if (typeof item.customLoadingRender === 'function') {
      return item.customLoadingRender(record);
    }
    return item.customLoadingRender;
  };
  const customRenderLabel = (item: any) => {
    if (typeof item.customRender === 'function') {
      return item.customRender(record);
    }
    return item.customRender;
  };
  return (
    <div class="action-list">
      {list.map((actionItem: MenuItem) => (
        <div class="action-item" onClick={() => actionItem.action?.(record)}>
          {actionItem.customLoadingRender ? (
            <a-tooltip
              placement="top"
              title={
                record.loading
                  ? actionItem.loadingText || null
                  : actionItem.label
              }
            >
              {customRenderLoadingLabel(actionItem)}
            </a-tooltip>
          ) : (
            actionItem.customRender ? customRenderLabel(actionItem) : actionItem.label
          )}
        </div>
      ))}
    </div>
  );
};

const DropdownList = defineComponent({
  name: 'DropdownList',
  props: {
    list: {
      type: Array as PropType<MenuItem[]>,
      default: () => [],
    },
    record: {
      type: Object,
      default: () => {},
    },
    selectedkey: {
      type: [Number, String],
      default: '',
    },
  },
  setup(props, { slots }) {
    const handleOpenChange = inject('handleOpenChange', null) as any;
    const openKeys = ref([]);

    const handleOpen = (e: any) => {
      if (handleOpenChange) {
        handleOpenChange?.(e);
      }
    };

    const vslots = {
      default: () => {
        return (
          slots?.default?.() || (
            <span class="icon-[ant-design--more-outlined]"></span>
          )
        );
      },
      overlay: () => {
        return (
          <a-menu
            onOpenChange={(e: any) => handleOpen(e)}
            openKeys={openKeys.value}
            selectedKeys={[props.selectedkey]}
            style={{ maxHeight: '300px', overflowY: 'auto' }}
          >
            {props.list.map((e: any) =>
              RenderMenuItem({ item: e, record: props.record }),
            )}
          </a-menu>
        );
      },
    };

    return () => {
      return <a-dropdown v-slots={vslots} />;
    };
  },
});

/**
 * dropdownMenu组件
 */
export const DropdownActionMenu = defineComponent({
  name: 'DropdownMenu',
  props: DropdownMenuProps,
  setup(props, { slots }) {
    const { dropdownList, actionList, record } = toRefs(props);

    return () => {
      return (
        <div class="drop-down-action-wrap">
          {actionList.value.length > 0 && (
            <ActionList list={actionList.value} record={record.value} />
          )}
          {actionList.value.length > 0 && dropdownList.value.length > 0 && (
            <a-divider type="vertical" />
          )}
          {dropdownList.value.length > 0 && dropdownList.value.length > 0 && (
            <DropdownList
              list={dropdownList.value}
              record={record.value}
              selectedkey={props.selectedKey}
              v-slots={slots}
            />
          )}
        </div>
      );
    };
  },
});
