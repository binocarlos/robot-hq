import { useState, useEffect, useRef } from "react"
import { useSSE } from "./useSSE"
import type { FleetSnapshot } from "@/types"

export interface FleetHistoryPoint {
  time: string
  eventsPerSec: number
  tasksDone: number
  cpuTemp: number
  latency: number
}

const MAX_HISTORY = 120

function generateSyntheticHistory(): FleetHistoryPoint[] {
  const now = Date.now()
  const points: FleetHistoryPoint[] = []

  for (let i = MAX_HISTORY; i > 0; i--) {
    const t = now - i * 250
    const elapsed = t / 1000

    // Mirror the backend wave patterns
    const wave1 = 8 * Math.sin((elapsed / 6) * 2 * Math.PI)
    const wave2 = 5 * Math.sin((elapsed / 10) * 2 * Math.PI)
    const wave3 = 3 * Math.sin((elapsed / 26) * 2 * Math.PI)
    const noise = (Math.random() - 0.5) * 10

    const taskWave = 4 * Math.sin((elapsed / 8) * 2 * Math.PI) + 3 * Math.sin((elapsed / 15) * 2 * Math.PI)
    const taskNoise = (Math.random() - 0.5) * 6

    const latWave = 15 * Math.sin((elapsed / 6) * 2 * Math.PI) + 10 * Math.sin((elapsed / 14) * 2 * Math.PI) + 8 * Math.sin((elapsed / 26) * 2 * Math.PI)
    const latNoise = (Math.random() - 0.5) * 20

    points.push({
      time: new Date(t).toLocaleTimeString(),
      eventsPerSec: Math.round(Math.max(0, 20 + wave1 + wave2 + wave3 + noise + 80) * 10) / 10,
      tasksDone: Math.max(0, Math.round(12 + taskWave + taskNoise + 10)),
      cpuTemp: Math.round((55 + (Math.random() - 0.5) * 6) * 10) / 10,
      latency: Math.round(Math.max(5, 45 + latWave + latNoise) * 10) / 10,
    })
  }

  return points
}

export function useFleetStream() {
  const { data, status } = useSSE<FleetSnapshot>("/api/stream/fleet")
  const [history, setHistory] = useState<FleetHistoryPoint[]>(() => generateSyntheticHistory())
  const lastTimestampRef = useRef<number>(0)

  useEffect(() => {
    if (!data || data.timestamp === lastTimestampRef.current) return
    lastTimestampRef.current = data.timestamp

    const point: FleetHistoryPoint = {
      time: new Date(data.timestamp).toLocaleTimeString(),
      eventsPerSec: data.events_per_sec,
      tasksDone: data.tasks_completed,
      cpuTemp: data.avg_cpu_temp,
      latency: data.network_latency,
    }

    setHistory((prev) => {
      const next = [...prev, point]
      return next.length > MAX_HISTORY ? next.slice(-MAX_HISTORY) : next
    })
  }, [data])

  return { snapshot: data, history, status }
}
