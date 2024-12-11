import { defineConfig } from "vite";

export default defineConfig({
  server:{
    port: 3000,
  },
  publicDir: "../../public/main",
  build: {
    outDir: "../../dist/main"
  }
});
