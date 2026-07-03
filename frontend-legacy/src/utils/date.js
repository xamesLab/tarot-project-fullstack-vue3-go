export function formatDateTime(value, locale = 'en') {
  if (!value) {
    return ''
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return ''
  }

  return new Intl.DateTimeFormat(locale === 'ru' ? 'ru-RU' : 'en-US', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(date)
}

export function formatDateOnly(value, locale = 'en') {
  if (!value) {
    return ''
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }

  return new Intl.DateTimeFormat(locale === 'ru' ? 'ru-RU' : 'en-US', {
    dateStyle: 'medium',
  }).format(date)
}

export function getAccountingDateForTimezone(timezone, date = new Date()) {
  const zoned = new Date(date.toLocaleString('en-US', { timeZone: timezone || 'UTC' }))

  if (zoned.getHours() < 3) {
    zoned.setDate(zoned.getDate() - 1)
  }

  const year = zoned.getFullYear()
  const month = String(zoned.getMonth() + 1).padStart(2, '0')
  const day = String(zoned.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}
