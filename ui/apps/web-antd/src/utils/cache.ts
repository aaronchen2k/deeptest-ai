import settings from '#/config/settings';

import { getCache, setCache } from './cache-local';

export const getExpandedKeysCache = async (type: string, id: number) => {
  window.console.log('getExpandedKeysCache');
  const key = `${type}-${id}`;

  const cachedData = await getCache(settings.expandedKeys);
  if (!cachedData || !cachedData[key]) {
    return [];
  }

  const keys = cachedData[key] ?? [];

  return [...keys];
};

export const setExpandedKeysCache = async (
  type: string,
  id: number,
  keys: string[],
) => {
  window.console.log('setExpandedKeysCache');
  if (!keys) keys = [];
  const key = `${type}-${id}`;

  let cachedData = await getCache(settings.expandedKeys);
  if (!cachedData) cachedData = {};

  const items = [] as any[];
  keys.forEach((item) => {
    items.push(item);
  });
  cachedData[key] = items;
  await setCache(settings.expandedKeys, cachedData);
};

// Tree Selected Key
export const getSelectedKeyCache = async (type: string, projectId: number) => {
  window.console.log('getSelectedKeyCache');
  const key = `${type}-${projectId}`;

  const cachedData = await getCache(settings.selectedKey);
  if (!cachedData || !cachedData[key]) {
    return null;
  }

  return cachedData[key];
};
export const setSelectedKeyCache = async (
  type: string,
  projectId: number,
  selectedKey: string,
) => {
  window.console.log('setSelectedKeyCache');
  const key = `${type}-${projectId}`;

  let cachedData = await getCache(settings.selectedKey);
  if (!cachedData) cachedData = {};

  cachedData[key] = selectedKey;
  await setCache(settings.selectedKey, cachedData);
};
