"use client";

import { useEffect, useState } from "react";
import Navbar from "@/components/Navbar";
import api from "@/lib/api";
import AnalyticsCharts from "@/components/AnalyticsCharts";

interface AuditLog {
  id: string;
  actor: string;
  action: string;
  target: string;
  createdAt: string;
}

export default function AuditLogsPage() {
  const [logs, setLogs] = useState<AuditLog[]>([]);

  useEffect(() => {
    api.get("/admin/audit-logs").then((res) => {
      setLogs(res.data);
    });
  }, []);

  return (
    <div className="min-h-screen bg-neutral-100">
      <Navbar />
   <div className="min-h-screen bg-neutral-100 p-8">
      <h1 className="text-3xl font-bold mb-6">
        Enterprise Security Audit Dashboard
      </h1>

      {/* Analytics Charts */}
      <AnalyticsCharts />

      {/* Export Buttons + Audit Table */}
      ...
    </div>
      <div className="pt-28 max-w-6xl mx-auto px-4">
        <h1 className="text-3xl font-semibold mb-6">Audit Logs</h1>

        <div className="bg-white rounded-2xl shadow overflow-hidden">
          <table className="w-full text-sm">
            <thead className="bg-neutral-50 border-b">
              <tr>
                <th className="p-3 text-left">Actor</th>
                <th className="p-3 text-left">Action</th>
                <th className="p-3 text-left">Target</th>
                <th className="p-3 text-left">Date</th>
              </tr>
            </thead>
            <tbody>
              {logs.map((log) => (
                <tr key={log.id} className="border-b">
                  <td className="p-3">{log.actor}</td>
                  <td className="p-3">{log.action}</td>
                  <td className="p-3">{log.target}</td>
                  <td className="p-3">
                    {new Date(log.createdAt).toLocaleString()}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}