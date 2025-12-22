import houdini from 'houdini/vite';
import devtoolsJson from 'vite-plugin-devtools-json';
import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [houdini(), tailwindcss(), sveltekit(), devtoolsJson()],
	server: {
		proxy: {
			'/query': 'http://127.0.0.1:8080/query'
		}
	}
});
