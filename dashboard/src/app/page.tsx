import ShortenForm from "@/components/ShortenForm";
import LinkList from "@/components/LinkList";

export default function Home() {
  return (
    <div className="space-y-10">
      <div>
        <h1 className="text-lg font-semibold tracking-tight">linksmith</h1>
      </div>
      <ShortenForm />
      <LinkList />
    </div>
  );
}
