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

export function parseNumber(n: any, def: number = 0): number {
  try {
    const v = parseFloat(n.toString());
    return isNaN(v) ? def : v;
  } catch {
    return def;
  }
}

export function parseDate(n: any): Date | undefined {
  const date = new Date(n.toString());
  if (date instanceof Date && date.getTime && !isNaN(date.getTime())) {
    return date;
  } else {
    return undefined;
  }
}

export function getThemeColor(theme: Theme, color: any): string {
  function getPropValue(obj: any, propName: any) {
    const vals = Object.getOwnPropertyNames(obj).filter(
      (p) => propName && p.toLowerCase() === propName.toLowerCase()
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
  let hash = decodeURIComponent(window.location.hash);
  return hash.length > 0 ? hash.substring(1) : hash;
}

export function isTrue(value: any) {
  return value !== undefined && value != null && value.toString().toLowerCase() === "true";
}

export function isFalse(value: any) {
  return value !== undefined && value != null && value.toString().toLowerCase() === "false";
}

export function getId(value: any) {
  return value.toString().replace(/:/g, "_");
}