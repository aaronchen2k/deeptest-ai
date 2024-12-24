export function isInArray(item: any, arr: any[]) {
  for (const element of arr) {
    if (item === element) {
      return true;
    }
  }
  return false;
}
