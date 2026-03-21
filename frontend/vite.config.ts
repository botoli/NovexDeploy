import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5432, // ваш текущий порт
    proxy: {
      "/auth": {
        target: "http://localhost:8888",
        changeOrigin: true,
      },
      "/api": {
        target: "http://localhost:8888",
        changeOrigin: true,
      },
      "/health": {
        target: "http://localhost:8888",
        changeOrigin: true,
      },
      "/projects": {
        target: "http://localhost:8888",
        changeOrigin: true,
      },
      "/users": {
        target: "http://localhost:8888",
        changeOrigin: true,
      },
      "/metrics": {
        target: "http://localhost:8888",
        changeOrigin: true,
      },
    },
  },
});
