"use client"

import { useEffect, useState } from "react"
import { CheckCircle, XCircle } from "lucide-react"
import api from "../../../lib/axios"

interface Server {
  address: string
  health: boolean
  connections: number
}

interface StatusResponse {
  servers: Server[]
  totalRequests: number
  strategy: string
  blockedIPs: string[]
}

export default function ServersPage() {
  const [status, setStatus] = useState<StatusResponse | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    api.get("/status")
      .then((res) => setStatus(res.data))
      .catch((err) => console.error("Failed to fetch status:", err))
      .finally(() => setLoading(false))
  }, [])

  return (
    <div className="p-6">
      <h1 className="text-3xl font-bold mb-6">Servers</h1>

      <div className="overflow-x-auto">
        {loading ? (
          <div className="text-center text-muted-foreground">Loading servers...</div>
        ) : !status || status.servers.length === 0 ? (
          <div className="text-center text-muted-foreground">No servers found.</div>
        ) : (
          <table className="min-w-full table-auto border-collapse rounded-md bg-card text-card-foreground shadow-md dark:bg-[#1e1e1e]">
            <thead className="bg-muted">
              <tr>
                <th className="px-6 py-3 text-left text-sm font-semibold">Address</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">Health</th>
                <th className="px-6 py-3 text-left text-sm font-semibold">Connections</th>
              </tr>
            </thead>
            <tbody>
              {status.servers.map((server) => (
                <tr key={server.address} className="border-t">
                  <td className="px-6 py-4">{server.address}</td>
                  <td className="px-6 py-4">
                    {server.health ? (
                      <div className="flex items-center gap-2 text-green-600 dark:text-green-400">
                        <CheckCircle className="h-5 w-5" />
                        Up
                      </div>
                    ) : (
                      <div className="flex items-center gap-2 text-red-600 dark:text-red-400">
                        <XCircle className="h-5 w-5" />
                        Down
                      </div>
                    )}
                  </td>
                  <td className="px-6 py-4">{server.connections}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
