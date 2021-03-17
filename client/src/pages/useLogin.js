import { API_BASE_URL } from "@shared/consts";

import useFetch from "@shared/hooks/useFetch";

export default function useMultiTorrentSearch(searchText, pageNo, sources) {
  const url = new URL("/api/v1/login", API_BASE_URL);
  url.searchParams.append("pageNo", pageNo);
  return url.toString();

  return useFetch(url);
}
