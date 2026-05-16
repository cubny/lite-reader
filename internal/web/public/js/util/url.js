// Extract a clean host (no `www.` prefix) from a URL string.
// Returns '' on bad input so callers can short-circuit cleanly.
export function hostFromUrl(u) {
  if (!u) return '';
  try {
    return new URL(u).hostname.replace(/^www\./, '');
  } catch {
    return '';
  }
}
