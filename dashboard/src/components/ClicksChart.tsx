"use client";

import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";
import { DailyClicks } from "@/lib/types";

export default function ClicksChart({
  data,
}: {
  data: DailyClicks[] | null;
}) {
  const clicks = data ?? [];

  if (clicks.length === 0) {
    return (
      <div className="border border-zinc-800 rounded-lg p-6">
        <p className="text-zinc-500 text-xs uppercase tracking-wider mb-4">
          Clicks Over Time
        </p>
        <div className="h-48 flex items-center justify-center">
          <p className="text-zinc-500 text-sm">No clicks yet</p>
        </div>
      </div>
    );
  }

  return (
    <div className="border border-zinc-800 rounded-lg p-6">
      <p className="text-zinc-500 text-xs uppercase tracking-wider mb-4">
        Clicks Over Time
      </p>
      <div className="h-48">
        <ResponsiveContainer width="100%" height="100%">
          <AreaChart data={clicks}>
            <defs>
              <linearGradient id="clicksFill" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="#f59e0b" stopOpacity={0.3} />
                <stop offset="95%" stopColor="#f59e0b" stopOpacity={0} />
              </linearGradient>
            </defs>
            <XAxis
              dataKey="date"
              stroke="#3f3f46"
              tick={{ fill: "#71717a", fontSize: 11 }}
              tickLine={false}
              axisLine={false}
            />
            <YAxis
              stroke="#3f3f46"
              tick={{ fill: "#71717a", fontSize: 11 }}
              tickLine={false}
              axisLine={false}
              allowDecimals={false}
            />
            <Tooltip
              contentStyle={{
                backgroundColor: "#18181b",
                border: "1px solid #3f3f46",
                borderRadius: "8px",
                fontSize: "12px",
              }}
              labelStyle={{ color: "#a1a1aa" }}
              itemStyle={{ color: "#f59e0b" }}
            />
            <Area
              type="monotone"
              dataKey="clicks"
              stroke="#f59e0b"
              fill="url(#clicksFill)"
              strokeWidth={2}
            />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
