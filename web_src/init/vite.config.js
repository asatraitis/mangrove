import { defineConfig } from "vite";
import path from "path";

export default defineConfig({
  resolve: {
    alias: {
      "@dto": path.resolve(__dirname, "../../internal/dto"),
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
