export interface User {
  id: string
  email: string
  username: string
}

interface PersistedAuth {
  user: User
  token: string
}

interface AuthState {
  user: User | null
  token: string | null
}

const STORAGE_KEY = 'chess_auth'

function loadPersisted(): AuthState {
  try {
    const raw = sessionStorage.getItem(STORAGE_KEY)
    if (raw) {
      const p: PersistedAuth = JSON.parse(raw)
      return { user: p.user, token: p.token }
    }
  } catch {
  }
  return { user: null, token: null }
}

function persist(auth: AuthState) {
  if (auth.user && auth.token) {
    sessionStorage.setItem(STORAGE_KEY, JSON.stringify({ user: auth.user, token: auth.token }))
  } else {
    sessionStorage.removeItem(STORAGE_KEY)
  }
}

let state: AuthState = loadPersisted()
const listeners = new Set<() => void>()

export function getAuth(): AuthState {
  return state
}

export function setAuth(user: User, token: string) {
  state = { user, token }
  persist(state)
  listeners.forEach(fn => fn())
}

export function clearAuth() {
  state = { user: null, token: null }
  persist(state)
  listeners.forEach(fn => fn())
}

export function subscribeAuth(fn: () => void) {
  listeners.add(fn)
  return () => listeners.delete(fn)
}

export function isAuthenticated(): boolean {
  return state.token !== null
}
