"use client";

import { useState } from "react";
import {
  AreaChart,
  Area,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  CartesianGrid,
} from "recharts";
import { DailyClicks } from "@/lib/types";
import { formatShortDate } from "@/lib/formatters";

type TimeRange = "7d" | "30d" | "all";

const ranges: { value: TimeRange; label: string }[] = [
  { value: "7d", label: "7d" },
  { value: "30d", label: "30d" },
  { value: "all", label: "All" },
];

function filterByRange(data: DailyClicks[], range: TimeRange): DailyClicks[] {
  if (range === "all") return data;
  const days = range === "7d" ? 7 : 30;
  const cutoff = new Date();
  cutoff.setDate(cutoff.getDate() - days);
  const cutoffStr = cutoff.toISOString().split("T")[0];
  return data.filter((d) => d.date >= cutoffStr);
}

function ChartTooltip({
  active,
  payload,
  label,
}: {
  active?: boolean;
  payload?: Array<{ value: number }>;
  label?: string;
}) {
  if (!active || !payload?.length || !label) return null;
  return (
    <div className="bg-zinc-900 border border-zinc-700 rounded-lg px-3 py-2 shadow-xl">
      <p className="text-zinc-400 text-[11px]">{formatShortDate(label)}</p>
      <p className="text-amber-400 font-mono text-sm font-semibold">
        {payload[0].value} clicks
      </p>
    </div>
  );
}

export default function ClicksChart({
  data,
}: {
  data: DailyClicks[] | null;
}) {
  const clicks = data ?? [];
  const [range, setRange] = useState<TimeRange>("30d");
  const filtered = filterByRange(clicks, range);

  return (
    <div
      className="border border-zinc-800 rounded-lg p-6 opacity-0 animate-fade-in-up"
      style={{ animationDelay: "240ms" }}
    >
      <div className="flex items-center justify-between mb-4">
        <p className="text-zinc-500 text-xs uppercase tracking-wider">
          Clicks Over Time
        </p>
        {clicks.length > 0 ? (
          <div className="flex gap-1">
            {ranges.map((r) => (
              <button
                key={r.value}
                onClick={() => setRange(r.value)}
                className={`px-2.5 py-1 text-[11px] rounded-md transition-colors cursor-pointer ${
                  range === r.value
                    ? "bg-zinc-700 text-zinc-100"
                    : "text-zinc-500 hover:text-zinc-300 hover:bg-zinc-800/50"
                }`}
              >
                {r.label}
              </button>
            ))}
          </div>
        ) : null}
      </div>

      {filtered.length === 0 ? (
        <div className="h-64 flex items-center justify-center">
          <p className="text-zinc-600 text-sm">
            {clicks.length === 0 ? "No clicks yet" : "No data in this range"}
          </p>
        </div>
      ) : (
        <div className="h-64">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart
              data={filtered}
              margin={{ top: 4, right: 4, left: -20, bottom: 0 }}
            >
              <defs>
                <linearGradient id="clicksFill" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="0%" stopColor="#f59e0b" stopOpacity={0.15} />
                  <stop offset="100%" stopColor="#f59e0b" stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid
                strokeDasharray="3 3"
                stroke="#27272a"
                vertical={false}
              />
              <XAxis
                dataKey="date"
                stroke="none"
                tick={{ fill: "#52525b", fontSize: 11 }}
                tickLine={false}
                tickFormatter={(value: string) => formatShortDate(value)}
                interval="preserveStartEnd"
              />
              <YAxis
                stroke="none"
                tick={{ fill: "#52525b", fontSize: 11 }}
                tickLine={false}
                allowDecimals={false}
              />
              <Tooltip
                content={<ChartTooltip />}
                cursor={{ stroke: "#3f3f46" }}
              />
              <Area
                type="monotone"
                dataKey="clicks"
                stroke="#f59e0b"
                fill="url(#clicksFill)"
                strokeWidth={2}
                dot={false}
                activeDot={{
                  r: 4,
                  fill: "#f59e0b",
                  stroke: "#18181b",
                  strokeWidth: 2,
                }}
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      )}
    </div>
  );
}
