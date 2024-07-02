import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

// https://vitejs.dev/config/
export default defineConfig({
	plugins: [svelte()],
	build: {
		outDir: "../bms_server/resources/dist",
		emptyOutDir: true,
		target: "esnext",
	},
});
