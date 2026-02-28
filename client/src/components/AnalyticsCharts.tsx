"use client"

import { useEffect, useState } from "react"
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  Legend,
  ResponsiveContainer,
} from "recharts"
import api from "@/lib/api"

interface TimelineData {
  date: string
  count: number
}

interface AdminData {
  admin: string
  count: number
}

interface EntityData {
  entity: string
  count: number
}

const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042", "#A020F0"]

export default function AnalyticsCharts() {
  const [timeline, setTimeline] = useState<TimelineData[]>([])
  const [adminData, setAdminData] = useState<AdminData[]>([])
  const [entityData, setEntityData] = useState<EntityData[]>([])

  useEffect(() => {
    fetchAnalytics()
  }, [])

  const fetchAnalytics = async () => {
    const [timelineRes, adminRes, entityRes] = await Promise.all([
      api.get("/admin/audit-analytics/timeline"),
      api.get("/admin/audit-analytics/admin-actions"),
      api.get("/admin/audit-analytics/entity-distribution"),
    ])

    setTimeline(timelineRes.data)
    setAdminData(adminRes.data)
    setEntityData(entityRes.data)
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
      {/* Timeline Chart */}
      <div className="bg-white rounded-2xl shadow p-4">
        <h3 className="font-semibold mb-2">Actions Timeline</h3>
        <ResponsiveContainer width="100%" height={250}>
          <LineChart data={timeline}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" />
            <YAxis allowDecimals={false} />
            <Tooltip />
            <Line
              type="monotone"
              dataKey="count"
              stroke="#8884d8"
              strokeWidth={2}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>

      {/* Admin Activity Bar Chart */}
      <div className="bg-white rounded-2xl shadow p-4">
        <h3 className="font-semibold mb-2">Admin Activity</h3>
        <ResponsiveContainer width="100%" height={250}>
          <BarChart data={adminData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="admin" />
            <YAxis allowDecimals={false} />
            <Tooltip />
            <Bar dataKey="count" fill="#00C49F" />
          </BarChart>
        </ResponsiveContainer>
      </div>

      {/* Entity Distribution Pie Chart */}
      <div className="bg-white rounded-2xl shadow p-4">
        <h3 className="font-semibold mb-2">Entity Distribution</h3>
        <ResponsiveContainer width="100%" height={250}>
          <PieChart>
            <Pie
              data={entityData}
              dataKey="count"
              nameKey="entity"
              cx="50%"
              cy="50%"
              outerRadius={80}
              label
            >
              {entityData.map((_, index) => (
                <Cell key={index} fill={COLORS[index % COLORS.length]} />
              ))}
            </Pie>
            <Tooltip />
            <Legend />
          </PieChart>
        </ResponsiveContainer>
      </div>
    </div>
  )
}