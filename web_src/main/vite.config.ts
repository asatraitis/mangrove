import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from "path";
import { TanStackRouterVite } from '@tanstack/router-plugin/vite'

// https://vite.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      "@dto": path.resolve(__dirname, "../../internal/dto"),
      "@websrc": path.resolve(__dirname, "../../web_src"),
      "@services": path.resolve(__dirname, "../../web_src/services"),
    }
  },
  server:{
    port: 3000,
    proxy: {
      '/v1': {
        target: "http://localhost:3030",
        changeOrigin: true,
        secure: false,
      },
    },
  },
  publicDir: "../../public/main",
  build: {
    outDir: "../../dist/main"
  },
  plugins: [
    TanStackRouterVite(),
    react()
  ],
})
