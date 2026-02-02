import { useCallback, useState } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { useLocale } from '../i18n/LocaleContext'
import { api } from '../api/client'
import { countries } from '../data/countries'
import { Layout } from '../components/Layout'
import styles from './Call.module.css'

type CallStatus = 'idle' | 'connecting' | 'talking' | 'finished'

export function Call() {
  const { token } = useAuth()
  const { t, locale } = useLocale()
  const [country, setCountry] = useState(countries[0])
  const [phoneNumber, setPhoneNumber] = useState('')
  const [status, setStatus] = useState<CallStatus>('idle')
  const [callId, setCallId] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  const handleStartCall = useCallback(async () => {
    if (!token) return
    setError(null)
    setStatus('connecting')
    try {
      const res = await api.startCall(token, {
        countryCode: country.dialCode,
        phoneNumber: phoneNumber.replace(/\D/g, ''),
      })
      setCallId(res.callId)
      setStatus('talking')
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Failed to start call')
      setStatus('idle')
    }
  }, [token, country, phoneNumber])

  const handleEndCall = useCallback(async () => {
    if (!token || !callId) return
    try {
      await api.endCall(token, { callId })
    } catch {
      // ignore
    }
    setCallId(null)
    setStatus('finished')
    setTimeout(() => setStatus('idle'), 2000)
  }, [token, callId])

  const getStatusText = () => {
    switch (status) {
      case 'connecting':
        return t.callStatus.connecting
      case 'talking':
        return t.callStatus.talking
      case 'finished':
        return t.callStatus.finished
      default:
        return ''
    }
  }

  return (
    <Layout>
      <div className={styles.page}>
        <h1>{t.call}</h1>

        {error && (
          <div className={styles.error}>
            {error}
          </div>
        )}

        <div className={styles.form}>
          <label>
            {t.country}
            <select
              value={country.code}
              onChange={(e) => {
                const c = countries.find((x) => x.code === e.target.value)
                if (c) setCountry(c)
              }}
              disabled={status === 'connecting' || status === 'talking'}
            >
              {countries.map((c) => (
                <option key={c.code} value={c.code}>
                  {c.flag} {locale === 'ru' ? c.nameRu : c.nameEn} ({c.dialCode})
                </option>
              ))}
            </select>
          </label>

          <label>
            {t.phoneNumber}
            <div className={styles.phoneRow}>
              <span className={styles.dialCode}>{country.dialCode}</span>
              <input
                type="tel"
                value={phoneNumber}
                onChange={(e) => setPhoneNumber(e.target.value.replace(/\D/g, ''))}
                placeholder="9123456789"
                disabled={status === 'connecting' || status === 'talking'}
              />
            </div>
          </label>

          {status === 'idle' && (
            <button
              className={styles.callButton}
              onClick={handleStartCall}
              disabled={!phoneNumber.trim()}
            >
              {t.startCall}
            </button>
          )}

          {(status === 'connecting' || status === 'talking') && (
            <button
              className={styles.endButton}
              onClick={handleEndCall}
            >
              {t.endCall}
            </button>
          )}
        </div>

        {getStatusText() && (
          <p className={styles.status}>{getStatusText()}</p>
        )}
      </div>
    </Layout>
  )
}
