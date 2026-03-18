import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from "recharts"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import type { FleetSnapshot } from "@/types"

const STATUS_CONFIG = [
  { key: "active_count", label: "Engaged", color: "var(--color-chart-2)" },
  { key: "idle_count", label: "Needs More Data", color: "#cbd5e1" },
  { key: "charging_count", label: "Power Napping", color: "var(--color-chart-3)" },
  { key: "error_count", label: "Marvin Depressed", color: "var(--color-chart-4)" },
  { key: "maintenance_count", label: "Spa Day", color: "var(--color-chart-5)" },
] as const

interface Props {
  snapshot: FleetSnapshot | null
}

export function StatusDistribution({ snapshot }: Props) {
  const data = snapshot
    ? STATUS_CONFIG.map((s) => ({
        name: s.label,
        value: snapshot[s.key],
        color: s.color,
      }))
    : []

  return (
    <Card className="flex-1">
      <CardHeader className="pb-2">
        <CardTitle className="text-base">Status Distribution</CardTitle>
      </CardHeader>
      <CardContent>
        {!snapshot ? (
          <div className="flex h-[200px] items-center justify-center text-muted-foreground text-sm">
            Loading...
          </div>
        ) : (
          <div className="flex items-center gap-4">
            <ResponsiveContainer width="100%" height={200}>
              <PieChart>
                <Pie
                  data={data}
                  cx="50%"
                  cy="50%"
                  innerRadius={50}
                  outerRadius={80}
                  paddingAngle={2}
                  dataKey="value"
                  isAnimationActive={false}
                >
                  {data.map((entry, i) => (
                    <Cell key={i} fill={entry.color} />
                  ))}
                </Pie>
                <Tooltip
                  contentStyle={{
                    backgroundColor: "var(--color-card)",
                    border: "1px solid var(--color-border)",
                    borderRadius: "8px",
                    fontSize: "12px",
                  }}
                />
              </PieChart>
            </ResponsiveContainer>
            <div className="flex flex-col gap-1.5 min-w-[120px]">
              {data.map((d) => (
                <div key={d.name} className="flex items-center gap-2 text-xs">
                  <span
                    className="h-2.5 w-2.5 rounded-full shrink-0"
                    style={{ backgroundColor: d.color }}
                  />
                  <span className="text-muted-foreground">{d.name}</span>
                  <span className="ml-auto font-medium tabular-nums">{d.value}</span>
                </div>
              ))}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
