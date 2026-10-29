# WORKFLOW — Aturan Main AI

Versi 1.0 · 29 Oktober 2026 · Unggulmart Marketing Workspace

Berlaku untuk AI coding assistant yang bekerja pada repository ini, termasuk ketika dokumen dibaca melalui Antigravity atau tool lain. Tempatkan `WORKFLOW.md` di root repository. Dokumen ini mengatur otonomi, komunikasi, izin, verifikasi, dan serah terima; bukan bukti bahwa aplikasi sudah diimplementasikan.

## 1. Prinsip kerja

Kerjakan permintaan sampai hasilnya konkret dan dapat diperiksa. Ambil keputusan implementasi rutin sendiri selama sesuai scope, dapat dipulihkan, dan tidak mengubah requirement. Tanyakan hanya hal yang benar-benar memerlukan keputusan pengguna atau otorisasi tambahan.

Jangan berhenti setelah membuat rencana jika pengguna meminta implementasi. Jangan memperluas tugas menjadi pembangunan seluruh aplikasi. Jangan menyamakan akses tool, credential yang tersedia, atau tombol deploy dengan izin menggunakan semuanya.

Instruksi platform dan batas keamanan tetap berlaku. Instruksi pengguna yang eksplisit menentukan scope; permintaan baru yang mengubah keputusan produk harus dicatat pada dokumen terkait. File ini tidak memberi hak untuk melewati pembatasan tool atau menyetujui tindakan sendiri.

## 2. Peran enam dokumen

| Dokumen | Menjawab pertanyaan |
|---|---|
| `PRD.md` | Apa yang dibangun, aturan bisnisnya, non-goals, dan acceptance criteria |
| `AGENTS.md` | Batas teknis dan aturan agent di repository |
| `ARCHITECTURE.md` | Di mana kode ditempatkan dan siapa memiliki tanggung jawab |
| `TODO.md` | Task/dependensi apa yang aktif, status, bukti, dan gate milestone |
| `SKILL.md` | Bagaimana menjalankan workflow berulang yang dipilih |
| `WORKFLOW.md` | Kapan bertindak, meminta keputusan/izin, berhenti, dan menyatakan selesai |

Seluruh dokumen Markdown proyek berada di root repository sesuai keputusan pengguna, termasuk SKILL.md, ADR (architecture-decisions.md), deployment.md, runbook.md, dan CALENDAR-DATA.md saat dibuat. Kontrak API berada di openapi.yaml; asset dataset di calendar-data/. Kode dan migration tetap di apps/ dan db/migrations/. Gunakan nama/kapitalisasi aktual tanpa membuat duplikat atau memindahkan dokumen ke folder khusus tanpa permintaan pengguna.

Nama `WORKFLOW.md` tidak menjamin AI tool otomatis memuatnya. Pada sesi baru, minta agent membaca file ini atau hubungkan melalui mekanisme instruksi tool yang tersedia. Jangan mengklaim pemasangan atau auto-discovery yang belum diverifikasi.

Jika dokumen bertentangan, gunakan masing-masing sesuai perannya. Catat kutipan/lokasi konflik dan dampak; jangan memilih aturan yang paling memudahkan coding. Lanjutkan bagian yang tidak terdampak, lalu minta keputusan untuk konflik material. TODO bukan izin menghapus acceptance; SKILL bukan izin deploy; architecture tidak mengganti requirement produk.

## 3. Mulai sesi dengan konteks aktual

1. Temukan repository dan branch/worktree yang dimaksud. Periksa status perubahan sebelum mengedit; pertahankan pekerjaan pengguna yang tidak terkait.
2. Baca AGENTS beserta instruksi lokal, WORKFLOW, task TODO terkait, lalu bagian PRD/architecture yang relevan dan workflow SKILL yang dibutuhkan. Pada sesi pertama, pahami cakupan seluruh dokumen sebelum memilih task.
3. Periksa implementasi, kontrak OpenAPI, migrations, scripts, lockfile dan CI yang terkait. Dokumen target tidak membuktikan file atau command sudah tersedia.
4. Identifikasi hasil yang diminta: diskusi, review saja, implementasi, bug fix, refactor, atau release. Jangan mengedit kode saat hanya diminta menjelaskan atau mereview.
5. Sampaikan singkat task, batas scope dan pemeriksaan yang akan dilakukan. Rencana rutin tidak memerlukan persetujuan tersendiri.

Jika pengguna mengatakan “lanjut” tanpa ID, gunakan konteks sesi. Jika konteks belum menentukan task, pilih task siap pertama pada milestone aktif dan sebutkan pilihannya. Jangan menanyakan pilihan yang sudah jelas dari TODO. Jangan menganggap “lanjut” sebagai izin mengerjakan seluruh milestone berikutnya atau merilis produksi.

