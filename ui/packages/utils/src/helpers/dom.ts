export function hasClass(element: any, cName: any) {
  if (!element) return false;
  return new RegExp(`(\\s|^)${cName}(\\s|$)`).test(element.className);
}

export function addClass(elements: any, cName: any) {
  if (!elements) return;
  if (!hasClass(elements, cName)) {
    elements.className += ` ${cName}`;
  }
}

export function removeClass(elements: any, cName: any) {
  if (!elements) return;
  if (hasClass(elements, cName)) {
    elements.className = elements.className.replace(
      new RegExp(`(\\s|^)${cName}(\\s|$)`),
      ' ',
    );
  }
}

export function getContextMenuStyle(x: number, y: number, height: number) {
  let top = y + 6;
  if (y + height > document.body.clientHeight)
    top = document.body.clientHeight - height;

  return {
    zIndex: 99,
    position: 'fixed',
    maxHeight: 40,
    textAlign: 'center',
    left: `${x + 10}px`,
    top: `${top}px`,
  };
}

export const findParentNodeByX = (
  node: Element,
  opts: {
    class?: string;
    id?: string;
  } = {},
): Element | undefined => {
  const { class: className, id } = opts;
  if (className) {
    for (let i = 0; i < node.classList.length; i++) {
      const item = node.classList[i];
      if (className === item) {
        return node;
      }
    }
  }
  if (id && node.id === id) {
    return node;
  }
  if (node === document.body) {
    return node;
  }
  return findParentNodeByX(node.parentNode as Element, opts);
};

export const getNodePath = (
  node: any,
  retPaths: string[],
  treeDataMap: any,
) => {
  if (!retPaths) retPaths = [];

  retPaths.unshift(node.title);

  if (
    node.parentId > 0 &&
    treeDataMap[node.parentId] &&
    treeDataMap[node.parentId].parentId > 0
  ) {
    getNodePath(treeDataMap[node.parentId], retPaths, treeDataMap);
  }
};

export function scroll(id: string): void {
  window.console.log('scroll');
  const elem = document.querySelector(`#${id}}`);
  if (elem) {
    setTimeout(() => {
      elem.scrollTop = elem.scrollHeight + 100;
    }, 300);
  }
}

export function scrollTo(id: string, top?: number): void {
  window.console.log('scrollTo');

  const elem = document.querySelector(`#${id}}`);
  if (elem) {
    setTimeout(() => {
      elem.scrollTop = elem.scrollHeight + (top || 100);

      window.console.log(elem.scrollHeight);
    }, 500);
  }
}

export function replaceLineBreak(str: string): string {
  if (!str) return '';

  let ret = str.replaceAll(' ', '&nbsp;');
  ret = ret.replaceAll('\n', '<br />');

  return ret;
}
