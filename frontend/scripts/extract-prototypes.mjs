import { mkdir, readFile, writeFile } from 'node:fs/promises';

const sources = [
  ['Arcana Home.html', 'home'],
  ['Arcana UI Library.html', 'ui-kit'],
];

await mkdir('reference', { recursive: true });
for (const [source, output] of sources) {
  const bundle = await readFile(source, 'utf8');
  const match = bundle.match(/<script type="__bundler\/template">(.*?)<\/script>/s);
  if (!match) throw new Error(`Template not found in ${source}`);
  const template = JSON.parse(match[1]);
  await writeFile(`reference/${output}.html`, template, 'utf8');
}

console.log('Prototype templates extracted to reference/.');
