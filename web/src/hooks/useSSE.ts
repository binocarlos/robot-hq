import { useEffect, useRef, useState, useCallback } from "react"

export type ConnectionStatus = "connected" | "reconnecting" | "disconnected"

export function useSSE<T>(url: string) {
  const [data, setData] = useState<T | null>(null)
  const [status, setStatus] = useState<ConnectionStatus>("disconnected")
  const eventSourceRef = useRef<EventSource | null>(null)
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)

  const connect = useCallback(() => {
    if (eventSourceRef.current) {
      eventSourceRef.current.close()
    }

    const es = new EventSource(url)
    eventSourceRef.current = es

    es.onopen = () => {
      setStatus("connected")
    }

    es.onmessage = (event) => {
      try {
        const parsed = JSON.parse(event.data) as T
        setData(parsed)
        setStatus("connected")
      } catch {
        // ignore parse errors
      }
    }

    es.onerror = () => {
      es.close()
      setStatus("reconnecting")
      reconnectTimeoutRef.current = setTimeout(() => {
        connect()
      }, 2000)
    }
  }, [url])

  useEffect(() => {
    connect()

    return () => {
      if (eventSourceRef.current) {
        eventSourceRef.current.close()
      }
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
    }
  }, [connect])

  return { data, status }
}
