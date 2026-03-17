import { relativeTime, truncateUrl, formatNumber, formatDate } from "@/lib/formatters";

describe("relativeTime", () => {
  it("returns 'just now' for dates less than a minute ago", () => {
    const now = new Date().toISOString();
    expect(relativeTime(now)).toBe("just now");
  });

  it("returns minutes for dates less than an hour ago", () => {
    const fiveMinutesAgo = new Date(Date.now() - 5 * 60 * 1000).toISOString();
    expect(relativeTime(fiveMinutesAgo)).toBe("5m ago");
  });

  it("returns hours for dates less than a day ago", () => {
    const threeHoursAgo = new Date(
      Date.now() - 3 * 60 * 60 * 1000
    ).toISOString();
    expect(relativeTime(threeHoursAgo)).toBe("3h ago");
  });

  it("returns days for dates less than 30 days ago", () => {
    const twoDaysAgo = new Date(
      Date.now() - 2 * 24 * 60 * 60 * 1000
    ).toISOString();
    expect(relativeTime(twoDaysAgo)).toBe("2d ago");
  });

  it("returns formatted date for dates older than 30 days", () => {
    const old = new Date("2023-01-15T00:00:00Z").toISOString();
    const result = relativeTime(old);
    expect(result).toMatch(/1\/15\/2023|15\/1\/2023|2023/);
  });
});

describe("truncateUrl", () => {
  it("returns url unchanged when shorter than max length", () => {
    expect(truncateUrl("https://example.com")).toBe("https://example.com");
  });

  it("truncates url and adds ellipsis when longer than max length", () => {
    const longUrl = "https://example.com/" + "a".repeat(100);
    const result = truncateUrl(longUrl, 30);
    expect(result.length).toBe(30);
    expect(result.endsWith("\u2026")).toBe(true);
  });

  it("uses default max length of 50", () => {
    const url = "https://example.com/" + "a".repeat(100);
    const result = truncateUrl(url);
    expect(result.length).toBe(50);
  });
});

describe("formatNumber", () => {
  it("returns number as-is when less than 1000", () => {
    expect(formatNumber(42)).toBe("42");
    expect(formatNumber(999)).toBe("999");
  });

  it("formats thousands", () => {
    expect(formatNumber(1500)).toBe("1.5K");
    expect(formatNumber(10000)).toBe("10.0K");
  });

  it("formats millions", () => {
    expect(formatNumber(2500000)).toBe("2.5M");
  });

  it("handles zero", () => {
    expect(formatNumber(0)).toBe("0");
  });
});

describe("formatDate", () => {
  it("formats ISO date string to readable format", () => {
    const result = formatDate("2024-03-15T10:30:00Z");
    expect(result).toMatch(/Mar 15, 2024/);
  });
});
