"use client";

import { use } from "react";
import Link from "next/link";
import dynamic from "next/dynamic";
import { API_URL } from "@/lib/api";
import { useStats, useUrls } from "@/lib/queries";
import { truncateUrl, formatDate } from "@/lib/formatters";
import StatsOverview from "@/components/StatsOverview";
import ReferrersTable from "@/components/ReferrersTable";
import CopyButton from "@/components/CopyButton";

const ClicksChart = dynamic(() => import("@/components/ClicksChart"), {
  ssr: false,
  loading: () => (
    <div className="border border-zinc-800 rounded-lg p-6">
      <div className="h-64 animate-pulse bg-zinc-800/30 rounded" />
    </div>
  ),
});

export default function LinkDetailPage({
  params,
}: {
  params: Promise<{ code: string }>;
}) {
  const { code } = use(params);
  const { data: stats, isLoading: statsLoading, isError } = useStats(code);
  const { data: urls } = useUrls();

  const shortUrl = `${API_URL}/${code}`;
  const urlInfo = urls?.find((u) => u.short_code === code);

  if (statsLoading) {
    return (
      <div className="space-y-6">
        <Link
          href="/"
          className="text-zinc-500 hover:text-zinc-300 text-sm transition-colors inline-block"
        >
          &larr; Back
        </Link>
        <div className="space-y-3">
          <div className="h-6 bg-zinc-800/50 rounded w-24 animate-pulse" />
          <div className="h-4 bg-zinc-800/30 rounded w-72 animate-pulse" />
          <div className="h-3 bg-zinc-800/20 rounded w-32 animate-pulse" />
        </div>
        <div className="grid grid-cols-3 gap-3">
          {[0, 1, 2].map((i) => (
            <div key={i} className="border border-zinc-800 rounded-lg p-4">
              <div className="h-3 bg-zinc-800/50 rounded w-16 mb-3 animate-pulse" />
              <div className="h-7 bg-zinc-800/50 rounded w-14 animate-pulse" />
            </div>
          ))}
        </div>
        <div className="border border-zinc-800 rounded-lg p-6">
          <div className="h-64 animate-pulse bg-zinc-800/30 rounded" />
        </div>
      </div>
    );
  }

  if (isError) {
    return (
      <div className="space-y-6">
        <Link
          href="/"
          className="text-zinc-500 hover:text-zinc-300 text-sm transition-colors inline-block"
        >
          &larr; Back
        </Link>
        <div className="border border-red-900/30 rounded-lg p-6 text-center">
          <p className="text-red-400 text-sm">
            Failed to load stats. Is the API running?
          </p>
        </div>
      </div>
    );
  }

  if (!stats) return null;

  return (
    <div className="space-y-6">
      <Link
        href="/"
        className="text-zinc-500 hover:text-zinc-300 text-sm transition-colors inline-block"
      >
        &larr; Back
      </Link>

      <div className="opacity-0 animate-fade-in-up">
        <div className="flex items-center gap-3 mb-2">
          <h1 className="text-lg font-mono text-amber-400 font-semibold">
            /{code}
          </h1>
          <CopyButton text={shortUrl} />
        </div>
        {urlInfo ? (
          <div className="space-y-1">
            <a
              href={urlInfo.original_url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-zinc-400 text-sm hover:text-zinc-200 transition-colors break-all inline-block"
              title={urlInfo.original_url}
            >
              {truncateUrl(urlInfo.original_url, 80)}
              <span className="text-zinc-600 ml-1">{"\u2197"}</span>
            </a>
            <p className="text-zinc-600 text-xs">
              Created {formatDate(urlInfo.created_at)}
            </p>
          </div>
        ) : null}
      </div>

      <StatsOverview stats={stats} />
      <ClicksChart data={stats.daily_clicks} />
      <ReferrersTable data={stats.top_referrers} />
    </div>
  );
}
