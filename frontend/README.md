# PLN Frontend - React + Vite

Frontend aplikasi PLN yang dibangun dengan **React 19** dan **Vite**, terintegrasi dengan Go backend yang mengimplementasikan SSDLC (Secure Software Development Lifecycle).

## 🚀 Teknologi

- **React 19.2.0** - Library UI modern
- **Vite** - Build tool yang sangat cepat
- **ESLint** - Linting untuk code quality
- **CSS3** - Styling dengan glassmorphism effect

## 📋 Prerequisites

- Node.js (v18 atau lebih baru)
- npm atau yarn
- Go backend berjalan di port 8080

## 🛠️ Installation

```bash
# Install dependencies
npm install
```

## 🎯 Available Scripts

### `npm run dev`

Menjalankan aplikasi dalam mode development.
Buka [http://localhost:3000](http://localhost:3000) untuk melihatnya di browser.

Hot Module Replacement (HMR) aktif - perubahan akan langsung terlihat tanpa refresh penuh.

### `npm run build`

Build aplikasi untuk production ke folder `build/`.
Build akan dioptimasi untuk performa terbaik.

```bash
npm run build
```

### `npm run preview`

Preview hasil build production secara lokal.

```bash
npm run preview
```

### `npm run lint`

Menjalankan ESLint untuk memeriksa kualitas kode.

```bash
npm run lint
```

## 🔒 Security Features (SSDLC Integration)

Frontend ini terintegrasi dengan backend Go yang mengimplementasikan:

- ✅ **Secure Headers** - CSP, HSTS, X-Frame-Options
- ✅ **JWT Authentication** - Token-based authentication
- ✅ **Rate Limiting** - Proteksi dari brute force
- ✅ **Input Validation** - Validasi di client dan server
- ✅ **CORS Protection** - Cross-Origin Resource Sharing yang aman
- ✅ **XSS Prevention** - Proteksi dari Cross-Site Scripting

## 📁 Struktur Folder

```
frontend/
├── public/           # Static assets
│   ├── manifest.json
│   └── robots.txt
├── src/
│   ├── App.js       # Main component
│   ├── App.css      # Main styles
│   ├── index.js     # Entry point
│   └── index.css    # Global styles
├── index.html       # HTML template (di root untuk Vite)
├── vite.config.js   # Vite configuration
├── eslint.config.js # ESLint configuration
└── package.json     # Dependencies
```

## 🔗 API Integration

Frontend berkomunikasi dengan backend melalui proxy yang dikonfigurasi di `vite.config.js`:

```javascript
proxy: {
  '/api': {
    target: 'http://localhost:8080',
    changeOrigin: true,
    secure: false,
  }
}
```

Contoh pemanggilan API:

```javascript
// Health check
fetch('/api/health')
  .then(res => res.json())
  .then(data => console.log(data));

// Login
fetch('/api/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username, password })
});
```

## 🎨 Styling

Aplikasi menggunakan modern CSS dengan:
- Gradient backgrounds
- Glassmorphism effects
- Responsive design
- Smooth animations

## 📝 Environment Variables

Buat file `.env.local` untuk konfigurasi environment:

```env
VITE_API_URL=http://localhost:8080
VITE_APP_NAME=PLN Secure App
```

Akses dalam kode:
```javascript
const apiUrl = import.meta.env.VITE_API_URL;
```

## 🚦 Menjalankan Full Stack

1. **Start Backend (Terminal 1)**:
   ```bash
   cd backend
   go run cmd/server/main.go
   ```

2. **Start Frontend (Terminal 2)**:
   ```bash
   cd frontend
   npm run dev
   ```

3. Buka browser di http://localhost:3000

## 📦 Build untuk Production

```bash
# Build
npm run build

# Preview build
npm run preview
```

Files production akan ada di folder `build/`.

## 🔧 Troubleshooting

### Port 3000 sudah digunakan
Edit `vite.config.js` dan ubah port:
```javascript
server: {
  port: 3001, // ganti port
}
```

### Backend tidak terkoneksi
- Pastikan Go backend berjalan di port 8080
- Check CORS configuration di backend
- Periksa proxy settings di `vite.config.js`

## 📚 Learn More

- [Vite Documentation](https://vitejs.dev/)
- [React Documentation](https://react.dev/)
- [ESLint Documentation](https://eslint.org/)

## 👨‍💻 Development

Dibuat dengan ❤️ mengikut best practices SSDLC
