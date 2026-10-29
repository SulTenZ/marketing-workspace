# marketing-workspace

Workspace internal Marketing Unggulmart untuk task, review opsional dua arah, event, kalender WITA, reminder, notifikasi, dan PWA. Stack yang direncanakan: SvelteKit/TypeScript, Go/Chi, dan PostgreSQL Neon.

Status per 29 Oktober 2026: fondasi lokal M1-01 tersedia. Web menampilkan halaman tahap pengembangan; server Go/Chi berjalan tanpa endpoint produk. Status task/gate serta bukti terbaru dikelola hanya di [TODO.md](TODO.md).

Seluruh dokumen Markdown proyek berada di root repository:

| Dokumen | Peran |
|---|---|
| [PRD.md](PRD.md) | Kebutuhan produk, aturan bisnis, milestone, dan pertanyaan terbuka |
| [ARCHITECTURE.md](ARCHITECTURE.md) | Ownership, struktur kode, dan aliran data |
| [AGENTS.md](AGENTS.md) | Aturan agent dan invariant teknis |
| [TODO.md](TODO.md) | Task, dependensi, keputusan terbuka, progres, dan bukti |
| [SKILL.md](SKILL.md) | Panduan workflow berulang untuk proyek ini |
| [WORKFLOW.md](WORKFLOW.md) | Otonomi, izin, verifikasi, dan serah terima |

Keputusan tooling dicatat di [architecture-decisions.md](architecture-decisions.md). Dokumen deployment.md, runbook.md, dan CALENDAR-DATA.md akan dibuat di root saat dibutuhkan milestone. Kontrak API akan berada di openapi.yaml; kontrak domain belum dibuat pada M1-01.

Kode berada di apps/web dan apps/api. SQL migration direncanakan di db/migrations/, dan asset dataset Hijriyah di calendar-data/. Default produk tetap mengikuti PRD sampai dikonfirmasi sesuai tenggatnya.

## Tooling yang dikunci

| Tool / library | Versi |
|---|---|
| Node.js LTS | 24.22.2 |
| pnpm | 11.10.0 |
| Go | 1.27.3 |
| SvelteKit / Svelte | 3.1.0 / 5.57.3 |
| TypeScript / Vite | 6.0.3 / 8.3.5 |
| Tailwind CSS / shadcn-svelte | 4.3.3 / 1.7.1 |
| Chi | 5.5.0 |

Node dipilih setelah memeriksa [rilis keamanan resmi 21 Oktober 2026](https://nodejs.org/en/blog/vulnerability/october-2026-security-releases), yang menyediakan perbaikan pada lini 24 mulai 24.22.1. Proyek memakai patch berikutnya, 24.22.2. Pin ini berlaku saat verifikasi, bukan jaminan bebas kerentanan baru; tinjau advisory kembali sebelum preview/production. Node 24.14.0 ditolak oleh pemeriksaan tooling.

Pin runtime ada di .nvmrc, .node-version, .go-version dan engines package.json. Dependency frontend memakai versi exact dan pnpm-lock.yaml; Go memakai go.mod/go.sum. Gunakan pnpm, bukan lockfile dari package manager lain.

## Menjalankan lokal

Pasang Node 24.22.2 melalui version manager pilihan (misalnya `nvm install` lalu `nvm use` dari root), pnpm 11.10.0, serta Go 1.27.3 dari [distribusi resmi Go](https://go.dev/dl/). Pada mesin sesi implementasi, Node dan Go juga sudah diekstrak ke .tools/ dari archive resmi dengan checksum SHA-256 terverifikasi. Folder ini diabaikan Git dan tidak ikut clone.

Untuk menggunakan toolchain lokal yang tersedia pada mesin ini, jalankan dari root pada setiap terminal:

```sh
export PATH="$PWD/.tools/node/bin:$PWD/.tools/go/bin:$PATH"
export GOCACHE="$PWD/.cache/go-build"
export GOMODCACHE="$PWD/.cache/go-mod"
```

Export PATH tersebut tidak mengganti Node global mesin. Bila memakai runtime yang sudah dipasang lewat version manager, cukup pastikan `node --version` dan `go version` cocok dengan tabel. Cache Go di .cache/ bersifat lokal dan tidak ikut Git.

Install dependency dan periksa versi:

```sh
pnpm install --frozen-lockfile
(cd apps/api && go mod download)
pnpm check:tooling
```

Terminal pertama menjalankan web:

```sh
pnpm dev:web --strictPort
```

Buka http://127.0.0.1:5173/. Terminal kedua menjalankan API:

```sh
pnpm dev:api
```

API mendengarkan pada 127.0.0.1:8080. Respons 404 pada `/` memang diharapkan: M1-01 hanya membuktikan server hidup, belum menyediakan health/auth/domain API. Belum ada akses database, proxy, migration, cron, push, atau credential yang diperlukan. Hentikan server dengan Ctrl-C. Pada sandbox yang membatasi port/socket, perintah dev dan smoke test memerlukan izin loopback.

## Pemeriksaan lokal

| Command dari root | Fungsi |
|---|---|
| `pnpm check:tooling` | Menolak runtime Node/Go yang berbeda dari pin |
| `pnpm format` / `pnpm format:api` | Merapikan source web / Go |
| `pnpm lint` | Prettier check dan ESLint web |
| `pnpm check` | Tooling, Svelte/TypeScript, gofmt check, go vet |
| `pnpm test` | Vitest dan go test |
| `pnpm build` | Build web dengan adapter Vercel dan binary Go lokal di .cache/ |
| `pnpm quality` | Seluruh checks, lint, test dan build di atas |
| `pnpm test:smoke` | Memeriksa SSR web dan respons Go, ketika kedua server lokal berjalan |
| `pnpm audit --audit-level high` | Memeriksa advisory dependency frontend pada registry |

Scaffold belum memiliki unit test bisnis: Vitest mengizinkan nol test file, dan go test melaporkan `[no test files]`. Itu tidak membuktikan permission/auth/transaction sudah diuji; smoke test hanya membuktikan startup dan respons lokal. Test bisnis ditambahkan bersama domain terkait. Svelte autofixer tidak tersedia pada sesi ini; compiler/svelte-check dipakai sebagai pemeriksaan alternatif.

Visual halaman fondasi diperiksa pada Chrome desktop. Sesi browser otomatis mencatat warning hydration pada favicon yang berubah menjadi SVG dengan atribut `data-codex-favicon-badge`; markup tersebut tidak berasal dari source proyek. Pemeriksaan ini bukan klaim runtime browser tanpa warning atau uji mobile; verifikasi kembali tanpa modifikasi favicon dari tooling saat menguji UI/auth berikutnya.

Audit 29 Oktober 2026 tidak menemukan high/critical, tetapi melaporkan satu advisory low pada deepmerge 4.3.1, dependency build transitif plugin Svelte. Versi patch yang disebut advisory (≥4.3.2) belum tersedia di registry saat verifikasi. Temuan dan tindak lanjut dilacak sebagai BUG-001 di TODO; dependency tidak dipatch sendiri atau advisory disembunyikan.

M1-02 belum dikerjakan. Wiring router reusable, config fail-closed, request ID, error convention dan runtime backend mengikuti task tersebut. Build adapter web lokal bukan bukti deployment Vercel, proxy cookie, atau kesiapan production; gate itu tetap terbuka di TODO.
