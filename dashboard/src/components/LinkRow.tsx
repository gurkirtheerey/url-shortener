"use client";

import { useState } from "react";
import Link from "next/link";
import { ShortenedUrl } from "@/lib/types";
import { API_URL } from "@/lib/api";
import { useDeleteUrl } from "@/lib/queries";
import { relativeTime, truncateUrl, formatNumber } from "@/lib/formatters";
import CopyButton from "./CopyButton";

export default function LinkRow({ url }: { url: ShortenedUrl }) {
  const [confirming, setConfirming] = useState(false);
  const deleteMutation = useDeleteUrl();
  const shortUrl = `${API_URL}/${url.short_code}`;
  const clickCount = url.click_count ?? 0;

  function handleDelete(e: React.MouseEvent) {
    e.preventDefault();
    e.stopPropagation();
    if (!confirming) {
      setConfirming(true);
      return;
    }
    deleteMutation.mutate(url.short_code);
    setConfirming(false);
  }

  return (
    <Link
      href={`/link/${url.short_code}`}
      className="group block border border-zinc-800 rounded-lg p-4 hover:border-zinc-700 transition-colors"
    >
      <div className="flex items-start justify-between gap-4">
        <div className="min-w-0 flex-1">
          <span className="font-mono text-sm text-amber-400">
            /{url.short_code}
          </span>
          <p className="text-zinc-500 text-xs mt-1 truncate">
            {truncateUrl(url.original_url, 80)}
          </p>
        </div>
        <CopyButton
          text={shortUrl}
          className="text-[11px] border border-zinc-700 hover:border-zinc-500 rounded-md px-2.5 py-1 text-zinc-400 hover:text-zinc-200 transition-colors cursor-pointer shrink-0"
        />
      </div>

      <div className="flex items-center justify-between mt-3 pt-3 border-t border-zinc-800/60">
        <div className="flex items-center gap-1.5 text-xs text-zinc-600">
          <span>{relativeTime(url.created_at)}</span>
          <span>·</span>
          <span>
            {formatNumber(clickCount)}{" "}
            {clickCount === 1 ? "click" : "clicks"}
          </span>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-xs text-zinc-500 group-hover:text-amber-400 transition-colors">
            Stats
          </span>
          <button
            onClick={handleDelete}
            onBlur={() => setConfirming(false)}
            disabled={deleteMutation.isPending}
            className={`text-xs transition-colors cursor-pointer ${
              confirming
                ? "text-red-400 hover:text-red-300"
                : "text-zinc-500 hover:text-red-400"
            }`}
          >
            {confirming ? "Confirm?" : "Delete"}
          </button>
        </div>
      </div>
    </Link>
  );
}
