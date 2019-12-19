import React, { useState, useMemo } from "react";
import Search from "@shared/components/Search";
import TorrentMetaList from "@shared/components/TorrentMetaList";

import { useFetch } from "./shared/hooks";
import { API_BASE_URL } from "./shared/consts";

import "./App.scss";
import StatusHandler from "@shared/components/StatusHandler";

const defaultSearchResultValue = [];
function useTorrentSearch(searchText) {
  const url = useMemo(() => {
    if (!searchText) return;
    const url = new URL("/api/v1/search", API_BASE_URL);
    url.searchParams.append("query", searchText);
    url.searchParams.append("site", "kickasstorrent");
    return url.toString();
  }, [searchText]);
  return useFetch(url, undefined, defaultSearchResultValue);
}

function App() {
  const [searchText, setSearchText] = useState("");
  const [searchResult, searchResultStatus] = useTorrentSearch(searchText);
  console.log("searchResult :", searchResult);
  return (
    <div className="app">
      <Search onSearch={setSearchText} />
      <section>
        <StatusHandler status={searchResultStatus}>
          {() => <TorrentMetaList items={searchResult} site="kickasstorrent" />}
        </StatusHandler>
      </section>
    </div>
  );
}

export default App;
