import { useCallback, useEffect, useRef, useState } from 'react'
import { useAuth } from '../contexts/AuthContext'
import { useLocale } from '../i18n/LocaleContext'
import { api } from '../api/client'
import { countries } from '../data/countries'
import { CountrySelect } from '../components/CountrySelect'
import { useDevices } from '../hooks/useDevices'
import { useWebRTC } from '../hooks/useWebRTC'
import { Layout } from '../components/Layout'
import styles from './Call.module.css'

const MIN_PHONE_LENGTH = 5

type CallStatus = 'idle' | 'connecting' | 'talking' | 'finished'

export function Call() {
  const { token } = useAuth()
  const { t, locale } = useLocale()
  const { microphoneOk, permissionDenied, error: deviceError, recheck: recheckDevices } = useDevices()
  const { startLocalStream, cleanup: cleanupWebRTC } = useWebRTC()
  const [country, setCountry] = useState(countries[0])
  const [phoneNumber, setPhoneNumber] = useState('')
  const [status, setStatus] = useState<CallStatus>('idle')
  const [callId, setCallId] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const streamRef = useRef<MediaStream | null>(null)

  const fullPhone = `${country.dialCode}${phoneNumber.replace(/\D/g, '')}`
  const canCall =
    phoneNumber.replace(/\D/g, '').length >= MIN_PHONE_LENGTH &&
    (microphoneOk === true || microphoneOk === null)

  const cleanup = useCallback(() => {
    if (streamRef.current) {
      streamRef.current.getTracks().forEach((tr) => tr.stop())
      streamRef.current = null
    }
    cleanupWebRTC()
  }, [cleanupWebRTC])

  useEffect(() => {
    return cleanup
  }, [cleanup])

  const handleStartCall = useCallback(async () => {
    if (!token) return
    setError(null)
    setStatus('connecting')
    try {
      let stream: MediaStream | null = null
      try {
        stream = await startLocalStream()
        streamRef.current = stream
      } catch (e) {
        setError(t.noMicrophone)
        setStatus('idle')
        return
      }
      const res = await api.initiateCall(token, { phone_number: fullPhone })
      setCallId(res.call_id)
      setStatus('talking')
    } catch (e) {
      setError(e instanceof Error ? e.message : t.connectionError)
      setStatus('idle')
      cleanup()
    }
  }, [token, country, fullPhone, startLocalStream, t.noMicrophone, t.connectionError, cleanup])

  const handleEndCall = useCallback(async () => {
    if (!token) return
    if (callId) {
      try {
        await api.terminateCall(token, { call_id: callId })
      } catch {
        //
      }
    }
    setCallId(null)
    setStatus('finished')
    cleanup()
    setTimeout(() => setStatus('idle'), 2000)
  }, [token, callId, cleanup])

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

        {deviceError && microphoneOk === false && (
          <div className={styles.warning}>
            {permissionDenied ? t.allowMicrophone : t.noMicrophone}
            <button type="button" className={styles.recheckBtn} onClick={recheckDevices}>
              {locale === 'ru' ? 'Проверить снова' : 'Check again'}
            </button>
          </div>
        )}

        {error && (
          <div className={styles.error} role="alert">
            {error}
          </div>
        )}

        <div className={styles.form}>
          <label>
            {t.country}
            <CountrySelect
              countries={countries}
              value={country}
              onChange={setCountry}
              locale={locale}
              disabled={status === 'connecting' || status === 'talking'}
            />
          </label>

          <label>
            {t.phoneNumber}
            <div className={styles.phoneRow}>
              <span className={styles.dialCode}>{country.dialCode}</span>
              <input
                type="tel"
                inputMode="numeric"
                value={phoneNumber}
                onChange={(e) => setPhoneNumber(e.target.value.replace(/\D/g, ''))}
                placeholder="9123456789"
                disabled={status === 'connecting' || status === 'talking'}
                maxLength={15}
              />
            </div>
            {phoneNumber.replace(/\D/g, '').length > 0 && phoneNumber.replace(/\D/g, '').length < MIN_PHONE_LENGTH && (
              <span className={styles.hint}>
                {locale === 'ru' ? `Минимум ${MIN_PHONE_LENGTH} цифр` : `At least ${MIN_PHONE_LENGTH} digits`}
              </span>
            )}
          </label>

          {status === 'idle' && (
            <button
              className={styles.callButton}
              onClick={handleStartCall}
              disabled={!canCall}
            >
              {t.startCall}
            </button>
          )}

          {(status === 'connecting' || status === 'talking') && (
            <button className={styles.endButton} onClick={handleEndCall}>
              {t.endCall}
            </button>
          )}
        </div>

        {getStatusText() && (
          <div className={styles.statusWrap}>
            <span
              className={`${styles.statusDot} ${status === 'connecting' ? styles.pulse : ''} ${status === 'talking' ? styles.active : ''} ${status === 'finished' ? styles.done : ''}`}
            />
            <p className={styles.status}>{getStatusText()}</p>
          </div>
        )}
      </div>
    </Layout>
  )
}
