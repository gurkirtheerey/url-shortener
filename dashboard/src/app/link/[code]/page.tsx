"use client";

import { use } from "react";
import Link from "next/link";
import dynamic from "next/dynamic";
import { API_URL } from "@/lib/api";
import { useStats, useUrls } from "@/lib/queries";
import StatsOverview from "@/components/StatsOverview";
import ReferrersTable from "@/components/ReferrersTable";

const ClicksChart = dynamic(() => import("@/components/ClicksChart"), {
  ssr: false,
  loading: () => (
    <div className="border border-zinc-800 rounded-lg p-6">
      <div className="h-48 animate-pulse bg-zinc-800/50 rounded" />
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

  return (
    <div className="space-y-6">
      <Link
        href="/"
        className="text-zinc-500 hover:text-zinc-300 text-sm transition-colors inline-block"
      >
        &larr; Back
      </Link>

      {statsLoading && (
        <div className="space-y-6">
          <div className="border border-zinc-800 rounded-lg p-6 animate-pulse">
            <div className="h-4 bg-zinc-800 rounded w-32 mb-4" />
            <div className="h-12 bg-zinc-800 rounded w-24" />
          </div>
          <div className="border border-zinc-800 rounded-lg p-6">
            <div className="h-48 animate-pulse bg-zinc-800/50 rounded" />
          </div>
        </div>
      )}

      {isError && (
        <div className="border border-red-900/50 rounded-lg p-4">
          <p className="text-red-400 text-sm">
            Failed to load stats. Is the API running?
          </p>
        </div>
      )}

      {stats && (
        <>
          <StatsOverview
            stats={stats}
            originalUrl={urlInfo?.original_url ?? ""}
            shortUrl={shortUrl}
            createdAt={urlInfo?.created_at ?? ""}
          />
          <ClicksChart data={stats.daily_clicks} />
          <ReferrersTable data={stats.top_referrers} />
        </>
      )}
    </div>
  );
}
