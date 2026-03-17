"use client";

import { FormEvent, useState } from "react";
import { useShortenUrl } from "@/lib/queries";
import CopyButton from "./CopyButton";

export default function ShortenForm() {
  const [url, setUrl] = useState("");
  const mutation = useShortenUrl();

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    const trimmed = url.trim();
    if (!trimmed) return;

    mutation.mutate(trimmed, {
      onSuccess: () => setUrl(""),
    });
  }

  return (
    <div>
      <form onSubmit={handleSubmit} className="flex gap-3">
        <input
          type="text"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="Paste a URL to shorten"
          className="flex-1 bg-zinc-900 border border-zinc-800 rounded-lg px-4 py-3 text-sm text-zinc-100 placeholder:text-zinc-500 focus:outline-none focus:border-amber-500/50 focus:ring-1 focus:ring-amber-500/25 transition-colors font-mono"
        />
        <button
          type="submit"
          disabled={mutation.isPending}
          className="bg-amber-500 hover:bg-amber-400 disabled:opacity-50 disabled:cursor-not-allowed text-zinc-950 font-medium text-sm px-6 py-3 rounded-lg transition-colors cursor-pointer"
        >
          {mutation.isPending ? "Shortening..." : "Shorten"}
        </button>
      </form>

      {mutation.isError && (
        <p className="text-red-400 text-xs mt-2">
          {mutation.error.message}
        </p>
      )}

      {mutation.isSuccess && (
        <div className="mt-3 flex items-center gap-3 bg-zinc-900 border border-zinc-800 rounded-lg px-4 py-3">
          <span className="text-amber-400 font-mono text-sm flex-1">
            {mutation.data.short_url}
          </span>
          <CopyButton text={mutation.data.short_url} />
        </div>
      )}
    </div>
  );
}
