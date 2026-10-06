import { defineConfig } from 'vite';
import cssInjectedByJsPlugin from 'vite-plugin-css-injected-by-js';
import { resolve } from 'path';

export default defineConfig(({ mode }) => {
  const isEmbed = mode === 'embed';

  return {
    plugins: [cssInjectedByJsPlugin()],
    resolve: {
      alias: {
        react: 'preact/compat',
        'react-dom': 'preact/compat',
      },
    },
    build: {
      outDir: 'dist',
      emptyOutDir: false,
      lib: isEmbed
        ? {
            entry: resolve(__dirname, 'src/embed.tsx'),
            name: 'SupportChat',
            fileName: () => 'widget.js',
            formats: ['iife'],
          }
        : {
            entry: resolve(__dirname, 'src/index.tsx'),
            name: 'SupportChatWidget',
            fileName: (format) => `index.${format}.js`,
            formats: ['es', 'cjs'],
          },
      rollupOptions: {
        external: isEmbed ? [] : ['react', 'react-dom'],
      },
    },
    server: {
      port: 3000,
    },
  };
});
