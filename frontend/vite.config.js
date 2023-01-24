import { sveltekit } from '@sveltejs/kit/vite';
import fs from 'fs'

const options = {
    key: fs.readFileSync('./localhost.key'),
    cert: fs.readFileSync('./localhost.crt'),
}

/** @type {import('vite').UserConfig} */
const config = {
	plugins: [sveltekit()],
	server: {
        https: options,
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
