import { useEffect, useState } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { useLocale } from '../i18n/LocaleContext'
import { api, type CallHistoryItem } from '../api/client'
import { Layout } from '../components/Layout'
import styles from './History.module.css'

function formatDate(iso: string) {
  const d = new Date(iso)
  return d.toLocaleString(undefined, {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function formatDuration(sec: number) {
  const m = Math.floor(sec / 60)
  const s = sec % 60
  return `${m}:${s.toString().padStart(2, '0')}`
}

function getStatusText(
  status: string,
  t: { statusCompleted: string; statusFailed: string; statusCanceled: string }
) {
  switch (status) {
    case 'completed':
      return t.statusCompleted
    case 'failed':
      return t.statusFailed
    case 'canceled':
      return t.statusCanceled
    default:
      return status
  }
}

export function History() {
  const { token } = useAuth()
  const { t } = useLocale()
  const [items, setItems] = useState<CallHistoryItem[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!token) return
    api
      .getHistory(token)
      .then((res) => setItems(res.calls))
      .catch((e) => setError(e instanceof Error ? e.message : 'Failed to load'))
      .finally(() => setLoading(false))
  }, [token])

  if (loading) {
    return (
      <Layout>
        <div className={styles.page}>
          <h1>{t.history}</h1>
          <p className={styles.loading}>...</p>
        </div>
      </Layout>
    )
  }

  return (
    <Layout>
      <div className={styles.page}>
        <h1>{t.history}</h1>

        {error && (
          <div className={styles.error}>{error}</div>
        )}

        {items.length === 0 && !error ? (
          <p className={styles.empty}>{t.historyEmpty}</p>
        ) : (
          <div className={styles.tableWrap}>
            <table className={styles.table}>
              <thead>
                <tr>
                  <th>{t.historyDate}</th>
                  <th>{t.phoneNumber}</th>
                  <th>{t.historyDuration}</th>
                  <th>{t.historyStatus}</th>
                </tr>
              </thead>
              <tbody>
                {items.map((item) => (
                  <tr key={item.callId}>
                    <td>{formatDate(item.startTime)}</td>
                    <td>{item.phoneNumber}</td>
                    <td>{formatDuration(item.duration)}</td>
                    <td>{getStatusText(item.status, t)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </Layout>
  )
}
