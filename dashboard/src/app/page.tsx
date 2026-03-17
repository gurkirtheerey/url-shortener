import ShortenForm from "@/components/ShortenForm";
import LinkList from "@/components/LinkList";

export default function Home() {
  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Shortener</h1>
        <p className="text-zinc-500 text-sm mt-1">
          Shorten URLs and track clicks.
        </p>
      </div>
      <ShortenForm />
      <LinkList />
    </div>
  );
}
