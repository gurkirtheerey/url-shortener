"use client";

import { useUrls } from "@/lib/queries";
import LinkRow from "./LinkRow";
import EmptyState from "./EmptyState";

function LinkListSkeleton() {
  return (
    <div className="space-y-3">
      {[1, 2, 3].map((i) => (
        <div
          key={i}
          className="border border-zinc-800 rounded-lg px-4 py-3 animate-pulse"
        >
          <div className="h-4 bg-zinc-800 rounded w-24 mb-2" />
          <div className="h-3 bg-zinc-800 rounded w-64" />
        </div>
      ))}
    </div>
  );
}

export default function LinkList() {
  const { data: urls, isLoading, isError } = useUrls();

  if (isLoading) return <LinkListSkeleton />;

  if (isError) {
    return (
      <div className="border border-red-900/50 rounded-lg p-4">
        <p className="text-red-400 text-sm">
          Failed to load links. Is the API running?
        </p>
      </div>
    );
  }

  if (!urls || urls.length === 0) return <EmptyState />;

  return (
    <div className="space-y-2">
      {urls.map((url) => (
        <LinkRow key={url.id} url={url} />
      ))}
    </div>
  );
}
