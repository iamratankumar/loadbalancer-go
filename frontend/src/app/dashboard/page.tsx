"use client"

import { useEffect, useState } from "react"

import { DashboardCard } from "@/components/DashboardCard"
import { Activity, Server, Shuffle } from "lucide-react"
import api from "../../../lib/axios"

interface StatusResponse {
  blockedIPs: string[]
  servers: { Address: string; Health: boolean }[]
  strategy: string
  totalRequests: number
}

export default function DashboardPage() {
  const [status, setStatus] = useState<StatusResponse | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const fetchStatus = () => {
      api.get("/status")
        .then((res) => setStatus(res.data))
        .catch((err) => console.error("Failed to fetch status:", err))
        .finally(() => setLoading(false))
    }

    fetchStatus() 

    const interval = setInterval(fetchStatus, 10000) 

    return () => clearInterval(interval)
  }, [])

  return (
    <div className="grid gap-6 md:grid-cols-3 p-6">
      {loading ? (
          <div className="col-span-3 flex items-center justify-center">
            <div className="flex flex-col items-center gap-3">
              <div className="h-8 w-8 animate-spin rounded-full border-4 border-muted border-t-primary" />
              <p className="text-sm text-muted-foreground">Loading dashboard...</p>
            </div>
          </div>
        ) : (

        <>
          <DashboardCard
            title="Total Requests"
            value={status?.totalRequests ?? "-"}
            icon={<Activity className="h-6 w-6 text-muted-foreground" />}
          />
          <DashboardCard
            title="Active Servers"
            value={status?.servers.filter((s) => s.Health).length ?? "-"}
            icon={<Server className="h-6 w-6 text-muted-foreground" />}
          />
          <DashboardCard
            title="Strategy"
            value={status?.strategy.replace("-", " ").toUpperCase() ?? "-"}
            icon={<Shuffle className="h-6 w-6 text-muted-foreground" />}
          />
        </>
      )}
    </div>
  )
}
