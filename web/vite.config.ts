import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig({
	// Vite offers a way to copy files as-is into
	// the dist/ folder. We disable that feature here.
	publicDir: false,
	build: {
		rollupOptions: {
			input: {
				index: "index.html",
				404: "404.html",
			},
		},
	},
	plugins: [react()],
	server: {
		proxy: {
			"/api": {
				target: "http://localhost:8080",
				changeOrigin: true,
				secure: false,
			},
		},
	},
});
