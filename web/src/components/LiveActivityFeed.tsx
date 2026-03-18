import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card"
import type { FleetEvent } from "@/types"

const typeStyles: Record<string, string> = {
  info: "text-chart-1",
  warning: "text-chart-3",
  error: "text-chart-4",
}

const typeDots: Record<string, string> = {
  info: "bg-chart-1",
  warning: "bg-chart-3",
  error: "bg-chart-4",
}

interface Props {
  events: FleetEvent[]
}

export function LiveActivityFeed({ events }: Props) {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-base">Live Activity Feed</CardTitle>
      </CardHeader>
      <CardContent>
        {events.length === 0 ? (
          <div className="flex h-[200px] items-center justify-center text-muted-foreground text-sm">
            Waiting for events...
          </div>
        ) : (
          <div className="space-y-1 max-h-[280px] overflow-y-auto pr-2">
            {events.slice(0, 20).map((event, i) => (
              <div
                key={`${event.timestamp}-${event.robot_id}-${i}`}
                className={`flex items-start gap-3 py-1.5 text-sm ${
                  i === 0 ? "animate-slide-in" : ""
                }`}
              >
                <span
                  className={`mt-1.5 h-2 w-2 rounded-full shrink-0 ${
                    typeDots[event.type] || typeDots.info
                  }`}
                />
                <span className="text-xs text-muted-foreground tabular-nums shrink-0 mt-0.5">
                  {new Date(event.timestamp).toLocaleTimeString()}
                </span>
                <span className="font-mono text-xs font-medium shrink-0">
                  {event.robot_id}
                </span>
                <span
                  className={`text-xs ${typeStyles[event.type] || typeStyles.info}`}
                >
                  {event.message}
                </span>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
