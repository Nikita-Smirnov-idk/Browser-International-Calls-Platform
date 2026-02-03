import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { useLocale } from '../i18n/LocaleContext'
import styles from './Layout.module.css'

interface LayoutProps {
  children: React.ReactNode
}

export function Layout({ children }: LayoutProps) {
  const { isAuthenticated, logout } = useAuth()
  const { t, locale, setLocale } = useLocale()
  const navigate = useNavigate()

  const handleLogout = async () => {
    await logout()
    navigate('/')
  }

  return (
    <div className={styles.wrapper}>
      <header className={styles.header}>
        <Link to="/" className={styles.logo}>
          {t.appName}
        </Link>
        <nav className={styles.nav}>
          {isAuthenticated ? (
            <>
              <Link to="/call">{t.call}</Link>
              <Link to="/history">{t.history}</Link>
              <button onClick={handleLogout} className={styles.logout}>
                {t.logout}
              </button>
            </>
          ) : (
            <>
              <Link to="/login">{t.login}</Link>
              <Link to="/register">{t.register}</Link>
            </>
          )}
          <button
            className={styles.lang}
            onClick={() => setLocale(locale === 'ru' ? 'en' : 'ru')}
          >
            {locale === 'ru' ? 'EN' : 'RU'}
          </button>
        </nav>
      </header>
      <main className={styles.main}>{children}</main>
    </div>
  )
}
