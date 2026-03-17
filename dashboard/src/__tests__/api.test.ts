import { fetchUrls, shortenUrl, fetchStats, deleteUrl } from "@/lib/api";

const mockFetch = jest.fn();
global.fetch = mockFetch;

beforeEach(() => {
  mockFetch.mockClear();
});

describe("fetchUrls", () => {
  it("fetches and returns URL list", async () => {
    const urls = [
      {
        id: 1,
        short_code: "abc",
        original_url: "https://example.com",
        created_at: "2024-01-01T00:00:00Z",
      },
    ];
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(urls),
    });

    const result = await fetchUrls();
    expect(result).toEqual(urls);
    expect(mockFetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/urls"
    );
  });

  it("returns empty array when API returns null", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(null),
    });

    const result = await fetchUrls();
    expect(result).toEqual([]);
  });

  it("throws on non-ok response", async () => {
    mockFetch.mockResolvedValueOnce({ ok: false });
    await expect(fetchUrls()).rejects.toThrow("Failed to fetch URLs");
  });
});

describe("shortenUrl", () => {
  it("posts URL and returns shorten response", async () => {
    const response = {
      short_url: "http://localhost:8080/abc",
      short_code: "abc",
      original_url: "https://example.com",
    };
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(response),
    });

    const result = await shortenUrl("https://example.com");
    expect(result).toEqual(response);
    expect(mockFetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/shorten",
      {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ url: "https://example.com" }),
      }
    );
  });

  it("throws with error message from API on failure", async () => {
    mockFetch.mockResolvedValueOnce({
      ok: false,
      json: () => Promise.resolve({ error: "url is required" }),
    });

    await expect(shortenUrl("")).rejects.toThrow("url is required");
  });
});

describe("fetchStats", () => {
  it("fetches stats for a given code", async () => {
    const stats = {
      total_clicks: 10,
      daily_clicks: [{ date: "2024-01-01", clicks: 5 }],
      top_referrers: [{ referrer: "google.com", count: 3 }],
    };
    mockFetch.mockResolvedValueOnce({
      ok: true,
      json: () => Promise.resolve(stats),
    });

    const result = await fetchStats("abc");
    expect(result).toEqual(stats);
    expect(mockFetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/urls/abc/stats"
    );
  });

  it("throws on non-ok response", async () => {
    mockFetch.mockResolvedValueOnce({ ok: false });
    await expect(fetchStats("abc")).rejects.toThrow("Failed to fetch stats");
  });
});

describe("deleteUrl", () => {
  it("sends DELETE request for a given code", async () => {
    mockFetch.mockResolvedValueOnce({ ok: true });

    await deleteUrl("abc");
    expect(mockFetch).toHaveBeenCalledWith(
      "http://localhost:8080/api/urls/abc",
      { method: "DELETE" }
    );
  });

  it("throws on non-ok response", async () => {
    mockFetch.mockResolvedValueOnce({ ok: false });
    await expect(deleteUrl("abc")).rejects.toThrow("Failed to delete URL");
  });
});
