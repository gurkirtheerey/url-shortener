export default function Footer() {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="border-t border-zinc-800/60 mt-20 pt-6 pb-8">
      <div className="flex items-center justify-between text-[11px] text-zinc-500">
        <div className="flex items-center gap-1.5">
          <span className="font-mono text-amber-500/70">&mdash;</span>
          <span className="tracking-wide">
            linksmith
          </span>
          <span className="text-zinc-700">/</span>
          <span className="font-mono text-zinc-600">{currentYear}</span>
        </div>
        <a
          href="https://github.com/gurkirtheerey/url-shortener"
          target="_blank"
          rel="noopener noreferrer"
          className="font-mono text-zinc-600 transition-colors hover:text-amber-500/80"
        >
          src
        </a>
      </div>
    </footer>
  );
}
