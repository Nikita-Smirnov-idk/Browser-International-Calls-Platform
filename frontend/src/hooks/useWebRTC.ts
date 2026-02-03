import { useCallback, useRef, useState } from 'react'

export function useWebRTC() {
  const pcRef = useRef<RTCPeerConnection | null>(null)
  const streamRef = useRef<MediaStream | null>(null)
  const [connectionState, setConnectionState] = useState<RTCPeerConnectionState | null>(null)

  const cleanup = useCallback(() => {
    if (streamRef.current) {
      streamRef.current.getTracks().forEach((t) => t.stop())
      streamRef.current = null
    }
    if (pcRef.current) {
      pcRef.current.close()
      pcRef.current = null
    }
    setConnectionState(null)
  }, [])

  const startLocalStream = useCallback(async (): Promise<MediaStream> => {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    streamRef.current = stream
    return stream
  }, [])

  const createConnection = useCallback(() => {
    const pc = new RTCPeerConnection()
    pcRef.current = pc
    pc.onconnectionstatechange = () => setConnectionState(pc.connectionState)
    return pc
  }, [])

  const attachStream = useCallback((pc: RTCPeerConnection, stream: MediaStream) => {
    stream.getTracks().forEach((track) => pc.addTrack(track, stream))
  }, [])

  return {
    pcRef,
    streamRef,
    connectionState,
    startLocalStream,
    createConnection,
    attachStream,
    cleanup,
  }
}