Jika repository atau dokumen wajib tidak tersedia, laporkan kebutuhan yang spesifik. Kerjakan pemeriksaan yang tidak bergantung pada bagian yang hilang; jangan mengarang path, acceptance, API, atau hasil sesi lama.

## 4. Matriks tindakan dan izin

Izin yang sudah diberikan tetap berlaku untuk target, environment, scope dan dampak yang sama. Jangan meminta ulang hanya karena sesi berlanjut. Perubahan target, biaya, audience, atau dampak material membutuhkan penilaian baru.

| Kondisi / tindakan | Keputusan AI |
|---|---|
| Membaca kode/docs, memeriksa diff/log aman, mencari sumber teknis resmi | Jalan sendiri dalam scope; jangan keluarkan secrets |
| Memilih nama helper, membagi komponen/service, memperbaiki typing atau validasi sesuai requirement | Jalan sendiri; ikuti pola repository |
| Implementasi fitur yang diminta, test, perbaikan akibat perubahan sendiri, docs/kontrak/client terkait | Jalan sendiri sampai acceptance dan checks terpenuhi |
| Memasang dependency yang diwajibkan stack atau kebutuhan tugasnya jelas | Jalan sendiri melalui package manager, cek kompatibilitas/lockfile; jangan menambah layanan berbayar atau infrastruktur terselubung |
| Menulis migration baru dan mengujinya di DB disposable/dev yang targetnya sudah dipastikan | Jalan sendiri; tidak otomatis berarti boleh menjalankan pada DB bersama/produksi |
| Mengisi status/bukti TODO berdasarkan pekerjaan aktual | Jalan sendiri; tidak mengubah scope atau menyatakan UAT pengguna tanpa bukti |
| Asumsi kecil, mudah dibalik, tidak memengaruhi permission/data/API publik | Pilih default PRD/pola existing, catat jika berpengaruh pada hasil |
| Ambiguitas mengubah business rule, permission, tanggal, retensi, data loss, atau acceptance | Klarifikasi keputusan; lanjutkan bagian independen |
| Menambah non-goal, mengganti stack/host/auth, dependency besar tanpa kebutuhan jelas, atau melompati gate | Siapkan bukti dan usulan; minta keputusan eksplisit sebelum perubahan material |
| Reset DB, truncate/drop data, restore, destructive migration, hapus file pengguna, rewrite history atau force push | Perlu otorisasi spesifik bila belum ada; pastikan target dan dampak terlebih dahulu |
| Deploy preview/staging terisolasi yang diminta atau menjadi scope eksplisit deployment spike | Jalan sesuai otorisasi, tanpa produksi dan tanpa biaya/audience baru yang belum disetujui |
| Deploy production, menjalankan migration produksi, mengaktifkan cron/push nyata, mengganti secrets/akses | Jalankan hanya jika tindakan dan environment itu sudah diotorisasi serta gate terpenuhi; jika belum, minta izin setelah persiapan lengkap |
| Membeli paket, mengubah billing, membuat resource berbayar, membuka aplikasi/data ke audience baru | Minta izin eksplisit dengan dampak/biaya yang diketahui |
| Mengirim pesan/email, membuat undangan, atau mengirim notifikasi uji ke pengguna nyata | Hanya jika pengiriman/penerimanya diotorisasi; default pengujian memakai fixture/fake transport |
| Commit, push, membuat PR, merge, atau mengubah pengaturan remote | Ikuti permintaan/otorisasi yang ada; implementasi lokal saja bukan otomatis izin publikasi atau merge |
| Tool menolak tindakan karena izin/kebijakan | Jangan bypass; pilih alternatif yang aman atau laporkan blokir beserta alasannya |

Melihat credential tidak memberi izin membaca data lain. Konfirmasi environment melalui konfigurasi yang tidak membocorkan secret. Jika target database tidak dapat dipastikan nonproduction, jangan menjalankan mutasi dengan asumsi “ini pasti dev”.

## 5. Cara meminta klarifikasi atau izin

Bedakan dua kebutuhan: **klarifikasi** menentukan hasil yang diinginkan; **izin** mengotorisasi tindakan yang dampaknya sudah diketahui. Jangan bertanya “boleh lanjut?” tanpa menjelaskan keputusan yang dibutuhkan.

Sebelum meminta izin tindakan berisiko, selesaikan persiapan yang aman dan diotorisasi: implementasi lokal, diff, test, migration plan, target environment, serta rencana recovery yang relevan. Jangan menjalankan tindakannya sebelum persetujuan.

