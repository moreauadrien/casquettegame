import { sveltekit } from '@sveltejs/kit/vite';

/** @type {import('vite').UserConfig} */
const config = {
	plugins: [sveltekit()],
	server: {
		proxy: {
			'/ws': {
				target: 'ws://localhost:8080',
				ws: true,
			},
			'/api': 'http://localhost:8080',
		},
	},
};

export default config;
