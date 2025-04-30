"use client"

import { useEffect, useRef, useState } from "react"
import api from "../../../lib/axios"

export default function LogsPage() {
  const [logs, setLogs] = useState<string[]>([])
  const [loading, setLoading] = useState(true)
  const logRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    api.get("/logs", { responseType: "text" })
      .then((res) => {
        const clean = res.data
          .replace(/^\[|\]$/g, "")
          .replace(/","/g, "\n")
          .replace(/\\n/g, "\n")
          .split("\n")
          .filter(Boolean)
        setLogs(clean)
      })
      .catch((err) => {
        console.error("Failed to fetch logs:", err)
      })
      .finally(() => {
        setLoading(false)
      })
  }, [])

  useEffect(() => {
    if (logRef.current) {
      logRef.current.scrollTo({
        top: logRef.current.scrollHeight,
        behavior: "smooth",
      })
    }
  }, [logs])

  const getTagStyle = (line: string) => {
    if (line.includes("[SUCCESS]")) return "bg-green-500 text-white"
    if (line.includes("[FAIL]") || line.includes("failed")) return "bg-red-500 text-white"
    if (line.includes("[REQUEST]")) return "bg-blue-500 text-white"
    if (line.includes("[STRATEGY SWITCH]")) return "bg-yellow-500 text-black"
    if (line.includes("[STARTUP]")) return "bg-purple-600 text-white"
    if (line.includes("[SNAPSHOT]")) return "bg-cyan-500 text-white"
    return "bg-muted text-foreground"
  }

  const extractTimestamp = (line: string) => {
    const match = line.match(/^(\d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2})/)
    return match ? match[1] : ""
  }

  const extractTag = (line: string) => {
    const match = line.match(/\[(.*?)\]/)
    return match ? match[0] : ""
  }

  const extractMessage = (line: string) => {
    return line
      .replace(/^(\d{4}\/\d{2}\/\d{2} \d{2}:\d{2}:\d{2})/, "")
      .replace(/\[(.*?)\]/, "")
      .trim()
  }

  return (
    <div className="p-6">
      <h1 className="text-3xl font-bold mb-6">Logs</h1>

      <div
        ref={logRef}
        className="rounded-xl border bg-card dark:bg-[#0f0f0f] text-card-foreground p-4 shadow-md overflow-y-auto h-[500px] font-mono text-sm space-y-2"
      >
        {loading ? (
          <div className="text-center text-muted-foreground">Loading logs...</div>
        ) : logs.length > 0 ? (
          logs.map((log, idx) => {
            const time = extractTimestamp(log)
            const tag = extractTag(log)
            const msg = extractMessage(log)

            return (
              <div
                key={idx}
                className="flex items-start gap-3 px-2 py-1 hover:bg-muted/10 transition rounded"
              >
                <span className="text-muted-foreground w-36 shrink-0">{time}</span>
                <span
                  className={`px-2 py-0.5 rounded text-xs font-semibold uppercase ${getTagStyle(
                    tag
                  )}`}
                >
                  {tag || "INFO"}
                </span>
                <span className="whitespace-pre-wrap break-words">{msg}</span>
              </div>
            )
          })
        ) : (
          <div className="text-center text-muted-foreground">No logs found.</div>
        )}
      </div>
    </div>
  )
}