Permintaan izin harus singkat dan konkret:

- Tindakan dan target tepatnya, termasuk environment/audience.
- Mengapa diperlukan dan hasil persiapan yang sudah tersedia.
- Dampak data, akses, biaya atau downtime yang relevan.
- Cara rollback/pemulihan atau keterbatasannya.
- Keputusan spesifik yang harus diberikan pengguna.

Contoh: “Migration dan tests staging sudah lulus. Langkah berikutnya menjalankan migration X pada database production Y sebelum deploy API. Perubahan ini menambah kolom nullable; rollback aplikasi tetap kompatibel. Apakah boleh menjalankan migration dan deploy versi Z ke production?”

Jika permintaan disebabkan aturan file atau penolakan tool, sebutkan sumber dan alasan sebenarnya. Jangan menciptakan approval gate untuk pilihan coding rutin. Diam, timeout, dan jawaban ambigu bukan persetujuan. Jangan meminta ulang izin yang sudah jelas; minta kembali hanya saat kondisi/dampaknya berubah.

## 6. Siklus eksekusi task

1. **Tentukan scope.** Resolve task ID, dependensi, acceptance dan milestone. Untuk permintaan ad hoc, tautkan atau tambahkan task/BUG yang sesuai tanpa mengubah requirement diam-diam.
2. **Mulai pekerjaan.** Isi PIC pelaksana aktual dan IN_PROGRESS pada task yang dikerjakan. Milestone yang dependensinya belum selesai tidak otomatis BLOCKED; catat blocker hanya jika konkret.
3. **Implementasikan satu alur lengkap.** Kontrak → migration bila perlu → policy/service/transaksi → handler → UI. Jangan membuat layar seolah berhasil sebelum backend benar-benar menyimpan.
4. **Verifikasi sesuai risiko.** Jalankan checks nyata dan acceptance, perbaiki kegagalan dalam scope, lalu ulangi pemeriksaan terdampak.
5. **Periksa diff akhir.** Pastikan tidak ada secrets, perubahan pengguna terhapus, generated drift, dependency tak perlu, atau perubahan di luar scope.
6. **Perbarui docs dan TODO.** Catat status, bukti, issue/gate, hitung ulang progres, dan langkah berikutnya.
7. **Serahkan hasil.** Jelaskan behavior akhir, bukti dan keterbatasan. Jangan menyebut “selesai” untuk tahap yang belum dibuktikan.

Default unit kerja adalah task yang diminta beserta pekerjaan pendukung yang diperlukan. Jika diminta menyelesaikan milestone, lanjutkan seluruh task dependennya sampai gate terpenuhi atau ada blocker konkret. Jika diminta membangun seluruh scope, tetap lewati gate berurutan. Gate teknis yang dapat dibuktikan AI tidak perlu persetujuan manusia tambahan; gate UAT/biaya/akses yang membutuhkan keputusan nyata tidak boleh diisi sendiri.

## 7. Verifikasi yang wajib sebanding dengan perubahan

Periksa command pada package.json, go.mod dan CI terlebih dahulu. `pnpm quality` hanya digunakan setelah didefinisikan. Gunakan SKILL verify-change dan AGENTS sebagai panduan rinci.

| Dampak perubahan | Bukti minimum yang relevan |
|---|---|
| Copy/style ringan | Pemeriksaan visual dan checks/build yang diperlukan; tidak perlu test cermin implementasi |
| Form/state/UI | Svelte/TypeScript checks, behavior loading/empty/error/success/conflict, akses keyboard dan build terkait |
| Go domain | Gofmt, go vet, Go tests sesuai AGENTS; policy dan transisi positif/negatif |
| API | OpenAPI valid, generated client konsisten, status/error/schema dan unknown-field handling sesuai kontrak |
| DB/transaksi | Migration pada PostgreSQL nyata; constraints, rollback, version/idempotency dan race terkait |
| Auth/review | Permission melalui API, self-approval/stale round/session revoke dan race yang terdampak |
| Kalender/reminder | WITA/date-only/all-day, revision/dedup/catch-up/reschedule; tidak hanya screenshot kalender |
| PWA/push | Cache privacy, opt-in/logout/retry/lease dan deep link; uji perangkat jika diwajibkan acceptance |
| Merge/release | Checks penuh wajib CI dan gate environment/release terkait |

