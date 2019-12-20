import React, { useState, useMemo } from "react";

import Search from "@shared/components/Search";
import TorrentMetaList from "@shared/components/TorrentMetaList";
import { API_BASE_URL } from "@shared/consts";
import useMultiFetch from "@shared/hooks/useMultiFetch";

import "./App.scss";

function useTorrentSearch(searchText, sources) {
  const urls = useMemo(() => {
    if (!searchText || !sources.length) return [];

    return sources.map(source => {
      const url = new URL("/api/v1/search", API_BASE_URL);
      url.searchParams.append("query", searchText);
      url.searchParams.append("site", source);
      return url.toString();
    });
  }, [searchText, sources]);

  return useMultiFetch(urls).map((ret, i) => [sources[i], ret]);
}

function App() {
  const [searchText, setSearchText] = useState("");
  const [selectedSources, setSelectedSources] = useState(["tpb", "kickasstorrent", "1337x"]);
  const searchResult = useTorrentSearch(searchText, selectedSources);
  console.log("searchResult :", searchResult);
  return (
    <div className="app">
      <Search onSearch={setSearchText} />
      <section>
        <TorrentMetaList searchResult={searchResult} />
      </section>
    </div>
  );
}

export default App;
