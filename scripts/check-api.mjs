import { execFileSync } from 'node:child_process';
import { fileURLToPath } from 'node:url';

const cwd = fileURLToPath(new URL('../apps/api/', import.meta.url));

try {
  const unformatted = execFileSync('gofmt', ['-l', '.'], { cwd, encoding: 'utf8' });
  if (unformatted.trim()) {
    console.error('Jalankan pnpm format:api untuk file berikut:\n' + unformatted);
    process.exit(1);
  }
  execFileSync('go', ['vet', './...'], { cwd, stdio: 'inherit' });
  console.log('PASS: gofmt dan go vet ./...');
} catch (error) {
  console.error(error.message);
  process.exit(1);
}
