import React, { useState, useMemo } from "react";

import Search from "@shared/components/Search";
import TorrentMetaList from "@shared/components/TorrentMetaList";
import { API_BASE_URL } from "@shared/consts";
import useMultiFetch from "@shared/hooks/useMultiFetch";

import "./App.scss";

function useTorrentSearch(searchText, pageNo, sources) {
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

function App() {
  const [searchText, setSearchText] = useState("");
  const [selectedSources, setSelectedSources] = useState(["tpb", "kickasstorrent", "1337x"]);
  const [currentPageNo, setCurrentPageNo] = useState(1);
  const searchResult = useTorrentSearch(searchText, currentPageNo, selectedSources);
  console.log("searchResult :", searchResult);
  return (
    <div className="app">
      <Search onSearch={setSearchText} />
      <section>
        {searchResult.length > 0 && (
          <TorrentMetaList searchResult={searchResult} currentPageNo={currentPageNo} onPage={setCurrentPageNo} />
        )}
      </section>
    </div>
  );
}

export default App;
