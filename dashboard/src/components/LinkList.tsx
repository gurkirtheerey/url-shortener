"use client";

import { useUrls } from "@/lib/queries";
import LinkRow from "./LinkRow";
import EmptyState from "./EmptyState";

function LinkListSkeleton() {
  return (
    <div>
      <div className="flex items-center gap-2 mb-4">
        <div className="h-4 bg-zinc-800/50 rounded w-20 animate-pulse" />
      </div>
      <div className="space-y-2">
        {[0, 1, 2].map((i) => (
          <div
            key={i}
            className="border border-zinc-800 rounded-lg p-4 animate-pulse"
          >
            <div className="h-4 bg-zinc-800/50 rounded w-16 mb-2" />
            <div className="h-3 bg-zinc-800/30 rounded w-64" />
            <div className="mt-3 pt-3 border-t border-zinc-800/60 flex justify-between">
              <div className="h-3 bg-zinc-800/30 rounded w-28" />
              <div className="h-3 bg-zinc-800/30 rounded w-20" />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

export default function LinkList() {
  const { data: urls, isLoading, isError } = useUrls();

  if (isLoading) return <LinkListSkeleton />;

  if (isError) {
    return (
      <div className="border border-red-900/30 rounded-lg p-4">
        <p className="text-red-400 text-sm">
          Failed to load links. Is the API running?
        </p>
      </div>
    );
  }

  if (!urls || urls.length === 0) return <EmptyState />;

  return (
    <div>
      <div className="flex items-center gap-2 mb-4">
        <h2 className="text-sm font-medium text-zinc-300">Your Links</h2>
        <span className="text-xs text-zinc-600 tabular-nums">{urls.length}</span>
      </div>
      <div className="space-y-2">
        {urls.map((url, index) => (
          <div
            key={url.id}
            className="opacity-0 animate-fade-in-up"
            style={{ animationDelay: `${index * 60}ms` }}
          >
            <LinkRow url={url} />
          </div>
        ))}
      </div>
    </div>
  );
}
