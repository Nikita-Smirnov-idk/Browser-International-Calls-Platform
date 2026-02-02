import {
  createContext,
  useCallback,
  useContext,
  useState,
  type ReactNode,
} from 'react'
import { translations, type Locale } from './translations'

const LOCALE_KEY = 'locale'
const DEFAULT: Locale = 'ru'

interface LocaleContextValue {
  locale: Locale
  setLocale: (l: Locale) => void
  t: (typeof translations)[Locale]
}

const LocaleContext = createContext<LocaleContextValue | null>(null)

export function LocaleProvider({ children }: { children: ReactNode }) {
  const [locale, setLocaleState] = useState<Locale>(() => {
    const s = localStorage.getItem(LOCALE_KEY)
    return (s === 'en' || s === 'ru' ? s : DEFAULT) as Locale
  })

  const setLocale = useCallback((l: Locale) => {
    setLocaleState(l)
    localStorage.setItem(LOCALE_KEY, l)
  }, [])

  return (
    <LocaleContext.Provider
      value={{
        locale,
        setLocale,
        t: translations[locale],
      }}
    >
      {children}
    </LocaleContext.Provider>
  )
}

export function useLocale() {
  const ctx = useContext(LocaleContext)
  if (!ctx) {
    throw new Error('useLocale must be used within LocaleProvider')
  }
  return ctx
}
