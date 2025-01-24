import { defineConfig } from "vite";
import path from "path";

export default defineConfig({
  resolve: {
    alias: {
      "@dto": path.resolve(__dirname, "../../internal/dto"),
      "@websrc":path.resolve(__dirname, "../../web_src"),
      "@services":path.resolve(__dirname, "../../web_src/services"),
    }
  },
  server:{
    port: 3001,
  },
  publicDir: "../../public/init",
  build: {
    outDir: "../../dist/init"
  }
});
