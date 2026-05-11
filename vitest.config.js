import { defineConfig } from 'vitest/config';
import path from 'node:path';

export default defineConfig({
  test: {
    environment: 'jsdom',
    globals: false,
    include: ['internal/web/public/js/**/*.test.js'],
    setupFiles: ['./internal/web/public/js/_test-setup.js'],
  },
  resolve: {
    alias: {
      preact: 'preact',
      'preact/hooks': 'preact/hooks',
      '@preact/signals': '@preact/signals',
      htm: 'htm',
      sortablejs: path.resolve(
        process.cwd(),
        'internal/web/public/assets/vendor/sortablejs@1.15.6/sortable.esm.js',
      ),
    },
  },
});
