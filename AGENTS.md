# AGENTS.md

Panduan ini berlaku untuk AI agent dan Codex yang bekerja di seluruh monorepo `marketing-workspace`.

## Source of Truth & Scope

- Baca `PRD.md` dan kode domain terkait sebelum mengubah behavior. PRD menjelaskan kebutuhan produk; file ini mengatur cara mengerjakannya.
- Seluruh dokumen Markdown proyek berada di root, termasuk `ARCHITECTURE.md`, `TODO.md`, `SKILL.md`, `WORKFLOW.md`, serta dokumen pendukung yang dibuat kemudian. Gunakan lokasi dan kapitalisasi tersebut; kontrak API berada di `openapi.yaml`. Kode dan migration tetap mengikuti struktur apps/ dan db/migrations/.
- Ikuti milestone PRD dan scope tugas aktif. Jangan membangun milestone berikutnya tanpa kebutuhan dari tugas.
- Bedakan keputusan tetap, default implementasi, dan pertanyaan terbuka dalam PRD. Jangan mengubah asumsi menjadi keputusan pengguna secara diam-diam.
- Jika dokumen dan kode bertentangan, jelaskan konflik; jangan mempertahankan bug hanya karena sudah ada atau mengubah aturan produk hanya agar test lulus.
- Periksa instruksi lokal tambahan jika tersedia. Instruksi subdirektori memperinci area tersebut tanpa menghapus invariant produk di root.

## Product Boundary

- Aplikasi ini adalah workspace internal Marketing Unggulmart untuk task, review, event, kalender, reminder, notifikasi, dan activity history.
- Dashboard menampilkan aktivitas pekerjaan. Jangan menambah Sales/Product Analytics, Campaign Analytics, CRM, Basket Analysis, Competitor Intelligence, Employee Engagement Tracker, HR, atau reporting BOD khusus.
- Tidak ada integrasi Affari/Google Calendar, WhatsApp automation, AI chatbot, native app, atau offline editing.
- Link hasil pekerjaan diperbolehkan; upload file, recurring task/event, multi-reviewer, dan fitur tambahan mengikuti pembatasan PRD.
- Bangun untuk satu tim kecil. Jangan menambahkan multi-tenancy, microservices, Redis, message broker, Kubernetes, Nx, atau Turborepo tanpa kebutuhan terukur dan keputusan arsitektur.

## Architecture & Ownership

- Pertahankan SvelteKit + TypeScript, Tailwind CSS, shadcn-svelte, FullCalendar, Go + Chi, pgx/pgxpool, dan PostgreSQL Neon.
- `apps/web` menangani UI; `apps/api` memiliki business logic, authorization, validasi transisi, dan persistence; `db/migrations` memiliki perubahan schema.
- Frontend tidak mengakses database dan tidak membawa database credential.
- SvelteKit server route hanya boleh menjadi adapter/proxy tipis bila diperlukan untuk same-origin API. Jangan membuat backend produk kedua melalui server actions/routes.
- Semua HTTP frontend melewati `src/lib/api`. Komponen tidak membuat jalur autentikasi, parsing error, atau request convention sendiri.
- `openapi.yaml` menjadi kontrak API saat tersedia. Lengkapi kontrak bersama perubahan endpoint; generate/reuse types, jangan menduplikasi model API di beberapa folder.
- Dokumentasikan keputusan penting di `architecture-decisions.md`; jangan mengganti stack atau menambah infrastruktur untuk menghindari memahami implementasi yang sudah ada.

## Svelte Architecture

