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
    proxy: {
      '/': {
        target: "http://localhost:3030",
        configure: (proxy, opts) => {
          proxy.on('proxyReq', (preq, req, res) => {
            console.log(req.method)
          })
        }
      }
    },
  },
  publicDir: "../../public/init",
  build: {
    outDir: "../../dist/init"
  }
});
