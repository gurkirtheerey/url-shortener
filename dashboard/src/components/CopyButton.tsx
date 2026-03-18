"use client";

import { useState } from "react";

const defaultClassName =
  "text-xs text-zinc-400 hover:text-amber-400 transition-colors cursor-pointer";

export default function CopyButton({
  text,
  className,
}: {
  text: string;
  className?: string;
}) {
  const [copied, setCopied] = useState(false);

  async function handleCopy(e: React.MouseEvent) {
    e.preventDefault();
    e.stopPropagation();
    try {
      await navigator.clipboard.writeText(text);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {
      // Clipboard API can fail if the user denies permission
    }
  }

  return (
    <button
      onClick={handleCopy}
      className={className ?? defaultClassName}
      title="Copy to clipboard"
    >
      {copied ? "Copied!" : "Copy"}
    </button>
  );
}
