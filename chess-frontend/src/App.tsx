import { Routes, Route, Navigate } from 'react-router-dom'
import { useState, useEffect, type ReactNode } from 'react'
import { subscribeAuth, isAuthenticated } from './store/authStore'
import { Login } from './pages/Login'
import { Register } from './pages/Register'
import { Lobby } from './pages/Lobby'
import { GameView } from './pages/GameView'

function ProtectedRoute({ children }: { children: ReactNode }) {
  return isAuthenticated() ? children : <Navigate to="/login" />
}

export default function App() {
  const [, setTick] = useState(0)

  useEffect(() => {
    const unsub = subscribeAuth(() => setTick(t => t + 1))
    return () => {
      unsub()
    }
  }, [])

  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />
      <Route
        path="/lobby"
        element={
          <ProtectedRoute>
            <Lobby />
          </ProtectedRoute>
        }
      />
      <Route
        path="/game/:id"
        element={
          <ProtectedRoute>
            <GameView />
          </ProtectedRoute>
        }
      />
      <Route path="*" element={<Navigate to="/lobby" />} />
    </Routes>
  )
}
