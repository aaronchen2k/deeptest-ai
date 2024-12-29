import { scrollTo } from '@vben/utils';

export function scroll() {
  setTimeout(() => {
    scrollTo('chat-messages', 0);
  }, 200);
}
export function getDocLink(doc: any): any {
  return {
    docId: doc.document_id,
    docName: doc.document_name,
    docType: doc.data_source_type,
    datasetId: doc.dataset_id,
    datasetName: doc.dataset_name,
  };
}

export function replaceImageUrl(str: string, imageRepoUrl: string) {
  window.console.log('replaceImageUrl');
  try {
    // ![海报](./img/poster.jpg \"海报\")
    str = str.replaceAll(
      /(\[.+?\])\(([^http].+?)\)/g,
      `$1(${imageRepoUrl}/$2)`,
    );
  } catch (error) {
    window.console.log('replaceImageUrl error', error);
  }

  window.console.log('******', str);

  return str;
}

export function replaceLinkWithoutTitle(str: string) {
  window.console.log('replaceLinkWithoutTitle');
  try {
    // // html page
    // str = str.replaceAll(
    //   /\[(\d+)-([^\]]+)\]\([^)]+\.html\)[\s\S]/g,
    //   '[$2](https://wiki.deeptestcloud.com/pages/viewpage.action?pageId=$1)',
    // );
    //
    // // diffpagesbyversion page
    // // ABC (/pages/diffpagesbyversion.action?pageId=5969977&selectedPageVersions=1&selectedPageVersions=2) 123
    // str = str.replaceAll(
    //   /([^\]])\((\/pages\/.+?\.action\?pageId=.+?)\)/g,
    //   '$1[链接](https://wiki.deeptestcloud.com$2)',
    // );
    //
    // change markdown link to html link.
    // str = str.replace(/([^\]])\((http.+?)\)/g, '$1<a href="$2" target="_blank">$2</a>')

    return str;
  } catch (error) {
    window.console.log('replaceLinkWithoutTitle error', error);
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

export function getMaterialIdInDocName(docName: string) {
  // 百年孤独_8.md
  const start = docName.lastIndexOf('@');
  const end = docName.lastIndexOf('.');

  const materialId = docName.slice(start + 1, end);

  const newDocName = `${docName.slice(0, start)}${docName.slice(end)}`;

  return { materialId, newDocName };
}
