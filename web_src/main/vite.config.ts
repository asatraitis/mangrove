import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from "path";

// https://vite.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      "@dto": path.resolve(__dirname, "../../internal/dto"),
    }
  },
  server:{
    port: 3000,
  },
  publicDir: "../../public/main",
  build: {
    outDir: "../../dist/main"
  },
  plugins: [react()],
})
