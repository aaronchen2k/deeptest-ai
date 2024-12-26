export function addSepIfNeeded(url: string) {
  if (url.lastIndexOf('/') !== url.length - 1) {
    url += '/';
  }

  return url;
}

export function setServeUrl(url: string) {
  const ifNeedPrefix = !url.startsWith('http');
  return `${ifNeedPrefix ? window.location.origin : ''}${url}`;
}

export function getUrlKey(name: string, url: string) {
  const regx = new RegExp(`[?|&]${name}=` + `([^&;]+?)(&|#|;|$)`) as any;

  return (
    decodeURIComponent(
      (regx.exec(url) || [undefined, ''])[1].replaceAll('+', '%20'),
    ) || null
  );
}
