# Architecture decisions — Unggulmart Marketing Workspace

## ADR-001 — Tooling lokal M1-01

Tanggal: 29 Oktober 2026 WITA. Status: dipakai dan diverifikasi untuk fondasi lokal; keputusan platform/cloud pada DEC-001 belum selesai.

Gunakan pnpm workspace untuk apps/web dan Go Modules mandiri untuk apps/api. Tidak menambah monorepo orchestrator. Frontend dibuat melalui sv 1.1.1 dengan template minimal TypeScript, Svelte 5 runes, Tailwind, ESLint, Prettier, Vitest, dan adapter Vercel. Fondasi shadcn-svelte 1.7.1 menggunakan preset Vega/neutral dan font Inter lokal dari package, tanpa mengaktifkan dark mode atau modul produk.

Versi terverifikasi: Node 24.22.2, pnpm 11.10.0, Go 1.27.3, SvelteKit 3.1.0, Svelte 5.57.3, TypeScript 6.0.3, Vite 8.3.5, Tailwind 4.3.3, Chi 5.5.0. Daftar dependency exact ada pada package.json web dan go.mod; dependency transitif dikunci melalui pnpm-lock.yaml dan go.sum.

Node global mesin awalnya 24.14.0. [Advisory resmi Oktober 2026](https://nodejs.org/en/blog/vulnerability/october-2026-security-releases) menyediakan patch lini 24 pada 24.22.1. Gunakan 24.22.2, yang tercantum pada [daftar rilis resmi](https://nodejs.org/en/blog/release), dan tolak 24.14.0 melalui engines serta check-tooling. Tidak menganggap satu patch menjamin bebas celah selamanya. Pemilihan ini tidak mengidentifikasi secara pasti berita viral yang dimaksud pengguna; yang dipenuhi adalah memakai lini LTS dengan patch setelah advisory terbaru yang diverifikasi.

Karena Go belum tersedia dan Node global lebih lama, pengujian memakai archive resmi Node/Go yang diekstrak lokal ke .tools/, tanpa mengubah instalasi global. SHA-256 diverifikasi terhadap [SHASUMS resmi Node](https://nodejs.org/dist/v24.22.2/SHASUMS256.txt) dan [metadata rilis Go](https://go.dev/dl/?mode=json):

| Archive yang digunakan | SHA-256 |
|---|---|
| node-v24.22.2-darwin-arm64.tar.gz | aa9549349e47dffaa42dcf95d6e06ade8c4bdafc980979883b88767670df385a |
| go1.27.3.darwin-arm64.tar.gz | dd3908cf20959d798c1e06ebe5b72213f0505ce01373b40ce3110975a6cd99c2 |

Toolchain/cache/build lokal diabaikan Git. Pengguna lain memasang versi yang sama melalui version manager/distribusi resmi; .tools/ tidak menjadi prasyarat tersembunyi di repository. Cara aktivasi dan command tersedia pada README.md.

Scaffold SvelteKit 3 menempatkan konfigurasi `sveltekit(...)` dan adapter dalam vite.config.ts, serta memakai package imports `#lib`/`#lib/*`. Ikuti hasil CLI ini daripada menambahkan svelte.config.js dari contoh versi lama; alias components.json disesuaikan ke #lib. Ownership lib/api, features, dan Go tetap mengikuti ARCHITECTURE.md. Belum ada SvelteKit server route, server action, SQL, atau backend produk kedua.

API M1-01 hanya listener Go/Chi loopback tanpa endpoint bisnis. Belum ada internal/app, adapter Vercel Go, config environment, auth, database, atau kontrak endpoint. Itu sengaja dibatasi ke M1-01; M1-02 membentuk wiring router reusable dan middleware. Adapter Vercel frontend sudah menghasilkan build lokal, tetapi routing Go, pairing environment, cookies, runtime limits dan deployment belum diverifikasi.

Audit dependency frontend saat setup melaporkan advisory low [GHSA-gcqx-w48m-w2h4](https://github.com/advisories/GHSA-gcqx-w48m-w2h4) pada deepmerge 4.3.1 transitif dari @sveltejs/vite-plugin-svelte 7.3.1. Pembaruan transitif diperiksa, tetapi registry masih menerbitkan 4.3.1 sebagai latest. Klarifikasi atribusi: ≥4.3.2 berasal dari metadata audit; advisory GitHub yang diperiksa mencantumkan belum ada patched version. Angka itu bukan bukti patch resmi tersedia. Tidak memakai override ke versi yang tidak tersedia, tidak mengubah source dependency, dan tidak menonaktifkan audit. BUG-001 tetap terbuka; tidak ada temuan high/critical pada audit M1-01. Keamanan domain, cloud, dan production belum dinilai melalui pemeriksaan ini.

## ADR-002 — Boundary HTTP dan konfigurasi M1-02

Tanggal: 29 Oktober 2026 WITA. Status: default implementasi yang dipakai dan diuji lokal; verifikasi platform/cloud tetap pada M1-10. Tidak mengubah permission, state machine atau keputusan produk terbuka.

internal/app memiliki Config/LoadConfig, constructor NewRouter, logger, dan factory cold start FromEnvironment. cmd/server dan api/index.go menggunakan constructor/factory tersebut. Tidak membuka pool, membaca DATABASE_URL/secrets, memasang route domain, atau menambahkan dependency pada task ini.

APP_ENV dan APP_ORIGIN wajib, tanpa fallback. Environment hanya development/preview/production; VERCEL_ENV opsional namun harus cocok jika ada. Development memakai origin loopback, preview/production HTTPS non-loopback. Bentuk origin canonical diperiksa tanpa request DNS/HTTP. Tidak menganggap validasi bentuk/marker sebagai bukti pasangan environment aktual; pairing web/API/database diuji kemudian.

Default launcher loopback 127.0.0.1:8080; HTTP_ADDR menolak bind ke audience publik. Request timeout default 10s (konfigurasi 1–30s), read-header timeout 5s, read/write timeout request+5s, idle 60s, max header 16 KiB. Port lokal 0 diperbolehkan untuk pengujian ephemeral. Default ini belum dibuktikan terhadap batas Vercel; adapter tidak memakai listener atau server IO timeout lokal. Seluruh variabel dan batas dicatat di README/.env.example; file .env tidak dibaca otomatis.

Config invalid menghentikan launcher sebelum listen. Adapter memakai sync.OnceValue tanpa listener/loop/IO database; kegagalan konfigurasi menjadi handler JSON 503 dan log startup berkode tetap, bukan panic/config dump. State cold start invalid dipertahankan sampai instance baru. Cron/push belum tersedia sehingga flag selain false ditolak, termasuk pada production, sampai milestone terkait mengimplementasikan lifecycle-nya.

Urutan boundary aktual: request ID server → security/no-store headers → safe access log → recovery → deadline → exact Origin guard untuk unsafe method → body limit → Chi mux. Boundary berada DI LUAR mux: test awal membuktikan fast path empty mux Chi melewati Mux.Use pada 404. Router kecil membungkus *chi.Mux agar startup dapat memasang route dan ServeHTTP selalu melewati chain yang sama. Tidak menambah dummy/public endpoint untuk mengatasi fast path tersebut.

Request ID memakai crypto/rand, diabaikan bila berasal dari header klien. Error envelope mengikuti PRD §14 dengan fields object dan pesan Bahasa Indonesia dari katalog tetap, tanpa arbitrary error text. Status 405/413/415 dan 503 timeout/cancellation diperjelas pada PRD sebagai perilaku HTTP umum. Tidak membuat kontrak endpoint/client generator sebelum baseline M1-04; endpoint domain belum tersedia.

Gunakan net/http.TimeoutHandler untuk buffering dan pembatalan context. Recovery ditempatkan di luarnya agar panic setelah partial write tidak mengirim success/Set-Cookie parsial. Guard completion menolak success sesudah context dibatalkan, termasuk race completion-versus-cancellation; late write ditolak. Cancellation parent yang menghasilkan 503 kosong dinormalisasi menjadi envelope JSON. Handler domain tetap wajib menghormati context dan tidak mengandalkan goroutine setelah response; streaming/Flusher/Hijacker bukan kebutuhan rilis ini.

Request JSON dibatasi 256 KiB dan UTF-8 valid; unknown fields, extra JSON value dan Content-Type salah ditolak tanpa detail decoder. Fixture route hanya berada pada test. Origin guard awal tidak menggantikan auth/session/CSRF dan tidak mempercayai forwarded headers atau membuka CORS. Semua endpoint privat tetap harus memasang auth saat domain tersedia.

Redaksi log dilakukan dengan allowlist metadata, bukan regex yang mencoba menebak secret. URL/path/query, IP, headers, body, error text dan panic value tidak masuk log. Kode HTTP generik, request ID, status dan latency tersedia; panic/config event memakai kode tetap. Raw transport diagnostics net/http lokal dibuang agar tidak membawa input/stack trace; lifecycle dan boundary tetap dicatat terstruktur. Tidak mengubah global slog/logger atau environment mesin.

Test memakai env/header/body sintetis dan httptest, termasuk handler Vercel serta server lokal tanpa deployment. Header Cookie/Set-Cookie fixture yang bertahan di router bukan bukti proxy cookie end-to-end. Go race detector, checks repository dan smoke TCP lokal melengkapi bukti M1-02; M1-03, CI, cloud routing, proxy cookie, session, DB dan cron tetap di luar scope. BUG-001 tetap terbuka dan dependency/pin tidak diubah.
