import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'
import { useLocale } from '../i18n/LocaleContext'
import styles from './Auth.module.css'

export function Login() {
  const { login, error, clearError } = useAuth()
  const { t } = useLocale()
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    clearError()
    setLoading(true)
    try {
      await login(email, password)
      navigate('/call')
    } catch {
      setLoading(false)
    }
  }

  return (
    <div className={styles.page}>
      <form className={styles.form} onSubmit={handleSubmit}>
        <h1>{t.login}</h1>
        {error && (
          <div className={styles.error}>
            {error}
            <button type="button" onClick={clearError} aria-label={t.close}>
              ×
            </button>
          </div>
        )}
        <label>
          {t.email}
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            autoComplete="email"
          />
        </label>
        <label>
          {t.password}
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            autoComplete="current-password"
          />
        </label>
        <button type="submit" disabled={loading}>
          {loading ? '...' : t.submitLogin}
        </button>
        <p className={styles.link}>
          {t.noAccount}{' '}
          <Link to="/register">{t.register}</Link>
        </p>
      </form>
    </div>
  )
}
