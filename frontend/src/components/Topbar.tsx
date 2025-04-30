"use client"

import { Moon, Sun } from "lucide-react"
import { useTheme } from "next-themes"

export function Topbar() {
  const { theme, setTheme } = useTheme()

  return (
    <header className="flex h-16 items-center justify-between border-b bg-background px-6 shadow-sm">
      <h1 className="text-2xl font-bold text-foreground">LoadBalancer Go</h1>

      {/* Dark Mode Toggle */}
      <button
        onClick={() => setTheme(theme === "dark" ? "light" : "dark")}
        className="flex items-center gap-2 rounded-md border px-3 py-1.5 text-sm font-medium shadow hover:bg-muted/50 transition"
      >
        {theme === "dark" ? <Sun className="h-5 w-5" /> : <Moon className="h-5 w-5" />}
        <span className="hidden sm:inline">
          {theme === "dark" ? "Light" : "Dark"}
        </span>
      </button>
    </header>
  )
}