Kegagalan tidak boleh diselesaikan dengan menghapus assertion bermakna, mematikan security check, mengganti requirement, atau mengklaim error existing tanpa bukti. Jangan menjalankan ulang suite yang sama tanpa perubahan/alasan jika sudah lulus. Jika tool/perangkat tidak tersedia, gunakan pemeriksaan alternatif yang valid dan laporkan batasnya; jangan menganggap alternatif otomatis menggantikan gate wajib.

Laporkan hasil sebagai PASS, FAIL atau NOT RUN. NOT RUN tidak sama dengan PASS. Pisahkan hasil lokal, CI, preview/staging, production, dan penerimaan pengguna.

## 8. Kapan dianggap selesai

### Task selesai

Task boleh DONE jika seluruh kondisi yang berlaku terpenuhi:

- Acceptance baris TODO dan requirement PRD terkait terbukti.
- Behavior berjalan pada jalur yang diminta, bukan placeholder atau hanya UI.
- Permission, validasi dan transisi berada di Go; tidak bergantung pada tombol disabled.
- Kontrak, schema, generated client, audit/inbox dan concurrency semantics konsisten bila terdampak.
- Checks wajib dan pengujian sesuai risiko lulus dengan bukti aktual.
- Dokumentasi/env example/runbook yang terdampak diperbarui.
- Diff akhir bersih dari perubahan tak terkait, kebocoran secret, dan regresi diketahui yang menggagalkan acceptance.
- TODO, ringkasan progres serta catatan bukti diperbarui.

Task lokal tidak perlu menunggu deploy jika acceptance-nya memang lokal. Task deployment tidak selesai hanya karena build lokal lulus. Dokumentasi atau rencana selesai bukan berarti implementasi selesai.

### Fitur dan milestone selesai

Fitur selesai jika seluruh task pembentuk alur fiturnya DONE dan acceptance lintas layer yang relevan lulus. Jangan menyebut fitur selesai jika API baru ada tetapi UI yang diwajibkan belum bekerja.

Milestone selesai jika seluruh task aktif dan exit gate PRD §24 terpenuhi dengan bukti. Jangan melompati M1–M7 atau mengisi gate karena persentase tracker 100%. Scope yang dibatalkan membutuhkan keputusan eksplisit dan tidak otomatis membuktikan gate.

Pertanyaan terbuka PRD dilacak melalui issue keputusan di TODO. Default memungkinkan implementasi, tetapi task/gate yang mencapai tenggat keputusan belum boleh DONE tanpa bukti konfirmasi yang diperlukan. Perubahan dokumentasi tidak dihitung sebagai penerimaan default oleh pengguna.

### Siap production dan sudah dirilis

“Siap production” memerlukan DoD PRD §25 serta gate operasional: CI, environment isolation, routing/cookie proxy, biaya/paket, kalender, keamanan, backup/restore, rollback, monitoring dan UAT yang relevan.

“Sudah dirilis” memerlukan otorisasi rilis, deployment berhasil pada environment yang disebut, dan smoke aktual. Penerimaan Sultan/Alvin hanya boleh dicatat dari konfirmasi/pelaksanaan nyata. Jangan mengubah hasil tes AI menjadi persetujuan pengguna.

## 9. Status tracker dan bukti

| Status TODO | Gunakan ketika |
|---|---|
| TODO | Belum mulai atau masih antre dependensi |
| IN_PROGRESS | Sedang dikerjakan oleh PIC aktual |
| VERIFY | Implementasi tersedia tetapi acceptance/checks wajib belum lengkap |
| BLOCKED | Ada hambatan konkret yang mencegah kemajuan task; catat owner dan pembuka blokir |
| DONE | Seluruh syarat selesai task terbukti |
| CANCELLED | Ada keputusan eksplisit membatalkan scope, alasan dan referensinya tercatat |

Status ini adalah status pekerjaan pengembangan, bukan state task produk. Pertahankan aturan hitung TODO: hanya leaf task aktif, gate tidak masuk denominator, tanpa double count. Catat tanggal WITA, file/commit bila tersedia, command/metode aktual dan hasil. Jangan mengarang commit, test, link preview, tanggal selesai, atau persentase.

Jika bug membatalkan acceptance task DONE, buka kembali task dan gate yang terdampak; tautkan BUG, hitung ulang progres serta nilai dampak ke milestone lanjutan. Jangan menyembunyikannya dengan mempertahankan angka 100%.

## 10. Hambatan, kegagalan, dan interupsi

Untuk blocker, catat expected/actual, bukti/log aman, penyebab yang diketahui versus hipotesis, upaya relevan dan tindakan pembuka. Jika butuh pengguna, minta informasi minimum yang tidak dapat ditemukan sendiri.

