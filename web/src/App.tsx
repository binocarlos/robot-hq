import { useState, useEffect } from "react"
import { useFleetStream } from "@/hooks/useFleetStream"
import { DashboardHeader } from "@/components/DashboardHeader"
import { FleetKPIBar } from "@/components/FleetKPIBar"
import { FleetActivityChart } from "@/components/FleetActivityChart"
import { CategoryBreakdown } from "@/components/CategoryBreakdown"
import { StatusDistribution } from "@/components/StatusDistribution"
import { CategoryCards } from "@/components/CategoryCards"
import { LiveActivityFeed } from "@/components/LiveActivityFeed"
import type { CategorySummary } from "@/types"

export default function App() {
  const { snapshot, history, status } = useFleetStream()
  const [categories, setCategories] = useState<CategorySummary[]>([])

  useEffect(() => {
    let cancelled = false

    async function fetchCategories() {
      try {
        const res = await fetch("/api/categories")
        if (res.ok && !cancelled) {
          setCategories(await res.json())
        }
      } catch {
        // retry on next interval
      }
    }

    fetchCategories()
    const interval = setInterval(fetchCategories, 2000)

    return () => {
      cancelled = true
      clearInterval(interval)
    }
  }, [])

  return (
    <div className="min-h-screen bg-background">
      <DashboardHeader connectionStatus={status} />
      <main className="mx-auto max-w-7xl space-y-6 p-6">
        <FleetKPIBar snapshot={snapshot} />

        <FleetActivityChart history={history} />

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <CategoryBreakdown categories={categories} />
          <StatusDistribution snapshot={snapshot} />
        </div>

        <CategoryCards categories={categories} />

        <LiveActivityFeed events={snapshot?.events ?? []} />
      </main>
    </div>
  )
}
