import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts"
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import type { FleetHistoryPoint } from "@/hooks/useFleetStream"

interface Props {
  history: FleetHistoryPoint[]
}

export function FleetActivityChart({ history }: Props) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-base">Fleet Activity</CardTitle>
      </CardHeader>
      <CardContent>
        {history.length === 0 ? (
          <div className="flex h-[300px] items-center justify-center text-muted-foreground text-sm">
            Waiting for data...
          </div>
        ) : (
          <ResponsiveContainer width="100%" height={300}>
            <LineChart data={history}>
              <CartesianGrid strokeDasharray="3 3" stroke="var(--color-border)" />
              <XAxis
                dataKey="time"
                tick={{ fontSize: 11, fill: "var(--color-muted-foreground)" }}
                tickLine={false}
                axisLine={false}
              />
              <YAxis
                tick={{ fontSize: 11, fill: "var(--color-muted-foreground)" }}
                tickLine={false}
                axisLine={false}
              />
              <Tooltip
                contentStyle={{
                  backgroundColor: "var(--color-card)",
                  border: "1px solid var(--color-border)",
                  borderRadius: "8px",
                  fontSize: "12px",
                }}
              />
              <Legend />
              <Line
                type="monotone"
                dataKey="eventsPerSec"
                stroke="#10b981"
                strokeWidth={2}
                dot={false}
                name="Events/sec"
                isAnimationActive={false}
              />
              <Line
                type="monotone"
                dataKey="tasksDone"
                stroke="#ec4899"
                strokeWidth={2}
                dot={false}
                name="Tasks Done"
                isAnimationActive={false}
              />
              <Line
                type="monotone"
                dataKey="cpuTemp"
                stroke="#f97316"
                strokeWidth={2}
                dot={false}
                name="CPU Temp"
                isAnimationActive={false}
              />
              <Line
                type="monotone"
                dataKey="latency"
                stroke="#3b82f6"
                strokeWidth={2}
                dot={false}
                name="Latency"
                isAnimationActive={false}
              />
            </LineChart>
          </ResponsiveContainer>
        )}
      </CardContent>
    </Card>
  )
}
