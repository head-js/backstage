import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import wails from "@wailsio/runtime/plugins/vite";

const port = Number(process.env.WAILS_VITE_PORT) || 5173;

export default defineConfig({
  plugins: [react(), wails("/bindings")],
  root: "src",
  clearScreen: false,
  server: {
    host: "127.0.0.1",
    port,
    strictPort: true,
  },
  build: {
    outDir: "../src-wails/frontend/dist",
    emptyOutDir: true,
  },
});
