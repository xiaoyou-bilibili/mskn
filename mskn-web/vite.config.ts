import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { resolve } from "path";

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [react()],
    resolve: {
        alias: {
            "@": resolve(__dirname, "src"), // 配置路径映射
        },
    },
    server: {
        proxy: {
            "/api": {
                // target: "http://127.0.0.1:8080",
                target: "http://mskn.xiaoyou.host",
                changeOrigin: true,
                // rewrite: (path) => path.replace(/^\/api/, ""),
            },
        },
    },
    optimizeDeps: {
        include: [
            `monaco-editor/esm/vs/language/json/json.worker`,
            `monaco-editor/esm/vs/language/css/css.worker`,
            `monaco-editor/esm/vs/language/html/html.worker`,
            `monaco-editor/esm/vs/language/typescript/ts.worker`,
            `monaco-editor/esm/vs/editor/editor.worker`
        ],
    },
})
