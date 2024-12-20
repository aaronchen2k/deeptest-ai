export function defineGlobalFuncs(app: any) {
  app.config.globalProperties.$info = $info;
  app.config.globalProperties.$warn = $warn;
  app.config.globalProperties.$error = $error;
}

export function $info(msg: string) {
  window.console.log(msg);
}
export function $warn(msg: string) {
  window.console.warn(msg);
}
export function $error(msg: string) {
  window.console.error(msg);
}
