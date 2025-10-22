# Frontend Integration (React)

This guide shows how to call the backend securely from your React app during development.

## CORS
Set `ALLOW_ORIGINS` in `.env` to your React dev URL (default is `http://localhost:3000`).

## Direct calls (no proxy)
```ts
// Simple GET
const res = await fetch('http://localhost:8080/api/v1/ping');
const data = await res.json();

// Login
const loginRes = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'admin', password: 'password123' }),
});
const { token } = await loginRes.json();

// Authenticated request
const meRes = await fetch('http://localhost:8080/api/v1/auth/me', {
  headers: { Authorization: `Bearer ${token}` },
});
```

## Using a dev proxy (recommended)
This removes cross-origin requests in development and avoids CORS preflights.

- Create React App (`src/setupProxy.js`):
```js
const { createProxyMiddleware } = require('http-proxy-middleware');
module.exports = function(app) {
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'http://localhost:8080',
      changeOrigin: true,
      secure: false,
    })
  );
};
```

- Vite (`vite.config.ts`):
```ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
});
```

Then call the API with relative URLs:
```ts
const res = await fetch('/api/v1/ping');
```

## Security notes
- Do not store JWTs in localStorage in production; prefer httpOnly cookies and CSRF protection for cookie-based auth.
- In production, host frontend and backend behind the same domain (or a trusted gateway) to simplify security and reduce CORS surface area.
