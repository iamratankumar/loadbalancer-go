"use client"

import { cn } from "@/lib/utils"

interface DashboardCardProps {
  title: string
  value: string | number
  icon?: React.ReactNode
}

export function DashboardCard({ title, value, icon }: DashboardCardProps) {
  return (
    <div
      className={cn(
        "flex flex-col justify-between rounded-xl border bg-card text-card-foreground shadow-md p-6 transition-all hover:shadow-xl hover:ring-2 hover:ring-primary/20 dark:bg-[#1e1e1e]"
      )}
    >
      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold">{title}</h2>
        {icon}
      </div>
      <div className="text-4xl font-bold">{value}</div>
    </div>
  )
}
