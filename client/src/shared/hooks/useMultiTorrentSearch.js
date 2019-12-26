import { useMemo } from "react";

import { API_BASE_URL } from "@shared/consts";

import useMultiFetch from "@shared/hooks/useMultiFetch";

export default function useMultiTorrentSearch(searchText, pageNo, sources) {
  const urls = useMemo(() => {
    if (!searchText || !sources.length) return [];

    return sources.map(source => {
      const url = new URL("/api/v1/search", API_BASE_URL);
      url.searchParams.append("query", searchText);
      url.searchParams.append("site", source);
      url.searchParams.append("pageNo", pageNo);
      return url.toString();
    });
  }, [searchText, sources, pageNo]);

  return useMultiFetch(urls).map(([result, status], i) => [{ ...result, source: sources[i] }, status]);
}
