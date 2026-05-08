import { describe, it, expect, afterEach } from 'vitest';
import { render, screen, cleanup } from '@testing-library/preact';
import { html } from '../util/html.js';
import { Spinner } from './Spinner.js';

describe('Spinner', () => {
  afterEach(() => cleanup());
  it('renders with data-testid', () => {
    render(html`<${Spinner} />`);
    expect(screen.getByTestId('spinner')).toBeInTheDocument();
  });
});
