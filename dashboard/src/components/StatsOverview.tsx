"use client";

import { UrlStats } from "@/lib/types";
import { formatNumber, formatDate } from "@/lib/formatters";

interface StatsOverviewProps {
  stats: UrlStats;
  originalUrl: string;
  shortUrl: string;
  createdAt: string;
}

export default function StatsOverview({
  stats,
  originalUrl,
  shortUrl,
  createdAt,
}: StatsOverviewProps) {
  return (
    <div className="border border-zinc-800 rounded-lg p-6">
      <div className="mb-6">
        <p className="text-zinc-500 text-xs uppercase tracking-wider mb-1">
          Total Clicks
        </p>
        <p className="text-5xl font-mono text-amber-400 font-bold">
          {formatNumber(stats.total_clicks)}
        </p>
      </div>

      <div className="space-y-3 text-sm">
        <div>
          <p className="text-zinc-500 text-xs">Short URL</p>
          <p className="text-zinc-200 font-mono">{shortUrl}</p>
        </div>
        <div>
          <p className="text-zinc-500 text-xs">Original URL</p>
          <p className="text-zinc-400 break-all">{originalUrl}</p>
        </div>
        <div>
          <p className="text-zinc-500 text-xs">Created</p>
          <p className="text-zinc-400">
            {createdAt ? formatDate(createdAt) : ""}
          </p>
        </div>
      </div>
    </div>
  );
}
