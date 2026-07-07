import { describe, expect, it } from 'vitest';
import { resolveApiAssetUrl } from '@/api/client';

describe('resolveApiAssetUrl', () => {
  it('resolves upload paths against the API server root', () => {
    const resolved = new URL(resolveApiAssetUrl('/uploads/cards/priestess.png'));
    const api = new URL(import.meta.env.VITE_API_URL || '/api', window.location.origin);

    expect(resolved.origin).toBe(api.origin);
    expect(resolved.pathname).toBe('/uploads/cards/priestess.png');
  });

  it('keeps absolute URLs unchanged', () => {
    const url = 'https://cdn.example/cards/priestess.png';
    expect(resolveApiAssetUrl(url)).toBe(url);
  });
});
