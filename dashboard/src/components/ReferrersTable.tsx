"use client";

import { ReferrerCount } from "@/lib/types";

export default function ReferrersTable({
  data,
}: {
  data: ReferrerCount[] | null;
}) {
  const referrers = data ?? [];
  const totalReferrerClicks = referrers.reduce(
    (sum, r) => sum + r.count,
    0
  );

  return (
    <div
      className="border border-zinc-800 rounded-lg p-6 opacity-0 animate-[fadeInUp_0.4s_ease-out_forwards]"
      style={{ animationDelay: "320ms" }}
    >
      <p className="text-zinc-500 text-xs uppercase tracking-wider mb-4">
        Top Referrers
      </p>

      {referrers.length === 0 ? (
        <div className="h-24 flex items-center justify-center">
          <p className="text-zinc-600 text-sm">No referrer data</p>
        </div>
      ) : (
        <div className="space-y-3">
          {referrers.map((referrer) => {
            const percentage =
              totalReferrerClicks > 0
                ? Math.round((referrer.count / totalReferrerClicks) * 100)
                : 0;

            return (
              <div key={referrer.referrer} className="group">
                <div className="flex items-center justify-between text-sm mb-1.5">
                  <span className="text-zinc-300 truncate mr-4 group-hover:text-zinc-100 transition-colors">
                    {referrer.referrer}
                  </span>
                  <div className="flex items-center gap-2.5 shrink-0">
                    <span className="text-zinc-600 text-xs tabular-nums">
                      {percentage}%
                    </span>
                    <span className="text-zinc-400 font-mono text-xs w-8 text-right tabular-nums">
                      {referrer.count}
                    </span>
                  </div>
                </div>
                <div className="h-1 bg-zinc-800/80 rounded-full overflow-hidden">
                  <div
                    className="h-full bg-amber-500/50 rounded-full transition-all duration-500 ease-out group-hover:bg-amber-400/60"
                    style={{ width: `${percentage}%` }}
                  />
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
