import { readFileSync } from 'node:fs';
import { execFileSync } from 'node:child_process';

const root = new URL('../', import.meta.url);
const expectedNode = readFileSync(new URL('.node-version', root), 'utf8').trim();
const expectedGo = readFileSync(new URL('.go-version', root), 'utf8').trim();

if (process.versions.node !== expectedNode) {
  console.error(`Gunakan Node.js ${expectedNode}; runtime aktif ${process.versions.node}.`);
  process.exit(1);
}

try {
  const actualGo = execFileSync('go', ['version'], { encoding: 'utf8' }).trim();
  if (actualGo.split(' ')[2] !== `go${expectedGo}`) {
    throw new Error(`Gunakan Go ${expectedGo}; runtime aktif: ${actualGo}.`);
  }
  console.log(`Node.js ${expectedNode}; ${actualGo}`);
} catch (error) {
  console.error(error.message);
  console.error('Pastikan Go tersedia di PATH sesuai README.md.');
  process.exit(1);
}
