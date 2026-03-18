import { Bot } from "lucide-react"
import type { ConnectionStatus } from "@/hooks/useSSE"

const statusColors: Record<ConnectionStatus, string> = {
  connected: "bg-green-500",
  reconnecting: "bg-yellow-500",
  disconnected: "bg-red-500",
}

const statusLabels: Record<ConnectionStatus, string> = {
  connected: "Live",
  reconnecting: "Reconnecting...",
  disconnected: "Disconnected",
}

interface Props {
  connectionStatus: ConnectionStatus
}

export function DashboardHeader({ connectionStatus }: Props) {
  return (
    <header className="flex items-center justify-between border-b border-border px-6 py-4">
      <div className="flex items-center gap-3">
        <Bot className="h-8 w-8 text-primary" />
        <h1 className="text-2xl font-bold tracking-tight text-foreground">
          Robot HQ
        </h1>
      </div>
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <span
          className={`h-2.5 w-2.5 rounded-full ${statusColors[connectionStatus]} ${
            connectionStatus === "connected" ? "animate-pulse-dot" : ""
          }`}
        />
        {statusLabels[connectionStatus]}
      </div>
    </header>
  )
}
