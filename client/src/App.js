import React, { useState, useMemo, useReducer, useEffect } from "react";
import Pagination from "rc-pagination";
import { Router } from "@reach/router";
import Search from "@shared/components/Search";
import TorrentMetaList from "@shared/components/TorrentMetaList";
import StatusHandler from "@shared/components/StatusHandler";

import { SourceMapProvider } from "@shared/contexts/SourceMap";
import { useFetch } from "@shared/hooks";
import "rc-pagination/assets/index.css";

import "./App.scss";
import Login from "pages/Login";
import ProtectedComponent from "@shared/components/ProtectedComponent";
import AuthContext from "@shared/contexts/Auth";

const initState = {
  searchText: "",
  pageNo: 1,
  sources: [],
  // Adding new property which is not used in TorrentMetaList may cause problem
};

const defaultSources = {};
function useSources() {
  return useFetch("/api/v1/sources", undefined, defaultSources);
}

function reducer(state, action) {
  switch (action.type) {
    case "SET_SEARCH_TEXT":
      return { ...state, searchText: action.payload };
    case "SET_PAGE_NO":
      return { ...state, pageNo: action.payload };
    case "SET_SOURCES":
      return { ...state, sources: action.payload };
    case "SET_SEARCH_PARAMS":
      return { ...state, ...action.payload, pageNo: 1 };
    default:
      return state;
  }
}

function Dashboard() {
  const [state, dispatch] = useReducer(reducer, initState);
  const [sources, sourcesStatus] = useSources();
  const [maxPages, setMaxPages] = useState([]);

  const updateMaxPage = maxPage => {
    if (!maxPages.find(o => o.source === maxPage.source)) {
      setMaxPages([...maxPages, maxPage]);
    }
  };

  const maxPageNo = maxPages.length > 0 ? Math.max(...maxPages.map(o => o.pages)) : 1;
  console.log("state :", state);
  console.log("maxPageNo :", maxPageNo);

  const params = useMemo(
    () =>
      state.sources
        .map(s => ({ source: s, pageNo: state.pageNo, query: state.searchText }))
        .filter(p => p.source === "tpb"),
    [state]
  );

  useEffect(() => {
    setMaxPages([]);
  }, [state.searchText, state.sources]);
  return (
    <StatusHandler status={sourcesStatus}>
      {() => (
        <SourceMapProvider value={sources}>
          <div className="app">
            <Search onSearch={payload => dispatch({ type: "SET_SEARCH_PARAMS", payload })} />
            {state.searchText && (
              <section>
                {params.map(params => (
                  <TorrentMetaList key={params.source} params={params} setTotalPageCount={updateMaxPage} />
                ))}
              </section>
            )}
            {params.length > 0 && (
              <Pagination
                current={state.pageNo}
                total={maxPageNo * 10}
                pageSize={10}
                onChange={pageNo => dispatch({ type: "SET_PAGE_NO", payload: pageNo })}
              />
            )}
          </div>
        </SourceMapProvider>
      )}
    </StatusHandler>
  );
}

function NotFound() {
  return "not found";
}

function App() {
  return (
    <AuthContext.Provider value={{ loggedIn: false }}>
      <Router id="router">
        <Login path="login" />
        <ProtectedComponent component={Dashboard} path="/" />
        <NotFound default />
      </Router>
    </AuthContext.Provider>
  );
}

export default App;
