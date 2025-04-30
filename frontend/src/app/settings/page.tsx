"use client"

import { useState, useEffect } from "react"

import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectItem
} from "@/components/ui/Select"
import { CheckCircle2 } from "lucide-react"
import api from "../../../lib/axios"

export default function SettingsPage() {
  const [strategy, setStrategy] = useState("round-robin")
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState("")

  useEffect(() => {
    api.get("/status")
      .then((res) => setStrategy(res.data.strategy))
      .catch(() => {})
  }, [])

  const handleChange = async (value: string) => {
    setSaving(true)
    try {
      await api.post("/set-strategy", { strategy: value })
      setStrategy(value)
      setMessage("Strategy updated successfully!")
    } catch (err) {
      console.error("Failed to update strategy:", err)
      setMessage("Failed to update strategy.")
    } finally {
      setSaving(false)
      setTimeout(() => setMessage(""), 3000)
    }
  }

  return (
    <div className="p-6">
      <h1 className="text-3xl font-bold mb-6">Settings</h1>

      <div className="rounded-xl border bg-card text-card-foreground p-6 shadow-md dark:bg-[#1e1e1e] space-y-6 max-w-lg">
        <div className="space-y-2">
          <h2 className="text-lg font-semibold">Load Balancing Strategy</h2>
          <p className="text-sm text-muted-foreground">
            Choose how requests are distributed across servers.
          </p>
        </div>

        <Select value={strategy} onValueChange={handleChange}>
          <SelectTrigger className="w-full">
            <SelectValue placeholder="Select a strategy" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="round-robin">Round Robin</SelectItem>
            <SelectItem value="least-connections">Least Connections</SelectItem>
          </SelectContent>
        </Select>

        {message && (
          <div className="flex items-center gap-2 rounded-md bg-green-100 px-3 py-2 text-green-800 dark:bg-green-900 dark:text-green-200">
            <CheckCircle2 className="h-4 w-4" />
            <span className="text-sm">{message}</span>
          </div>
        )}
      </div>
    </div>
  )
}
