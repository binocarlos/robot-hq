import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import type { CategorySummary } from "@/types"

interface Props {
  categories: CategorySummary[]
}

function BatteryBar({ level }: { level: number }) {
  const color =
    level > 60 ? "bg-chart-2" : level > 30 ? "bg-chart-3" : "bg-chart-4"

  return (
    <div className="flex items-center gap-2">
      <div className="flex-1 h-2 rounded-full bg-muted overflow-hidden">
        <div
          className={`h-full rounded-full transition-all duration-500 ${color}`}
          style={{ width: `${level}%` }}
        />
      </div>
      <span className="text-xs tabular-nums text-muted-foreground w-9 text-right">
        {level.toFixed(0)}%
      </span>
    </div>
  )
}

export function CategoryCards({ categories }: Props) {
  if (categories.length === 0) {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {Array.from({ length: 8 }).map((_, i) => (
          <Card key={i}>
            <CardContent className="p-4">
              <div className="h-24 animate-pulse rounded bg-muted" />
            </CardContent>
          </Card>
        ))}
      </div>
    )
  }

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      {categories.map((cat) => (
        <Card key={cat.name} className="hover:border-chart-1/30 transition-colors">
          <CardContent className="p-4 space-y-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <span className="text-lg">{cat.icon}</span>
                <span className="font-semibold text-sm">{cat.name}</span>
              </div>
              <span className="text-xs text-muted-foreground tabular-nums">
                {cat.count} bots
              </span>
            </div>

            <BatteryBar level={cat.avg_battery} />

            <div className="flex items-center gap-1.5 flex-wrap">
              <Badge variant="secondary" className="text-[10px] px-1.5 py-0">
                {cat.active_count} engaged
              </Badge>
              {cat.error_count > 0 && (
                <Badge variant="destructive" className="text-[10px] px-1.5 py-0">
                  {cat.error_count} marvin
                </Badge>
              )}
              {cat.charging_count > 0 && (
                <Badge variant="outline" className="text-[10px] px-1.5 py-0">
                  {cat.charging_count} napping
                </Badge>
              )}
            </div>

            <div className="flex justify-between text-xs text-muted-foreground">
              <span>Task rate: {cat.avg_task_rate}%</span>
              <span>CPU: {cat.avg_cpu_temp}°C</span>
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
