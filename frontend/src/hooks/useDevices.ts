import { useCallback, useEffect, useState } from 'react'

export interface DeviceState {
  microphoneOk: boolean | null
  speakersOk: boolean | null
  permissionDenied: boolean
  checking: boolean
  error: string | null
}

export function useDevices() {
  const [state, setState] = useState<DeviceState>({
    microphoneOk: null,
    speakersOk: null,
    permissionDenied: false,
    checking: false,
    error: null,
  })

  const check = useCallback(async () => {
    if (!navigator.mediaDevices?.getUserMedia) {
      setState((s) => ({
        ...s,
        microphoneOk: false,
        error: 'MediaDevices not supported',
      }))
      return
    }
    setState((s) => ({ ...s, checking: true, error: null, permissionDenied: false }))
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      const tracks = stream.getTracks()
      tracks.forEach((t) => t.stop())
      const devices = await navigator.mediaDevices.enumerateDevices()
      const hasOutput = devices.some((d) => d.kind === 'audiooutput')
      setState({
        microphoneOk: true,
        speakersOk: hasOutput,
        permissionDenied: false,
        checking: false,
        error: null,
      })
    } catch (e) {
      const err = e instanceof Error ? e : new Error('Unknown error')
      const denied = err.name === 'NotAllowedError' || err.name === 'PermissionDeniedError'
      setState({
        microphoneOk: false,
        speakersOk: null,
        permissionDenied: denied,
        checking: false,
        error: err.message,
      })
    }
  }, [])

  useEffect(() => {
    check()
  }, [check])

  return { ...state, recheck: check }
}
