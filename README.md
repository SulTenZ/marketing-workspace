# marketing-workspace

Workspace internal Marketing Unggulmart untuk task, review opsional dua arah, event, kalender WITA, reminder, notifikasi, dan PWA. Stack yang direncanakan: SvelteKit/TypeScript, Go/Chi, dan PostgreSQL Neon.

Status per 29 Oktober 2026: fondasi lokal M1-01/M1-02 tersedia. Web menampilkan halaman tahap pengembangan; lokal dan adapter Vercel Go memakai router HTTP bersama, konfigurasi eksplisit, request ID, error JSON, timeout dan log aman. Endpoint produk belum tersedia. Status task/gate serta bukti terbaru dikelola hanya di [TODO.md](TODO.md).

Seluruh dokumen Markdown proyek berada di root repository:

| Dokumen | Peran |
|---|---|
| [PRD.md](PRD.md) | Kebutuhan produk, aturan bisnis, milestone, dan pertanyaan terbuka |
| [ARCHITECTURE.md](ARCHITECTURE.md) | Ownership, struktur kode, dan aliran data |
| [AGENTS.md](AGENTS.md) | Aturan agent dan invariant teknis |
| [TODO.md](TODO.md) | Task, dependensi, keputusan terbuka, progres, dan bukti |
| [SKILL.md](SKILL.md) | Panduan workflow berulang untuk proyek ini |
| [WORKFLOW.md](WORKFLOW.md) | Otonomi, izin, verifikasi, dan serah terima |

Keputusan tooling dan boundary HTTP dicatat di [architecture-decisions.md](architecture-decisions.md). Dokumen deployment.md, runbook.md, dan CALENDAR-DATA.md akan dibuat di root saat dibutuhkan milestone. Kontrak API akan berada di openapi.yaml pada M1-04; error umum mengikuti PRD §14. Belum ada endpoint domain/client API yang digenerate.

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

Buka http://127.0.0.1:5173/. Pada terminal kedua, aktifkan toolchain yang sama, lalu isi konfigurasi development secara eksplisit sebelum menjalankan API:

```sh
export APP_ENV=development
export APP_ORIGIN=http://127.0.0.1:5173
export PUSH_ENABLED=false
export CRON_ENABLED=false
pnpm dev:api
```

API mendengarkan pada 127.0.0.1:8080. GET `/` atau path domain yang belum tersedia menghasilkan 404 JSON dengan `error.code=NOT_FOUND` dan `error.request_id` yang sama dengan header X-Request-ID. POST/PATCH/PUT/DELETE tanpa Origin yang cocok ditolak 403, termasuk ke route yang belum tersedia. Belum ada health/auth/domain API, akses database, proxy, migration, cron, push, atau credential yang diperlukan. Hentikan server dengan Ctrl-C/SIGTERM; launcher membatalkan context request dan memberi shutdown budget 5 detik. Pada sandbox yang membatasi port/socket, perintah dev dan smoke test memerlukan izin loopback.

## Konfigurasi HTTP M1-02

[.env.example](.env.example) hanya referensi placeholder. Launcher tidak membaca .env secara otomatis; isi environment proses melalui export seperti contoh lokal atau konfigurasi deployment. Tidak ada fallback ke production/development jika nilai wajib hilang atau salah. Jangan menyimpan nilai secret pada file yang ikut Git.

| Variabel | Aturan |
|---|---|
| APP_ENV | Wajib: development, preview, atau production |
| APP_ORIGIN | Wajib: satu origin canonical, tanpa path, trailing slash, query, fragment, userinfo, wildcard, atau default port eksplisit. Development memakai localhost/IP loopback; preview/production memakai HTTPS non-loopback. |
| VERCEL_ENV | Opsional untuk host portable/lokal; bila disediakan provider harus sama dengan APP_ENV. Nilai kosong ditolak. |
| HTTP_ADDR | Opsional; default 127.0.0.1:8080, wajib IP loopback dan port 0–65535. Port 0 untuk port lokal otomatis. Hanya launcher lokal yang menggunakannya; adapter tidak membuat listener. |
| REQUEST_TIMEOUT | Opsional; default 10s, rentang 1s–30s |
| LOG_LEVEL | Opsional; default info; debug/info/warn/error |
| PUSH_ENABLED / CRON_ENABLED | Belum diimplementasikan: tidak diisi atau false. true, kosong, dan nilai lain ditolak pada M1-02, termasuk production. |

String kosong pada konfigurasi opsional yang diisi juga ditolak. Konfigurasi invalid membuat launcher berhenti dengan exit 1 sebelum membuka listener; adapter serverless mengembalikan 503 SERVICE_UNAVAILABLE dengan request ID dan tanpa membocorkan nilainya. Adapter menginisialisasi handler sekali per instance; perubahan environment memerlukan instance/cold start baru.

