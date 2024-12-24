export type Recordable<T = any> = {
  [x: string]: T;
};

export type MenuItem = {
  action?: (...args: any[]) => void;
  /** 权限编码 */
  auth?: string;
  checkExecClickAble?: boolean;
  children?: MenuItem[];
  customLoadingRender?: any;
  customRender?: any;
  /** 描述 */
  desc: any;
  disabled?: boolean;
  /** 显示图标，只支持图片 */
  icon?: string;
  /** 是否渲染 */
  ifShow?: ((record: Recordable, action?: MenuItem) => boolean) | boolean;
  key?: number | string;
  /** 操作名称 */
  label?: any;
  loadingText?: string;
  renderChildren?: (record: Recordable) => any[];
  /** 判断是否展示按钮， 优先以这个为准 */
  show?: ((record: Recordable, action?: MenuItem) => boolean) | boolean; // 部分按钮展示条件特殊，这个作为一些备用项来判断按钮的展示与否
  tip?: string;
  /** 提示 */
  tooltip?: string;
  value?: number | string;
};
