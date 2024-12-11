import { defineConfig } from "vite";

export default defineConfig({
  server:{
    port: 3001
  },
  publicDir: "../../public/init",
  build: {
    outDir: "../../dist/init"
  }
});
