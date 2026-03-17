"use client";

import { useState } from "react";
import Link from "next/link";
import { ShortenedUrl } from "@/lib/types";
import { API_URL } from "@/lib/api";
import { useDeleteUrl } from "@/lib/queries";
import { relativeTime, truncateUrl } from "@/lib/formatters";
import CopyButton from "./CopyButton";

export default function LinkRow({ url }: { url: ShortenedUrl }) {
  const [confirming, setConfirming] = useState(false);
  const deleteMutation = useDeleteUrl();
  const shortUrl = `${API_URL}/${url.short_code}`;

  function handleDelete() {
    if (!confirming) {
      setConfirming(true);
      return;
    }
    deleteMutation.mutate(url.short_code);
    setConfirming(false);
  }

  return (
    <div className="group flex items-center justify-between gap-4 border border-zinc-800 rounded-lg px-4 py-3 hover:border-zinc-700 transition-colors">
      <div className="min-w-0 flex-1">
        <div className="flex items-center gap-2">
          <span className="font-mono text-sm text-amber-400">
            /{url.short_code}
          </span>
          <span className="text-zinc-600 text-xs">
            {relativeTime(url.created_at)}
          </span>
        </div>
        <p className="text-zinc-400 text-xs mt-0.5 truncate">
          {truncateUrl(url.original_url, 80)}
        </p>
      </div>

      <div className="flex items-center gap-3 shrink-0">
        <CopyButton text={shortUrl} />
        <Link
          href={`/link/${url.short_code}`}
          className="text-xs text-zinc-400 hover:text-amber-400 transition-colors"
        >
          Stats
        </Link>
        <button
          onClick={handleDelete}
          onBlur={() => setConfirming(false)}
          disabled={deleteMutation.isPending}
          className={`text-xs transition-colors cursor-pointer ${
            confirming
              ? "text-red-400 hover:text-red-300"
              : "text-zinc-400 hover:text-red-400"
          }`}
        >
          {confirming ? "Confirm?" : "Delete"}
        </button>
      </div>
    </div>
  );
}
