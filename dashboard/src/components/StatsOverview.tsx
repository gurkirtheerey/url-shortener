"use client";

import { UrlStats, DailyClicks } from "@/lib/types";
import { formatNumber, formatShortDate } from "@/lib/formatters";

interface DerivedMetrics {
  avgPerDay: number;
  bestDay: DailyClicks | null;
  trend: number | null;
}

function deriveMetrics(
  dailyClicks: DailyClicks[] | null,
  totalClicks: number
): DerivedMetrics {
  if (!dailyClicks || dailyClicks.length === 0) {
    return { avgPerDay: 0, bestDay: null, trend: null };
  }

  const avgPerDay = totalClicks / dailyClicks.length;

  const bestDay = dailyClicks.reduce((best, day) =>
    day.clicks > best.clicks ? day : best
  );

  const sorted = [...dailyClicks].sort((a, b) =>
    a.date.localeCompare(b.date)
  );
  const last7 = sorted.slice(-7);
  const prev7 = sorted.slice(-14, -7);

  if (prev7.length === 0) {
    return { avgPerDay, bestDay, trend: null };
  }

  const last7Total = last7.reduce((sum, d) => sum + d.clicks, 0);
  const prev7Total = prev7.reduce((sum, d) => sum + d.clicks, 0);

  const trend =
    prev7Total === 0
      ? last7Total > 0
        ? 100
        : 0
      : Math.round(((last7Total - prev7Total) / prev7Total) * 100);

  return { avgPerDay, bestDay, trend };
}

export default function StatsOverview({ stats }: { stats: UrlStats }) {
  const { avgPerDay, bestDay, trend } = deriveMetrics(
    stats.daily_clicks,
    stats.total_clicks
  );

  return (
    <div className="grid grid-cols-3 gap-3">
      <div
        className="border border-amber-900/30 rounded-lg p-4 bg-gradient-to-br from-amber-500/[0.04] to-transparent hover:border-amber-900/50 transition-colors opacity-0 animate-fade-in-up"
      >
        <p className="text-zinc-500 text-[10px] uppercase tracking-wider mb-2">
          Total Clicks
        </p>
        <p className="text-2xl font-mono text-amber-400 font-semibold">
          {formatNumber(stats.total_clicks)}
        </p>
        {trend !== null ? (
          <p
            className={`text-xs mt-1.5 ${
              trend >= 0 ? "text-emerald-400" : "text-red-400"
            }`}
          >
            {trend >= 0 ? "\u2191" : "\u2193"} {Math.abs(trend)}% vs last week
          </p>
        ) : null}
      </div>

      <div
        className="border border-zinc-800 rounded-lg p-4 hover:border-zinc-700 transition-colors opacity-0 animate-fade-in-up"
        style={{ animationDelay: "80ms" }}
      >
        <p className="text-zinc-500 text-[10px] uppercase tracking-wider mb-2">
          Avg / Day
        </p>
        <p className="text-2xl font-mono text-zinc-100 font-semibold">
          {avgPerDay > 0 ? avgPerDay.toFixed(1) : "0"}
        </p>
      </div>

      <div
        className="border border-zinc-800 rounded-lg p-4 hover:border-zinc-700 transition-colors opacity-0 animate-fade-in-up"
        style={{ animationDelay: "160ms" }}
      >
        <p className="text-zinc-500 text-[10px] uppercase tracking-wider mb-2">
          Best Day
        </p>
        {bestDay ? (
          <>
            <p className="text-2xl font-mono text-zinc-100 font-semibold">
              {formatNumber(bestDay.clicks)}
            </p>
            <p className="text-zinc-500 text-xs mt-1.5">
              {formatShortDate(bestDay.date)}
            </p>
          </>
        ) : (
          <p className="text-2xl font-mono text-zinc-600">--</p>
        )}
      </div>
    </div>
  );
}
