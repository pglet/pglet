import { Theme, SharedColors } from "@fluentui/react";

export function defaultPixels(size: any) {
  if (!size) {
    return size;
  }

  const m = size.toString().match(/^\d*(\.\d+)?$/);
  if (m) {
    // just number
    return `${size}px`;
  }
  return size;
}

export function parseNumber(n: any): number {
  try {
    const v = parseFloat(n.toString());
    return isNaN(v) ? 0 : v;
  } catch {
    return 0;
  }
}

export function getThemeColor(theme: Theme, color: any): string {
  function getPropValue(obj: any, propName: any) {
    const vals = Object.getOwnPropertyNames(obj).filter(
      (p) => p.toLowerCase() === propName.toLowerCase()
    );
    if (vals.length > 0) {
      return obj[vals[0]] ? obj[vals[0]].toString() : "";
    }
    return "";
  }

  let result = getPropValue(theme.palette, color);
  if (result === "") {
    result = getPropValue(SharedColors, color);
  }
  return result !== "" ? result : color;
}

export function getWindowHash() {
  let hash = window.location.hash;
  return hash.length > 0 ? hash.substring(1) : hash;
}