- Untuk setup baru, gunakan Svelte 5 runes; pada repository yang sudah terbentuk, periksa versi dan pola existing sebelum mengubahnya.
- `src/routes` menyusun page/layout dan navigasi. UI domain dan orkestrasi query/mutation berada di `src/lib/features/<domain>`; komponen bersama di `src/lib/components`.
- Pisahkan server state dari state UI seperti dialog, filter, dan draft form. Jangan menjadikan localStorage sebagai sumber kebenaran pekerjaan atau sesi.
- Jangan memasang TanStack Query atau state library hanya karena dipakai proyek lain. Gunakan pendekatan existing dan tambahkan dependency hanya jika masalahnya jelas.
- Integrasikan FullCalendar pada client lifecycle, bersihkan instance saat unmount, dan gunakan fitur standard sesuai kebutuhan PRD.
- Tangani loading, empty, error, success, dan version conflict. Refresh data setelah mutasi; polling hanya saat tab aktif sesuai PRD.
- Environment variable public hanya untuk nilai yang aman berada di browser. Origin backend privat tetap konfigurasi server/proxy.

## Go Backend

- Susun domain `auth`, `users`, `tasks`, `events`, dan `notifications` dengan handler, service, repository, serta model yang sederhana.
- Handler menangani HTTP; service menjalankan aturan bisnis/transaksi; repository menjalankan SQL berparameter. Hindari generic repository atau interface berlapis tanpa kebutuhan.
- Ambil actor dari sesi, bukan user_id yang dipercaya dari request. Validasi input, unknown fields, permission, state, dan version di server.
- Mutasi bisnis, audit history, dan penulisan inbox/outbox yang relevan harus commit dalam transaksi yang sama.
- Gunakan context timeout/cancellation, error terstruktur, dan request_id. Jangan membocorkan SQL, stack trace, credential, atau isi sensitif ke response/log.
- Reuse router `http.Handler` untuk adapter Vercel dan server lokal. Jangan menjalankan HTTP listener atau background loop permanen di handler serverless.
- Jangan mengandalkan goroutine setelah response untuk reminder atau pengiriman push yang harus persisten.

## Roles, Tasks & Review Invariants

- Supervisor memiliki seluruh kemampuan Staff ditambah administrasi. Permission operasional mengikuti creator, PIC, dan reviewer sesuai matrix PRD.
- Semua pengguna aktif dapat melihat pekerjaan tim; hak baca tidak otomatis memberi hak edit.
- Staff boleh mereview pekerjaan Supervisor. Jangan mengubah review menjadi approval hierarkis wajib.
- `requires_review` default OFF. PIC dapat melakukan direct completion melalui transisi sah ketika review tidak diwajibkan.
- `Mark as Done` dan `Approve` adalah aksi berbeda. Simpan completed_by, completed_at, dan completion_method `direct` atau `review-approved` dengan benar.
- Status hanya `todo`, `in_progress`, `review`, dan `done`. Overdue adalah kondisi turunan; archived adalah flag. Jangan menambah status untuk menggantikan keduanya.
- Gunakan endpoint aksi dan state machine backend, bukan PATCH status bebas. Drag Kanban tunduk pada aturan yang sama.
- Review membutuhkan reviewer aktif selain PIC/pengaju. Self-approval dilarang; self-completion melalui jalur direct diperbolehkan.
- Saat Review, blok perubahan substantif sampai cancel review. Saat Done, perubahan substantif membutuhkan reopen. Komentar mengikuti izin PRD.
- Setiap resubmit membuat review round baru. Approval memvalidasi round aktif, reviewer, content revision, dan version; jangan menerima keputusan dari notifikasi lama.
- Request Changes, cancel, reopen, perubahan kebijakan review, dan administrative override membutuhkan alasan/komentar sesuai PRD.
- Reopen membersihkan completion terkini, bukan menghapus sejarah. Administrative override tidak boleh memalsukan reviewer atau memintas invariant approval.
- Gunakan optimistic concurrency, transaksi, dan lock sesuai risiko. Uji approve-versus-cancel dan retry; jangan hanya mengandalkan tombol disabled.

## Calendar & Reminders

