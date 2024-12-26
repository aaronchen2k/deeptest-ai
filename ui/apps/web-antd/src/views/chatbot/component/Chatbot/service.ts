import { scrollTo } from '@vben/utils';

export function scroll() {
  setTimeout(() => {
    scrollTo('chat-messages', 0);
  }, 200);
}

export function getDocDesc(str: string) {
  if (str.length < 30) {
    return str;
  }

  const first = str.slice(0, 16);
  const last = str.slice(Math.max(0, str.length - 8));

  return `${first} ... ${last}`;
}

export function getDocLink(source: any): any {
  // docs/66666666-API - 20226666 - 文档中心 - 用户知识库.html

  // eslint-disable-next-line regexp/no-super-linear-backtracking
  const regex = /.+?(\d+)-(.+?)-.*\.(html)/g;

  const matches = regex.exec(source) as any;
  if (matches && matches.length > 3) {
    return {
      pageId: matches[1],
      pageTitle: matches[2].trim(),
      pageType: matches[3].trim(),
    };
  }

  return {};
}

export function replaceLinkWithoutTitle(str: string) {
  window.console.log('replaceLinkWithoutTitle');
  try {
    // html page
    str = str.replaceAll(
      /\[(\d+)-([^\]]+)\]\([^)]+\.html\)[\s\S]/g,
      '[$2](https://wiki.deeptestcloud.com/pages/viewpage.action?pageId=$1)',
    );

    // diffpagesbyversion page
    // ABC (/pages/diffpagesbyversion.action?pageId=5969977&selectedPageVersions=1&selectedPageVersions=2) 123
    str = str.replaceAll(
      /([^\]])\((\/pages\/.+?\.action\?pageId=.+?)\)/g,
      '$1[链接](https://wiki.deeptestcloud.com$2)',
    );

    // change markdown link to html link.
    // str = str.replace(/([^\]])\((http.+?)\)/g, '$1<a href="$2" target="_blank">$2</a>')

    return str;
  } catch (error) {
    window.console.log('replace error', error);
  }
}

export const getLatestRobotMsg = function (msgs: any) {
  if (msgs.length === 0) return -1;

  for (let i = msgs.length - 1; i >= 0; i--) {
    if (msgs[i].type === 'robot') {
      return i;
    }
  }

  return -1;
};

export const setSelectionRange = function (ctrl: any, pos: any) {
  window.console.log('setSelectionRange', ctrl, pos);

  setTimeout(() => {
    if (ctrl.setSelectionRange) {
      ctrl.focus();
      ctrl.setSelectionRange(-1, -1);
    } else if (ctrl.createTextRange) {
      const range = ctrl.createTextRange();
      range.collapse(true);
      range.moveEnd('character', pos);
      range.moveStart('character', pos);
      range.select();
    }
  }, 100);
};

export function isUnderRobotMsg(elem: any) {
  const parent = elem.parentNode;
  if (!parent) {
    return false;
  }

  if (parent.classList.contains('markdown-container')) {
    return true;
  }

  return isUnderRobotMsg(parent);
}
