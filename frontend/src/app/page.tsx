"use client"

import { useEffect } from "react"

import { Button } from "@/components/ui/button"
import { Server } from "lucide-react"
import api from "../../lib/axios"

export default function Home() {
  useEffect(() => {
    api.get("/status").then((res) => {
      console.log("Load Balancer Status:", res.data)
    }).catch((err) => {
      console.error("API Error:", err)
    })
  }, [])

  return (
    <></>
  )
}
