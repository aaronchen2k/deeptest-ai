import {
  computed,
  defineComponent,
  inject,
  type PropType,
  ref,
  toRefs,
} from 'vue';

import { Dropdown, Menu, MenuItem } from 'ant-design-vue';

import { type MenuItem as MenuItemType, type Recordable } from './type';

import './index.less';

const DropdownMenuProps = {
  dropdownList: {
    type: Array as PropType<MenuItemType[]>,
    default: [],
  }, // 下拉菜单
  actionList: {
    type: Array as PropType<MenuItemType[]>,
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
  item: MenuItemType;
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

  const renderLabel = () => {
    return typeof item.label === 'function' ? (
      item.label(record)
    ) : (
      <div class="label-desc">
        <div class="label">{item.label}</div>
        <div class="desc">{item.desc}</div>
      </div>
    );
  };

  return !item.children || item.children?.length === 0 ? (
    <MenuItem key={item.key} onClick={(e: any) => handleClick(e)}>
      <span class="drop-down-menu-text">{renderLabel()}</span>
    </MenuItem>
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

const DropdownList = defineComponent({
  name: 'DropdownList',
  props: {
    list: {
      type: Array as PropType<MenuItemType[]>,
      default: () => [],
    },
    record: {
      type: Object,
      default: () => {},
    },
    selectedKey: {
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
          <Menu
            onOpenChange={(e: any) => handleOpen(e)}
            openKeys={openKeys.value}
            selectedKeys={[props.selectedKey]}
            style={{ maxHeight: '300px', overflowY: 'auto' }}
          >
            {props.list.map((e: any) =>
              RenderMenuItem({ item: e, record: props.record }),
            )}
          </Menu>
        );
      },
    };

    return () => {
      return <Dropdown v-slots={vslots} />;
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
    const { dropdownList, record } = toRefs(props);

    const ifShow = (actionItem: MenuItemType, props: any) => {
      if (typeof actionItem.ifShow === 'boolean') {
        return actionItem.ifShow;
      }
      if (typeof actionItem.ifShow === 'function') {
        return actionItem.ifShow(props.record);
      }
      return true;
    };

    const filterAction = (e: any, props: any) => {
      return ifShow(e, props);
    };
    const filteredDropDownList = computed(() =>
      dropdownList.value.filter((e) => filterAction(e, props)),
    );

    return () => {
      return (
        <div class="drop-down-action-wrap">
          {filteredDropDownList.value.length > 0 &&
            filteredDropDownList.value.length > 0 && (
              <DropdownList
                list={filteredDropDownList.value}
                record={record.value}
                selectedKey={props.selectedKey}
                v-slots={slots}
              />
            )}
        </div>
      );
    };
  },
});