Jangan mengulang aksi gagal secara buta, menaikkan hak akses untuk melewati penolakan, mencoba ulang pengiriman dengan key baru, atau melakukan rollback luas pada file pengguna. Setelah upaya masuk akal tidak memberi bukti baru, ubah strategi atau laporkan blocker. Tetap lanjutkan bagian independen yang aman.

Jika ditemukan kegagalan existing di luar scope, dokumentasikan terpisah. Perbaiki hanya jika diperlukan untuk acceptance dan perubahan masih proporsional; bila memperluas scope secara material, minta keputusan. Jangan menyatakan semua checks lulus hanya karena error berasal dari area lain.

Jika pengguna menyela atau mengoreksi, tanggapi singkat dan masukkan arahan ke task aktif. Berhenti bila diminta berhenti. Saat sesi terputus, simpan handoff minimal di TODO: task/status, perubahan yang sudah ada, verifikasi, blocker, dan next action. Lanjutan sesi memeriksa state aktual; tidak langsung mengulang aksi eksternal atau migration.

## 11. Komunikasi dan laporan akhir

Gunakan Bahasa Indonesia yang jelas. Sebutkan hasil/keputusan lebih dulu, lalu alasan dan bukti yang diperlukan. Saat pekerjaan panjang, kirim update ketika ada temuan, perubahan pendekatan, blocker, atau hasil verifikasi; jangan membiarkan pengguna menunggu lama tanpa konteks. Jangan membanjiri pengguna dengan setiap command.

Laporan implementasi memuat perubahan behavior, lokasi penting, verifikasi dan risiko yang masih relevan. Gunakan tabel ketika melaporkan progres:

| Task ID | Status sebelum → sesudah | Bukti / hasil pemeriksaan | Langkah berikut / blocker |
|---|---|---|---|

Jika ada permintaan izin yang tertunda, laporkan persiapan yang sudah selesai dan tindakan tepat yang menunggu izin. Jangan menyembunyikan pekerjaan belum selesai di balik kalimat “berhasil”. Jangan menawarkan kelanjutan untuk pekerjaan yang sebenarnya sudah diotorisasi dan masih dapat dikerjakan.

## 12. Contoh keputusan sehari-hari

| Situasi | Respons yang diharapkan |
|---|---|
| “Kerjakan M2-03”; dependensi/gate terpenuhi | Implementasikan, uji, update TODO; tidak perlu persetujuan nama helper atau struktur form |
| “Review security login” | Review dan laporkan temuan; jangan mengubah auth atau mengirim serangan uji ke production |
| Tests menunjukkan self-approval lolos | Perbaiki jika scope fix diotorisasi, uji negatif/race; jangan menghapus larangan self-approval |
| M4 UI selesai tetapi perangkat UTC memperlihatkan tanggal salah | Task terkait tetap VERIFY/IN_PROGRESS; fitur kalender belum selesai |
| Build lulus tetapi PostgreSQL integration belum dapat dijalankan | Laporkan NOT RUN dan kendala; task yang membutuhkan bukti transaksi belum DONE |
| M3 belum lulus tetapi muncul ide upload file | Catat sebagai usulan di luar scope; jangan implementasikan atau mengubah PRD sendiri |
| “Siapkan release” tanpa perintah deploy | Selesaikan readiness, diff, tests dan rencana; minta izin hanya untuk eksekusi yang belum diotorisasi |
| “Deploy versi X ke production Y” sudah eksplisit dan gate lengkap | Jalankan dalam cakupan izin itu, verifikasi deployment/smoke; jangan meminta konfirmasi identik lagi |
| Seluruh file spesifikasi sudah dibuat, belum ada kode | Dokumen selesai; progres implementasi tidak berubah menjadi DONE |

## 13. Aktivasi untuk sesi baru

Gunakan instruksi pembuka berikut di AI coding tool yang dipakai, dengan menyesuaikan lokasi file aktual:

> Baca AGENTS.md dan WORKFLOW.md terlebih dahulu, lalu PRD, architecture, TODO, dan workflow relevan dalam SKILL.md. Periksa state repository aktual. Kerjakan task [ID] sampai acceptance dan checks wajib terpenuhi, dalam scope serta batas izin dokumen. Jalankan keputusan rutin sendiri; ajukan hanya keputusan material atau tindakan yang belum diotorisasi. Perbarui TODO dengan bukti aktual dan laporkan hasilnya. Jangan deploy atau mengubah production kecuali itu termasuk otorisasi tugas.

Instruksi ini adalah cara meminta agent membaca dokumen; bukan konfigurasi integrasi otomatis atau perintah untuk mulai coding ketika pengguna hanya meminta dokumen.
