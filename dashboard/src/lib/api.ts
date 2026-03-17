import { ShortenedUrl, ShortenResponse, UrlStats } from "./types";

export const API_URL =
  process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export async function fetchUrls(): Promise<ShortenedUrl[]> {
  const response = await fetch(`${API_URL}/api/urls`);
  if (!response.ok) throw new Error("Failed to fetch URLs");
  const data = await response.json();
  return data ?? [];
}

export async function shortenUrl(url: string): Promise<ShortenResponse> {
  const response = await fetch(`${API_URL}/api/shorten`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ url }),
  });
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error ?? "Failed to shorten URL");
  }
  return response.json();
}

export async function fetchStats(code: string): Promise<UrlStats> {
  const response = await fetch(`${API_URL}/api/urls/${code}/stats`);
  if (!response.ok) throw new Error("Failed to fetch stats");
  return response.json();
}

export async function deleteUrl(code: string): Promise<void> {
  const response = await fetch(`${API_URL}/api/urls/${code}`, {
    method: "DELETE",
  });
  if (!response.ok) throw new Error("Failed to delete URL");
}
