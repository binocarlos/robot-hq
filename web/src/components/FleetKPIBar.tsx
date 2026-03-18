import { Card, CardContent } from "@/components/ui/card"
import {
  Activity,
  Battery,
  AlertTriangle,
  Pause,
  Zap,
  Users,
} from "lucide-react"
import type { FleetSnapshot } from "@/types"

interface KPICardProps {
  label: string
  value: string | number
  icon: React.ReactNode
  color?: string
}

function KPICard({ label, value, icon, color = "text-foreground" }: KPICardProps) {
  return (
    <Card className="flex-1 min-w-[140px]">
      <CardContent className="p-4">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm text-muted-foreground">{label}</p>
            <p className={`text-2xl font-bold tabular-nums ${color}`}>
              {value}
            </p>
          </div>
          <div className="text-muted-foreground">{icon}</div>
        </div>
      </CardContent>
    </Card>
  )
}

interface Props {
  snapshot: FleetSnapshot | null
}

export function FleetKPIBar({ snapshot }: Props) {
  if (!snapshot) {
    return (
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
        {Array.from({ length: 6 }).map((_, i) => (
          <Card key={i} className="flex-1 min-w-[140px]">
            <CardContent className="p-4">
              <div className="h-12 animate-pulse rounded bg-muted" />
            </CardContent>
          </Card>
        ))}
      </div>
    )
  }

  return (
    <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
      <KPICard
        label="Total Robots"
        value={snapshot.total_robots.toLocaleString()}
        icon={<Users className="h-5 w-5" />}
      />
      <KPICard
        label="Engaged"
        value={snapshot.active_count.toLocaleString()}
        icon={<Activity className="h-5 w-5" />}
        color="text-chart-2"
      />
      <KPICard
        label="Needs More Data"
        value={snapshot.idle_count.toLocaleString()}
        icon={<Pause className="h-5 w-5" />}
        color="text-muted-foreground"
      />
      <KPICard
        label="Power Napping"
        value={snapshot.charging_count.toLocaleString()}
        icon={<Zap className="h-5 w-5" />}
        color="text-chart-3"
      />
      <KPICard
        label="Marvin Depressed"
        value={snapshot.error_count}
        icon={<AlertTriangle className="h-5 w-5" />}
        color="text-chart-4"
      />
      <KPICard
        label="Avg Battery"
        value={`${snapshot.avg_battery}%`}
        icon={<Battery className="h-5 w-5" />}
        color="text-chart-1"
      />
    </div>
  )
}