Batas IO berasal dari konfigurasi server. Launcher memakai read-header timeout 5s, read/write timeout REQUEST_TIMEOUT+5s, idle timeout 60s, dan batas header 16 KiB. Batas request JSON 256 KiB; decoder menolak UTF-8 invalid, unknown fields, dan lebih dari satu nilai JSON. Handler domain berikutnya tetap wajib memvalidasi field bisnis dan menggunakan decoder bersama. Respons 405 menyertakan Allow untuk metode route yang sah, termasuk nested route.

Timeout/cancellation menghasilkan JSON 503 REQUEST_TIMEOUT/REQUEST_CANCELLED. Response handler dibuffer oleh net/http.TimeoutHandler; partial success/panic atau late write tidak boleh mengganti error yang committed. Handler tetap wajib meneruskan context ke IO/SQL dan berhenti ketika dibatalkan. Goroutine request standar ini tidak dipakai sebagai worker sesudah response, tidak menyediakan streaming/Hijacker/Flusher, dan bukan scheduler/background job.

Log JSON hanya memuat metadata yang diizinkan: waktu/level/message, request ID, action, method standar, status, error_code generik HTTP_XXX, latency, serta kode startup/panic tetap. Tidak mencatat URL/path/query, IP/forwarded headers, body, Cookie/Set-Cookie, Authorization/CSRF, arbitrary error atau panic value. Request ID dibuat server; header request ID dari klien diabaikan. Boundary tidak mempercayai forwarded headers dan tidak membuka CORS. Origin guard ini belum menggantikan session/CSRF/authorization yang dikerjakan di milestone auth.

Pemeriksaan APP_ENV/VERCEL_ENV dan bentuk origin tidak membuktikan pasangan web/API benar secara operasional. Go/Chi routing cloud, pairing preview, Cookie/Set-Cookie melalui proxy, serta batas durasi Vercel tetap harus dibuktikan pada M1-10. Batas timeout lokal belum dianggap verifikasi runtime cloud.

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
| `pnpm test:smoke` | Memeriksa SSR web, API JSON/request ID/no-store dan penolakan Origin, ketika kedua server lokal berjalan |
| `pnpm audit --audit-level high` | Memeriksa advisory dependency frontend pada registry |

Go kini memiliki test config, boundary HTTP, adapter, dan launcher: empty-router 404, error/Origin, header/cookie forwarding fixture, redaksi, panic/partial response, timeout/cancellation/late write, body limit/UTF-8/unknown fields serta request ID konkuren. Jalankan `go test -race ./...` dari apps/api untuk race detector. Test tidak memakai database atau route fixture pada runtime aplikasi.

Unit test bisnis dan frontend belum ada: Vitest masih mengizinkan nol test file. Test foundation/smoke ini belum membuktikan permission/auth/transaction atau readiness cloud; test bisnis ditambahkan bersama domain terkait. Svelte autofixer tidak tersedia pada sesi ini; compiler/svelte-check dipakai sebagai pemeriksaan alternatif.

Visual halaman fondasi diperiksa pada Chrome desktop. Sesi browser otomatis mencatat warning hydration pada favicon yang berubah menjadi SVG dengan atribut `data-codex-favicon-badge`; markup tersebut tidak berasal dari source proyek. Pemeriksaan ini bukan klaim runtime browser tanpa warning atau uji mobile; verifikasi kembali tanpa modifikasi favicon dari tooling saat menguji UI/auth berikutnya.

Audit M1-01 pada 29 Oktober 2026 tidak menemukan high/critical, tetapi melaporkan satu advisory low pada deepmerge 4.3.1, dependency build transitif plugin Svelte. Metadata audit saat itu menyebut ≥4.3.2; [advisory GitHub yang diperiksa](https://github.com/advisories/GHSA-gcqx-w48m-w2h4) mencantumkan belum ada patched version. Jangan menganggap 4.3.2 sudah tersedia atau terbukti memperbaiki temuan. BUG-001 tetap terbuka; tidak ada dependency, lockfile, atau pin runtime yang diubah pada M1-02, dan tidak ada klaim audit ulang dependency pada task ini.

M1-03 belum dimulai. Tidak ada database/migration/session atau endpoint domain yang ditambahkan. Build dan test adapter lokal bukan bukti deployment Vercel, proxy cookie, atau kesiapan production; gate itu tetap terbuka di TODO. Tidak ada perubahan konfigurasi global mesin, commit, push atau deploy pada pengerjaan M1-02.
