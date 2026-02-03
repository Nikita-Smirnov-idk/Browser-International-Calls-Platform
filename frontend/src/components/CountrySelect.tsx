import { useRef, useState, useCallback, useEffect } from 'react'
import type { Country } from '../data/countries'
import styles from './CountrySelect.module.css'

interface CountrySelectProps {
  countries: Country[]
  value: Country
  onChange: (country: Country) => void
  locale: 'ru' | 'en'
  disabled?: boolean
}

export function CountrySelect({
  countries,
  value,
  onChange,
  locale,
  disabled,
}: CountrySelectProps) {
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState('')
  const listRef = useRef<HTMLUListElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  const name = locale === 'ru' ? value.nameRu : value.nameEn
  const filtered = query.trim()
    ? countries.filter((c) => {
        const n = locale === 'ru' ? c.nameRu : c.nameEn
        const code = c.code.toLowerCase()
        const dial = c.dialCode
        const q = query.toLowerCase()
        return n.toLowerCase().includes(q) || code.includes(q) || dial.includes(q)
      })
    : countries

  const select = useCallback(
    (c: Country) => {
      onChange(c)
      setOpen(false)
      setQuery('')
    },
    [onChange]
  )

  useEffect(() => {
    if (!open) return
    const handle = (e: MouseEvent) => {
      const el = (e.target as Node)
      if (listRef.current?.contains(el) || inputRef.current?.contains(el)) return
      setOpen(false)
    }
    document.addEventListener('mousedown', handle)
    return () => document.removeEventListener('mousedown', handle)
  }, [open])

  return (
    <div className={styles.wrap}>
      <button
        type="button"
        className={styles.trigger}
        onClick={() => !disabled && setOpen((o) => !o)}
        disabled={disabled}
        aria-expanded={open}
        aria-haspopup="listbox"
      >
        <span className={styles.flag}>{value.flag}</span>
        <span className={styles.label}>{name} ({value.dialCode})</span>
      </button>
      {open && (
        <div className={styles.dropdown}>
          <input
            ref={inputRef}
            type="text"
            className={styles.search}
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder={locale === 'ru' ? 'Поиск...' : 'Search...'}
            autoFocus
          />
          <ul ref={listRef} className={styles.list} role="listbox">
            {filtered.map((c) => (
              <li
                key={c.code}
                role="option"
                aria-selected={c.code === value.code}
                className={c.code === value.code ? styles.selected : ''}
                onClick={() => select(c)}
              >
                <span className={styles.flag}>{c.flag}</span>
                {locale === 'ru' ? c.nameRu : c.nameEn} ({c.dialCode})
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  )
}
