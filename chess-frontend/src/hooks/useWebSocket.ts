import { useEffect, useRef, useState, useCallback } from 'react'
import { getAuth } from '../store/authStore'

export function useWebSocket(onMessage: (data: unknown) => void) {
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)
  const closedRef = useRef(false)
  const [connected, setConnected] = useState(false)
  const callbackRef = useRef(onMessage)
  callbackRef.current = onMessage

  useEffect(() => {
    const { token } = getAuth()
    if (!token) return

    closedRef.current = false

    function connect() {
      if (closedRef.current) return
      const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const ws = new WebSocket(`${proto}//${window.location.host}/ws?token=${token}`)
      wsRef.current = ws

      ws.onopen = () => {
        setConnected(true)
      }
      ws.onerror = _err => {
      }
      ws.onclose = event => {
        setConnected(false)
        wsRef.current = null
        if (!closedRef.current && event.code !== 1000) {
          reconnectTimerRef.current = setTimeout(connect, 3000)
        }
      }
      ws.onmessage = event => {
        try {
          const data = JSON.parse(event.data)
          callbackRef.current(data)
        } catch {
        }
      }
    }

    connect()

    return () => {
      closedRef.current = true
      if (reconnectTimerRef.current) {
        clearTimeout(reconnectTimerRef.current)
        reconnectTimerRef.current = null
      }
      const ws = wsRef.current
      if (ws) {
        ws.onopen = null
        ws.onclose = null
        ws.onerror = null
        ws.onmessage = null
        if (ws.readyState === WebSocket.OPEN) {
          ws.close()
        } else if (ws.readyState === WebSocket.CONNECTING) {
          ws.close()
        }
      }
      wsRef.current = null
      setConnected(false)
    }
  }, [])

  const send = useCallback((data: unknown) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(data))
    }
  }, [])

  return { connected, send }
}
