import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';

export default defineConfig({
  base: process.env.VITE_BASE_PATH || '/ui/',
  plugins: [react()],
  server: {
    host: '127.0.0.1',
    port: 8080,
    proxy: Object.fromEntries(['/api', '/es', '/version', '/swagger'].map(path => [path, {
      target: process.env.ZINC_API_TARGET || 'http://localhost:4080',
      changeOrigin: true,
    }])),
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    clearMocks: true,
    coverage: { reporter: ['text', 'html'] },
  },
});
