import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import basicSsl from '@vitejs/plugin-basic-ssl'

export default defineConfig({
    plugins: [sveltekit(), basicSsl()],
    server: {
        proxy: {
            '/ws': {
                target: 'ws://localhost:8080',
                ws: true,
            },
            '/api': 'http://localhost:8080',
        },
    },
});