- Zona aplikasi selalu `Asia/Makassar`. Jangan mengikuti timezone perangkat atau host secara implisit.
- Simpan instant sebagai `timestamptz`; date-only sebagai `date`. Jangan mengonversi date-only menjadi UTC midnight yang dapat menggeser tanggal.
- Akhir all-day event eksklusif. Task date-only overdue setelah hari tersebut berakhir di WITA.
- Deadline task tampil sebagai projection kalender dari task asli; jangan membuat record calendar_event duplikat.
- Gregorian adalah sumber tanggal. Hijriyah informasi sekunder dengan metode, cakupan, dan sumber yang jelas; fallback harus berlabel perkiraan.
- Jangan menyamakan konversi Hijriyah algoritmis dengan penetapan resmi hari besar Indonesia atau menggeser event otomatis saat referensi Hijriyah berubah.
- Default reminder event besar H-7/H-3/H-1 dan event kecil H-1 mengikuti jam/penerima PRD. H-N berbasis hari kalender WITA, bukan sekadar kelipatan 24 jam.
- Reschedule/archive membatalkan reminder lama yang belum dikirim. Hormati revision, dedupe key, dan aturan missed reminder; jangan mengirim seluruh backlog secara mendadak.
- Cron wajib autentikasi khusus, bounded batch, idempotency, dan recovery. Jangan mengasumsikan scheduler selalu tepat waktu atau otomatis retry.

## Authentication, PWA & Notifications

- Gunakan server-side session dengan cookie HttpOnly/Secure sesuai environment; jangan memakai JWT/session token di localStorage.
- Unsafe browser requests wajib validasi Origin dan CSRF sesuai kontrak backend. CORS dan client route guard bukan authorization.
- Hanya 401 sesi yang memicu login ulang; tangani 403, 409, 422, 429, network, dan 5xx sesuai maknanya.
- Tidak ada registrasi publik. Bootstrap dan recovery akun melalui jalur operasional terkontrol, bukan endpoint publik tersembunyi.
- Deaktivasi, reset password, role change, dan perlindungan Supervisor aktif terakhir mengikuti PRD; jangan menghapus history pengguna.
- Terapkan guard deaktivasi bersama domainnya: PIC task terbuka pada M2, review pending pada M3, PIC event mendatang/berlangsung pada M4. Uji race dengan assignment; jangan menunda guard task sampai M3.
- Notifikasi inbox disimpan meskipun push gagal. Penerima mengikuti assignment aktual, bukan jabatan.
- Push opt-in, gunakan VAPID dengan private key hanya di backend. Validasi subscription endpoint dan lindungi pengiriman dari SSRF.
- Retry push memakai state persisten, lease, dan batas percobaan; disable subscription invalid. Jangan menjanjikan exactly-once pada perangkat.
- Gunakan satu service worker untuk cache asset publik dan push. Jangan cache API/halaman privat atau mengantre perubahan offline.
- Logout membersihkan state klien dan menonaktifkan subscription terkait sesuai PRD. Payload lockscreen tidak memuat detail pekerjaan sensitif.

## UI & Copy

- User-facing copy memakai Bahasa Indonesia, waktu jelas WITA, light mode, sidebar desktop, dan navigasi mobile sesuai PRD.
- Tampilan sederhana, tenang, accessible; hindari animasi berlebihan, glow dekoratif, dan dashboard analytics yang tidak diperlukan.
- Status/priority tidak hanya dibedakan warna. Sediakan label, keyboard access, fokus terlihat, dan alternatif tombol untuk drag-and-drop.
- Gunakan primitive aksesibilitas matang dari foundation yang dipilih untuk dialog, menu, dan combobox; jangan hand-roll tanpa alasan.
- Tampilkan `Tandai Selesai`, `Minta Review`, `Setujui`, dan `Minta Revisi` hanya pada konteks yang sah. Respons backend tetap menentukan hasil akhir.
- Jangan menampilkan placeholder, data demo, atau toast sukses seolah operasi production sudah berhasil.

