export interface ShortenedUrl {
  id: number;
  short_code: string;
  original_url: string;
  created_at: string;
  click_count: number;
}

export interface ShortenResponse {
  short_url: string;
  short_code: string;
  original_url: string;
}

export interface DailyClicks {
  date: string;
  clicks: number;
}

export interface ReferrerCount {
  referrer: string;
  count: number;
}

export interface UrlStats {
  total_clicks: number;
  daily_clicks: DailyClicks[] | null;
  top_referrers: ReferrerCount[] | null;
}
