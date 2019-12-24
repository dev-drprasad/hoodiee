import React, { useState, useMemo, useContext, useCallback, useReducer, useEffect } from "react";

import Search from "@shared/components/Search";
import TorrentMetaList from "@shared/components/TorrentMetaList";
import StatusHandler from "@shared/components/StatusHandler";
import { mergeStatuses } from "@shared/utils";

import "./App.scss";

const initState = {
  searchText: "",
  pageNo: 1,
  sources: [],
  // Adding new property which is not used in TorrentMetaList may cause problem
};

const sources = [
  { id: "tpb", name: "The Pirate Bay" },
  { id: "kickasstorrent", name: "KickAssTorrent" },
  { id: "1337x", name: "1337x" },
];

function reducer(state, action) {
  switch (action.type) {
    case "SET_SEARCH_TEXT":
      return { ...state, searchText: action.payload };
    case "SET_PAGE_NO":
      return { ...state, pageNo: action.payload };
    case "SET_SOURCES":
      return { ...state, sources: action.payload };
    case "SET_STATE":
      return { ...state, ...action.payload };
    default:
      return state;
  }
}

function App() {
  const [state, dispatch] = useReducer(reducer, initState);
  const [maxPages, setMaxPages] = useState([]);

  const updateMaxPage = maxPage => {
    if (!maxPages.find(o => o.source === maxPage.source)) {
      setMaxPages([...maxPages, maxPage]);
    }
  };

  const maxPageNo = maxPages.length > 0 ? Math.max(...maxPages.map(o => o.pages)) : 0;
  console.log("state :", state);
  console.log("maxPageNo :", maxPageNo);

  const params = useMemo(() => state.sources.map(s => ({ source: s, pageNo: state.pageNo, query: state.searchText })), [
    state,
  ]);

  useEffect(() => {
    setMaxPages([]);
  }, [params]);

  console.log("params :", params);

  return (
    <div className="app">
      <Search onSearch={payload => dispatch({ type: "SET_STATE", payload })} sources={sources} />
      {state.searchText && (
        <section>
          {params.map(params => (
            <TorrentMetaList key={params.source} params={params} setPageNo={updateMaxPage} />
          ))}
        </section>
      )}
      {maxPageNo > 0 && (
        <ul className="pagination">
          {Array.from({ length: maxPageNo }, (_, i) => i + 1).map(i => (
            <li key={i} className={`pagination-item ${state.pageNo === 1 ? "pagination-item-active" : ""}`}>
              {<button onClick={() => dispatch({ type: "SET_PAGE_NO", payload: i })}>{i}</button>}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

export default App;
