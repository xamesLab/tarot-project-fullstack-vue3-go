const EOCD_SIGNATURE = 0x06054b50
const CENTRAL_DIRECTORY_SIGNATURE = 0x02014b50
const LOCAL_FILE_HEADER_SIGNATURE = 0x04034b50
const MAX_EOCD_SEARCH = 0xffff + 22

const utf8Decoder = new TextDecoder('utf-8')

export async function readDeckMetaFromZip(file) {
  if (!file) {
    return null
  }

  const buffer = await file.arrayBuffer()
  const view = new DataView(buffer)
  const eocdOffset = findEndOfCentralDirectory(view)
  if (eocdOffset < 0) {
    return null
  }

  const centralDirectoryOffset = view.getUint32(eocdOffset + 16, true)
  const totalEntries = view.getUint16(eocdOffset + 10, true)
  const entries = parseCentralDirectory(buffer, view, centralDirectoryOffset, totalEntries)
  const metaEntry = entries.find((entry) => basename(entry.name).toLowerCase() === '0_meta.json')
  if (!metaEntry) {
    return null
  }

  const entryBytes = await extractZipEntry(buffer, metaEntry)
  const parsed = JSON.parse(utf8Decoder.decode(entryBytes))

  return {
    nameEn: normalizeString(parsed?.en?.name),
    nameRu: normalizeString(parsed?.ru?.name),
    descriptionEn: normalizeString(parsed?.en?.text),
    descriptionRu: normalizeString(parsed?.ru?.text),
    author: normalizeString(parsed?.author),
    year: normalizeScalar(parsed?.year),
    tags: normalizeTags(parsed?.tags),
  }
}

function findEndOfCentralDirectory(view) {
  const minOffset = Math.max(0, view.byteLength - MAX_EOCD_SEARCH)
  for (let offset = view.byteLength - 22; offset >= minOffset; offset -= 1) {
    if (view.getUint32(offset, true) === EOCD_SIGNATURE) {
      return offset
    }
  }

  return -1
}

function parseCentralDirectory(buffer, view, offset, totalEntries) {
  const entries = []
  let cursor = offset

  for (let index = 0; index < totalEntries; index += 1) {
    if (view.getUint32(cursor, true) !== CENTRAL_DIRECTORY_SIGNATURE) {
      break
    }

    const compressionMethod = view.getUint16(cursor + 10, true)
    const compressedSize = view.getUint32(cursor + 20, true)
    const fileNameLength = view.getUint16(cursor + 28, true)
    const extraFieldLength = view.getUint16(cursor + 30, true)
    const commentLength = view.getUint16(cursor + 32, true)
    const localHeaderOffset = view.getUint32(cursor + 42, true)
    const fileNameBytes = new Uint8Array(buffer, cursor + 46, fileNameLength)
    const name = utf8Decoder.decode(fileNameBytes)

    entries.push({
      name,
      compressionMethod,
      compressedSize,
      localHeaderOffset,
    })

    cursor += 46 + fileNameLength + extraFieldLength + commentLength
  }

  return entries
}

async function extractZipEntry(buffer, entry) {
  const view = new DataView(buffer)
  const localOffset = entry.localHeaderOffset
  if (view.getUint32(localOffset, true) !== LOCAL_FILE_HEADER_SIGNATURE) {
    throw new Error('Invalid local file header')
  }

  const fileNameLength = view.getUint16(localOffset + 26, true)
  const extraFieldLength = view.getUint16(localOffset + 28, true)
  const dataStart = localOffset + 30 + fileNameLength + extraFieldLength
  const compressed = new Uint8Array(buffer, dataStart, entry.compressedSize)

  if (entry.compressionMethod === 0) {
    return compressed
  }

  if (entry.compressionMethod === 8) {
    if (typeof DecompressionStream === 'undefined') {
      throw new Error('DecompressionStream is not available')
    }

    const stream = new Blob([compressed]).stream().pipeThrough(new DecompressionStream('deflate-raw'))
    return new Uint8Array(await new Response(stream).arrayBuffer())
  }

  throw new Error(`Unsupported ZIP compression method: ${entry.compressionMethod}`)
}

function basename(value) {
  return value.split('/').pop() || value
}

function normalizeString(value) {
  return typeof value === 'string' ? value.trim() : ''
}

function normalizeScalar(value) {
  if (typeof value === 'number') {
    return String(value)
  }

  return normalizeString(value)
}

function normalizeTags(value) {
  if (Array.isArray(value)) {
    return value.map((item) => normalizeString(item)).filter(Boolean).join(', ')
  }

  return normalizeString(value)
}