## Database, Dependencies & Deployment

- Perubahan schema menggunakan versioned SQL migration. Jangan mengedit migration yang sudah diterapkan; buat migration baru.
- Ikuti urutan schema PRD §13.5. Completion CHECK dasar tersedia sejak M2; kolom/FK review, event, dan reminder ditambahkan saat tabel tujuannya tersedia, dengan uji upgrade antarmilestone.
- Pertahankan FK, constraints, indexes yang relevan, audit append-only, dan archive semantics. Jangan hard-delete pekerjaan/user untuk memudahkan implementasi.
- Runtime memakai pooled database connection dengan pool kecil reusable; migration/backup memakai direct connection. Jangan migrasi pada cold start atau build frontend.
- Dev, preview, dan production wajib terpisah; preview memakai data sintetis dan tidak memicu cron/push production.
- Target deployment tetap dua Vercel Projects. Buktikan Go routing, proxy Cookie/Set-Cookie, environment pairing, dan runtime limits sebelum mengklaim siap production.
- Jangan mengasumsikan hosting bisnis atau cron gratis. Verifikasi paket dan biaya ketika deployment menjadi scope tugas.
- Gunakan pnpm untuk frontend dan Go Modules untuk backend. Tambah dependency melalui package manager, periksa kompatibilitas, dan sertakan lockfile/go.sum.
- Jangan hardcode secret, credential, production origin, atau API URL environment. Perbarui `.env.example` dengan placeholder saja.
- Jangan menjalankan migration destruktif, restore database, atau perubahan akses production di luar scope yang diotorisasi.

## Quality Gates

- Baca scripts dalam package.json, go.mod, dan workflow CI sebelum memilih command. Jangan mengarang nama script atau mengklaim tool tersedia.
- Saat setup awal, sediakan dan dokumentasikan command quality frontend untuk format/lint, Svelte/TypeScript checks, test, dan build. `pnpm quality` boleh dipakai setelah benar-benar didefinisikan.
- Untuk perubahan frontend, jalankan checks dan test relevan melalui script repository. Gunakan Svelte autofixer jika tool tersedia; jika tidak, gunakan compiler/typecheck dan sebutkan keterbatasannya.
- Untuk backend, jalankan pemeriksaan gofmt, `go vet ./...`, dan `go test ./...` dari `apps/api`; tambahkan integration tests ketika perubahan menyentuh transaksi/schema.
- Perubahan kontrak memerlukan validasi OpenAPI dan regenerasi client. Perubahan migration diuji pada PostgreSQL, bukan hanya mock/SQLite.
- Prioritaskan test permission, direct completion, review dua arah, self-approval rejection, stale version, idempotency, WITA, serta reminder reschedule/retry.
- Test harus sebanding dengan risiko. Perubahan copy/style kecil cukup pemeriksaan visual/build relevan; jangan menambah test yang hanya menyalin implementasi.
- Jalankan checks penuh yang diwajibkan CI sebelum merge/release. Jika tidak dapat dijalankan, nyatakan command, penyebab, dan risiko yang belum terverifikasi; jangan menyebut semuanya lulus.

## Documentation & Handoff

- Perbarui PRD/OpenAPI/README/runbook ketika mengubah requirement, route, env, auth flow, schema, deployment, atau architecture boundary.
- Jangan mengubah PRD untuk membenarkan pelanggaran implementasi. Catat perubahan produk yang diminta secara eksplisit.
- Pertahankan perubahan pengguna yang tidak terkait; hindari refactor besar, formatting massal, atau file hasil generator yang tidak diperlukan.
- Jika diminta membuat commit, gunakan Conventional Commits, misalnya `feat(tasks): add optional review` atau `fix(reminders): prevent duplicate notifications`.
- Saat menyerahkan hasil, jelaskan perilaku yang berubah, cara validasi, dan keterbatasan yang tersisa. Bedakan implementasi lokal, preview, dan production deployment.
