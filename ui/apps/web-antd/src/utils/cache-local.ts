import { toRaw } from '@vue/reactivity';
import localforage from 'localforage';

export const getCache = async (key: string): Promise<any | null> => {
  try {
    const ret = await localforage.getItem(key);
    return ret;
  } catch (error) {
    window.console.log('getCache err', error);
    return false;
  }
};

export const setCache = async (key: string, val: any): Promise<boolean> => {
  try {
    await localforage.setItem(key, toRaw(val));
    return true;
  } catch (error) {
    window.console.log('setCache err', error);
    return false;
  }
};
