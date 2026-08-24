import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import fs from "fs";
import path from "path";

// https://vite.dev/config/
export default defineConfig(({ command }) => ({
  plugins: [react()],
  base: '/app/',
  resolve:{
    alias:{
        "@": path.resolve(
            __dirname,
            "./src"
        ),
      },
  },
  server:  command === 'serve' ?{
        // host:"0.0.0.0",
        host:"localhost",
        port:5173,
        https:{
            cert: fs.readFileSync(
                "../../localhost+2.pem"
            ),

            key: fs.readFileSync(
                "../../localhost+2-key.pem"
            )
        }
    }: undefined,
}))
