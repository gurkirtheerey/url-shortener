export default function EmptyState() {
  return (
    <div className="border border-dashed border-zinc-800 rounded-lg py-16 text-center">
      <p className="text-zinc-400 text-sm">No links yet</p>
      <p className="text-zinc-600 text-xs mt-1.5">
        Paste a URL above to create your first short link.
      </p>
    </div>
  );
}
