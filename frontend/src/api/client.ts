const API_BASE = ''

export interface RegisterRequest {
  email: string
  password: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface AuthResponse {
  accessToken: string
}

export interface StartCallRequest {
  countryCode: string
  phoneNumber: string
}

export interface StartCallResponse {
  callId: string
  rtcConfig?: Record<string, unknown>
}

export interface EndCallRequest {
  callId: string
}

export interface CallHistoryItem {
  callId: string
  phoneNumber: string
  startedAt: string
  durationSeconds: number
  status: 'completed' | 'failed' | 'canceled'
}

export interface HealthResponse {
  status: string
}

async function request<T>(
  path: string,
  options: RequestInit & { token?: string } = {}
): Promise<T> {
  const { token, ...init } = options
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(init.headers as Record<string, string>),
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }
  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({}))
    throw new Error((err as { error?: string }).error || res.statusText)
  }
  if (res.status === 204) {
    return undefined as T
  }
  return res.json()
}

export const api = {
  register: (body: RegisterRequest) =>
    request<void>('/auth/register', { method: 'POST', body: JSON.stringify(body) }),

  login: (body: LoginRequest) =>
    request<AuthResponse>('/auth/login', { method: 'POST', body: JSON.stringify(body) }),

  logout: (token: string) =>
    request<void>('/auth/logout', { method: 'POST', token }),

  startCall: (token: string, body: StartCallRequest) =>
    request<StartCallResponse>('/calls/start', { method: 'POST', body: JSON.stringify(body), token }),

  endCall: (token: string, body: EndCallRequest) =>
    request<void>('/calls/end', { method: 'POST', body: JSON.stringify(body), token }),

  getHistory: (token: string) =>
    request<CallHistoryItem[]>('/calls/history', { token }),

  health: () => request<HealthResponse>('/system/health'),
}
