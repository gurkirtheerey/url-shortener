"use client";

import { ReferrerCount } from "@/lib/types";

export default function ReferrersTable({
  data,
}: {
  data: ReferrerCount[] | null;
}) {
  const referrers = data ?? [];

  if (referrers.length === 0) {
    return (
      <div className="border border-zinc-800 rounded-lg p-6">
        <p className="text-zinc-500 text-xs uppercase tracking-wider mb-4">
          Top Referrers
        </p>
        <div className="h-24 flex items-center justify-center">
          <p className="text-zinc-500 text-sm">No referrer data</p>
        </div>
      </div>
    );
  }

  const maxCount = referrers[0].count;

  return (
    <div className="border border-zinc-800 rounded-lg p-6">
      <p className="text-zinc-500 text-xs uppercase tracking-wider mb-4">
        Top Referrers
      </p>
      <div className="space-y-3">
        {referrers.map((referrer) => (
          <div key={referrer.referrer}>
            <div className="flex items-center justify-between text-sm mb-1">
              <span className="text-zinc-300 truncate mr-4">
                {referrer.referrer}
              </span>
              <span className="text-zinc-400 font-mono text-xs shrink-0">
                {referrer.count}
              </span>
            </div>
            <div className="h-1.5 bg-zinc-800 rounded-full overflow-hidden">
              <div
                className="h-full bg-amber-500/60 rounded-full transition-all"
                style={{
                  width: `${(referrer.count / maxCount) * 100}%`,
                }}
              />
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
