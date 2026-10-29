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

Audit dependency frontend saat setup melaporkan advisory low [GHSA-gcqx-w48m-w2h4](https://github.com/advisories/GHSA-gcqx-w48m-w2h4) pada deepmerge 4.3.1 transitif dari @sveltejs/vite-plugin-svelte 7.3.1. Pembaruan transitif diperiksa, tetapi registry masih menerbitkan 4.3.1 sebagai latest sementara advisory menyebut patch ≥4.3.2. Tidak memakai override ke versi yang tidak tersedia, tidak mengubah source dependency, dan tidak menonaktifkan audit. BUG-001 melacak pembaruan saat patch tersedia; tidak ada temuan high/critical pada audit M1-01. Keamanan domain, cloud, dan production belum dinilai melalui pemeriksaan ini.
