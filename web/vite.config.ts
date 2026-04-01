import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/locations': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/rooms': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/landlords': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/leases': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/bills': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/deposits': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/print': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/income': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/operation-logs': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/oauth2': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/health': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/metrics': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
})
