import {
  createContext,
  useCallback,
  useContext,
  useState,
  type ReactNode,
} from 'react'
import { api } from '../api/client'

const TOKEN_KEY = 'accessToken'

interface AuthContextValue {
  token: string | null
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string) => Promise<void>
  logout: () => Promise<void>
  error: string | null
  clearError: () => void
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(() =>
    localStorage.getItem(TOKEN_KEY)
  )
  const [error, setError] = useState<string | null>(null)

  const login = useCallback(async (email: string, password: string) => {
    setError(null)
    try {
      const { accessToken } = await api.login({ email, password })
      localStorage.setItem(TOKEN_KEY, accessToken)
      setToken(accessToken)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Login failed')
      throw e
    }
  }, [])

  const register = useCallback(async (email: string, password: string) => {
    setError(null)
    try {
      await api.register({ email, password })
      await login(email, password)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Registration failed')
      throw e
    }
  }, [login])

  const logout = useCallback(async () => {
    if (token) {
      try {
        await api.logout(token)
      } catch {
        // ignore
      }
      localStorage.removeItem(TOKEN_KEY)
      setToken(null)
    }
  }, [token])

  const clearError = useCallback(() => setError(null), [])

  return (
    <AuthContext.Provider
      value={{
        token,
        isAuthenticated: !!token,
        login,
        register,
        logout,
        error,
        clearError,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return ctx
}
