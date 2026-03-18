export interface Robot {
  id: string
  category: string
  status: "active" | "idle" | "charging" | "error" | "maintenance"
  battery_level: number
  cpu_temp: number
  task_completion_rate: number
  distance_traveled_today: number
  error_count: number
  uptime_hours: number
  lat: number
  lng: number
  last_event?: string
  last_event_time?: number
}

export interface CategorySummary {
  name: string
  icon: string
  count: number
  active_count: number
  idle_count: number
  charging_count: number
  error_count: number
  maintenance_count: number
  avg_battery: number
  avg_cpu_temp: number
  avg_task_rate: number
  total_errors: number
}

export interface FleetEvent {
  timestamp: number
  robot_id: string
  category: string
  message: string
  type: "info" | "warning" | "error"
}

export interface FleetSnapshot {
  timestamp: number
  total_robots: number
  active_count: number
  idle_count: number
  charging_count: number
  error_count: number
  maintenance_count: number
  avg_battery: number
  avg_task_rate: number
  total_error_count: number
  events_per_sec: number
  tasks_completed: number
  alerts_active: number
  avg_cpu_temp: number
  network_latency: number
  events: FleetEvent[]
}
