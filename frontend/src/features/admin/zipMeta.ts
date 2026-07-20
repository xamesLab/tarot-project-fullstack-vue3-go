const EOCD_SIGNATURE = 0x06054b50;
const CENTRAL_DIRECTORY_SIGNATURE = 0x02014b50;
const LOCAL_FILE_HEADER_SIGNATURE = 0x04034b50;
const MAX_EOCD_SEARCH = 0xffff + 22;

const utf8Decoder = new TextDecoder('utf-8');

export interface DeckZipMeta {
  nameEn: string;
  nameRu: string;
  descriptionEn: string;
  descriptionRu: string;
  author: string;
  year: string;
  tags: string;
}

interface ZipEntry {
  name: string;
  compressionMethod: number;
  compressedSize: number;
  localHeaderOffset: number;
}

export async function readDeckMetaFromZip(file: File): Promise<DeckZipMeta | null> {
  const buffer = await file.arrayBuffer();
  const view = new DataView(buffer);
  const eocdOffset = findEndOfCentralDirectory(view);
  if (eocdOffset < 0) return null;

  const centralDirectoryOffset = view.getUint32(eocdOffset + 16, true);
  const totalEntries = view.getUint16(eocdOffset + 10, true);
  const entries = parseCentralDirectory(buffer, view, centralDirectoryOffset, totalEntries);
  const metaEntry = entries.find((entry) => basename(entry.name).toLowerCase() === '0_meta.json');
  if (!metaEntry) return null;

  const entryBytes = await extractZipEntry(buffer, metaEntry);
  const parsed = JSON.parse(utf8Decoder.decode(entryBytes)) as Record<string, unknown>;
  const en = typeof parsed.en === 'object' && parsed.en ? parsed.en as Record<string, unknown> : {};
  const ru = typeof parsed.ru === 'object' && parsed.ru ? parsed.ru as Record<string, unknown> : {};
  return {
    nameEn: normalizeString(en.name),
    nameRu: normalizeString(ru.name),
    descriptionEn: normalizeString(en.text),
    descriptionRu: normalizeString(ru.text),
    author: normalizeString(parsed.author),
    year: typeof parsed.year === 'number' ? String(parsed.year) : normalizeString(parsed.year),
    tags: Array.isArray(parsed.tags)
      ? parsed.tags.map(normalizeString).filter(Boolean).join(', ')
      : normalizeString(parsed.tags),
  };
}

function findEndOfCentralDirectory(view: DataView) {
  const minOffset = Math.max(0, view.byteLength - MAX_EOCD_SEARCH);
  for (let offset = view.byteLength - 22; offset >= minOffset; offset -= 1) {
    if (view.getUint32(offset, true) === EOCD_SIGNATURE) return offset;
  }
  return -1;
}

function parseCentralDirectory(buffer: ArrayBuffer, view: DataView, offset: number, totalEntries: number) {
  const entries: ZipEntry[] = [];
  let cursor = offset;
  for (let index = 0; index < totalEntries; index += 1) {
    if (cursor + 46 > view.byteLength || view.getUint32(cursor, true) !== CENTRAL_DIRECTORY_SIGNATURE) break;
    const compressedSize = view.getUint32(cursor + 20, true);
    const fileNameLength = view.getUint16(cursor + 28, true);
    const extraFieldLength = view.getUint16(cursor + 30, true);
    const commentLength = view.getUint16(cursor + 32, true);
    entries.push({
      name: utf8Decoder.decode(new Uint8Array(buffer, cursor + 46, fileNameLength)),
      compressionMethod: view.getUint16(cursor + 10, true),
      compressedSize,
      localHeaderOffset: view.getUint32(cursor + 42, true),
    });
    cursor += 46 + fileNameLength + extraFieldLength + commentLength;
  }
  return entries;
}

async function extractZipEntry(buffer: ArrayBuffer, entry: ZipEntry) {
  const view = new DataView(buffer);
  if (view.getUint32(entry.localHeaderOffset, true) !== LOCAL_FILE_HEADER_SIGNATURE) throw new Error('Invalid ZIP header');
  const fileNameLength = view.getUint16(entry.localHeaderOffset + 26, true);
  const extraFieldLength = view.getUint16(entry.localHeaderOffset + 28, true);
  const dataStart = entry.localHeaderOffset + 30 + fileNameLength + extraFieldLength;
  const compressed = new Uint8Array(buffer, dataStart, entry.compressedSize);
  if (entry.compressionMethod === 0) return compressed;
  if (entry.compressionMethod === 8 && typeof DecompressionStream !== 'undefined') {
    const stream = new Blob([compressed]).stream().pipeThrough(new DecompressionStream('deflate-raw'));
    return new Uint8Array(await new Response(stream).arrayBuffer());
  }
  throw new Error('Unsupported ZIP compression');
}

function basename(value: string) { return value.split('/').pop() || value; }
function normalizeString(value: unknown) { return typeof value === 'string' ? value.trim() : ''; }
