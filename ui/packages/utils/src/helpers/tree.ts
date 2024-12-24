import { cloneDeep } from 'lodash';

import { isInArray } from './array';

export function getSelectedTreeNode(
  checkedKeys: any,
  treeDataMapValue: any,
): any[] {
  const childrenMap = {} as any; // nodes that is other's child
  checkedKeys.forEach((id: number) => {
    if (treeDataMapValue[id].children) {
      treeDataMapValue[id].children.forEach((child: any) => {
        getChildren(treeDataMapValue[child.id], childrenMap);
      });
    }
  });
  const selectedNodes = [] as any[];

  Object.keys(treeDataMapValue).forEach((id) => {
    if (!childrenMap[id] && isInArray(id, checkedKeys)) {
      // in array and except other's child
      const node = treeDataMapValue[id];
      if (!node.isDir || node.children) {
        selectedNodes.push(node);
      }
    }
  });

  return selectedNodes;
}

function getChildren(node: any, mp: any) {
  mp[node.id] = true;

  if (node.children) {
    node.children.forEach((child: any) => {
      getChildren(child, mp);
    });
  }
}

export function filterTreeNodes(
  treeDataValue: any,
  keywords: string,
): number[] {
  if (!treeDataValue || treeDataValue.length === 0) return [];

  const flattenTreeList = flattenTree(treeDataValue[0]);

  let parentKeys: any = [];
  for (let i = 0; i < flattenTreeList.length; i++) {
    const node = flattenTreeList[i] as any;

    const text = node.title ?? node.name;
    // 兼容大小写问题
    if (text.toLowerCase().includes(keywords.toLowerCase().trim())) {
      parentKeys.push(node.parentId);
      parentKeys = [
        ...parentKeys,
        ...findParentIds(node.parentId, flattenTreeList),
      ];
    }
  }
  parentKeys = [...new Set(parentKeys)];

  return parentKeys;
}

function flattenTree(tree: any) {
  const nodes: Array<any> = [];

  function traverse(node: any) {
    nodes.push(node);
    if (node.children) {
      node.children.forEach((element: any) => {
        traverse(element);
      });
    }
  }

  traverse(tree);

  return nodes;
}

export function findParentIds(nodeId: number, tree: any) {
  let current: any = tree.find((node: any) => node.id === nodeId);
  const parentIds: Array<any> = [];
  while (current && current.parentId) {
    parentIds.unshift(current.parentId); // unshift 方法可以将新元素添加到数组的开头
    current = tree.find((node: any) => node.id === current.parentId);
  }
  return parentIds;
}

/**
 * @desc 根据关键词过滤树节点
 * @param {Array} children 树节点
 * @param {String} keyword 关键词
 * @param {String} field 搜索字段
 * @return {Array} 过滤后的树节点
 * */
export function filterByKeyword(
  children: any[],
  keyword: string,
  field = 'title',
) {
  if (!keyword.trim()) return children;

  function filterChildren(node: any) {
    if (node?.children?.length) {
      node.children = node.children.filter((child: any) => {
        return filterChildren(child);
      });
    }
    return hasChildrenByKeyword(node, keyword, field);
  }

  return children.filter((menu) => {
    return filterChildren(menu);
  });
}

/**
 * @desc 该节点下是否包含关键词
 * @param {Object} node 节点
 * @param {String} keyword 关键词
 * @param {String} field 搜索字段
 * @return {Boolean} 是否包含关键词
 * */
function hasChildrenByKeyword(node: any, keyword: string, field = 'title') {
  let result = false;

  // 定义递归函数，用于遍历树节点
  function traverse(node: any) {
    if (node?.[field]?.toLowerCase()?.includes(keyword.toLowerCase().trim())) {
      result = true;
      return;
    }
    // 递归处理子节点
    if (node?.children?.length > 0) {
      for (const child of node.children) {
        traverse(child);
      }
    }
  }

  // 调用递归函数，开始遍历
  traverse(node);
  return result;
}

export function findPath(nodeId: number, nodes: any[]): number[] {
  for (const node of nodes) {
    if (node.id === nodeId) {
      return [node.id];
    }
    if (node.children) {
      const path = findPath(nodeId, node.children);
      if (path.length > 0) {
        return [node.id, ...path];
      }
    }
  }

  return [];
}

export const getAllTabsId = (data: any) => {
  let result: any[] = [];
  data.forEach((el: any) => {
    if (el.id !== 0) {
      result.push(el.id);
    } else if (el.id === 0 && el.children) {
      result = [...result, ...getAllTabsId(el.children)];
    }
  });
  return result;
};

/**
 * 树结构
 * @param {Array} data 树的结构
 * @param {String} key 当前节点id
 * @param {String} callback 回调函数
 * @param {String} defaultKey 默认节点
 * @returns Array
 */
export const loopTree = (
  data: any,
  currKey: number,
  callback: any,
  defaultKey: number,
) => {
  // 循环树节点
  data.forEach((item: any, index: number, arr: any[]) => {
    if (item[defaultKey] === currKey) {
      return callback(item, index, arr);
    }
    if (item.children) {
      return loopTree(item.children, currKey, callback, defaultKey);
    }
  });
  return [...data];
};

export const removeLeafNode = (data: any) => {
  const arrayData = cloneDeep(data);
  arrayData.forEach((e: any) => {
    e.children = (e.children || []).filter((e: any) => e.entityId === 0);
    if (e.children) {
      e.children = removeLeafNode(e.children);
    }
  });

  return [...arrayData];
};

/**
 * 将树结构转化为map结构
 * @param treeData
 * @returns
 */
export const transTreeNodesToMap = (treeData: any[]): any => {
  const nodesMap: any = {};
  treeData.forEach((node) => {
    if (!nodesMap[node.id]) {
      nodesMap[node.id] = node;
    }
    if (Array.isArray(node.children)) {
      const res = transTreeNodesToMap(node.children);
      Object.assign(nodesMap, res);
    }
  });
  return nodesMap;
};

export function genNodeMap(treeNode: any, ids?: number[]): any {
  const mp = {};
  getNodeMap(treeNode, mp, ids);

  return mp;
}

export function getNodeMap(treeNode: any, mp: any, ids?: number[]): void {
  if (!treeNode) return;

  treeNode.entity = null;
  mp[treeNode.id] = treeNode;
  if (ids && treeNode.entityCategory !== 'processor_group') {
    ids.push(treeNode.id);
    // console.log('===', treeNode.entityCategory)
  }

  if (treeNode.children) {
    treeNode.children.forEach((item: any) => {
      getNodeMap(item, mp, ids);
    });
  }
}

export function expandAllKeys(treeMap: any, isExpand: boolean): number[] {
  const keys = new Array<number>();
  if (!isExpand) return keys;

  Object.keys(treeMap).forEach((key: string) => {
    if (!keys.includes(+key)) keys.push(+key);
  });

  return keys;
}

export function expandOneKey(
  treeMap: any,
  key: number,
  expandedKeys: number[],
) {
  if (!expandedKeys.includes(key)) expandedKeys.push(key);

  if (treeMap[key]) {
    const parentId = treeMap[key].parentId;
    if (parentId) {
      expandOneKey(treeMap, parentId, expandedKeys);
    }
  }
}
