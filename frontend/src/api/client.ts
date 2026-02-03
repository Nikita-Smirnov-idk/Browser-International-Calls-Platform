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
  token: string
}

export interface InitiateCallRequest {
  phone_number: string
}

export interface InitiateCallResponse {
  call_id: string
  session_id: string
  sdp_offer: string
  status: string
  start_time: string
  voice_token?: string
}

export interface TerminateCallRequest {
  call_id: string
}

export interface TerminateCallResponse {
  call_id: string
  duration: number
  status: string
}

export interface CallHistoryItem {
  callId: string
  phoneNumber: string
  startTime: string
  duration: number
  status: string
}

export interface HistoryResponse {
  calls: CallHistoryItem[]
  total: number
  page: number
  limit: number
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
    const body = err as { message?: string; error?: string }
    throw new Error(body.message || body.error || res.statusText)
  }
  if (res.status === 204) {
    return undefined as T
  }
  return res.json()
}

export const api = {
  register: (body: RegisterRequest) =>
    request<AuthResponse>('/api/auth/register', { method: 'POST', body: JSON.stringify(body) }),

  login: (body: LoginRequest) =>
    request<AuthResponse>('/api/auth/login', { method: 'POST', body: JSON.stringify(body) }),

  logout: (token: string) =>
    request<void>('/api/auth/logout', { method: 'POST', token }),

  initiateCall: (token: string, body: InitiateCallRequest) =>
    request<InitiateCallResponse>('/api/calls/initiate', { method: 'POST', body: JSON.stringify(body), token }),

  terminateCall: (token: string, body: TerminateCallRequest) =>
    request<TerminateCallResponse>('/api/calls/terminate', { method: 'POST', body: JSON.stringify(body), token }),

  getHistory: (token: string) =>
    request<HistoryResponse>('/api/calls/history', { token }),

  health: () => request<HealthResponse>('/system/health'),
}
