const RTL_RE = /[֑-߿‏‫‮יִ-﷽ﹰ-ﻼ]/;

export function detectDir(text) {
  return text && RTL_RE.test(text) ? 'rtl' : 'ltr';
}
