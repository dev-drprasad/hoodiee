import React, { useMemo, useState, useEffect } from "react";
import { useFetch } from "@shared/hooks";
import { API_BASE_URL } from "@shared/consts";
import StatusHandler from "./StatusHandler";
import Spinner from "./Spinner";

import "./TorrentMetaList.scss";
import useMultiFetch from "@shared/hooks/useMultiFetch";

function useDetail({ site, detailURL }) {
  const url = useMemo(() => {
    if (!site || !detailURL) return;
    const url = new URL("/api/v1/detail", API_BASE_URL);
    url.searchParams.append("url", detailURL);
    url.searchParams.append("site", site);
    return url.toString();
  }, [site, detailURL]);
  console.log("url :", url);
  return useFetch(url);
}

function MagnetURIFetcher({ info }) {
  const [props, setProps] = useState({});
  const [detailInfo = {}, infoStatus] = useDetail(props);

  const magnetURI = info.magnetURI || detailInfo.magnetURI;

  return (
    <span className="magnet-uri">
      {infoStatus.code === "LOADING" ? (
        <Spinner />
      ) : magnetURI ? (
        <a href={magnetURI}>magnet</a>
      ) : (
        <button onClick={() => setProps({ site: info.source, detailURL: info.URL })}>fetch magnet</button>
      )}
    </span>
  );
}

function useMultiTorrentSearch(searchText, pageNo, sources) {
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

function useTorrentSearch(params) {
  const url = useMemo(() => {
    if (!params.query || !params.source) return;

    const url = new URL("/api/v1/search", API_BASE_URL);
    url.searchParams.append("query", params.query);
    url.searchParams.append("site", params.source);
    url.searchParams.append("pageNo", params.pageNo);
    return url;
  }, [params]);
  console.log("url :", url);
  return useFetch(url);
}

function TorrentMetaList({ params, setPageNo }) {
  const [searchResult, status] = useTorrentSearch(params);

  useEffect(() => {
    if (status.isSuccess) {
      setPageNo({ source: params.source, pages: searchResult.pages });
    }
  }, [status, setPageNo, searchResult, params.source]);

  return (
    <>
      <h4>
        {params.source} (pages: {(searchResult || {}).pages})
      </h4>
      <StatusHandler status={status}>
        {() => (
          <ul className="meta-list">
            {searchResult.items.map(item => (
              <li className="meta-list-item">
                <h4 className="name">
                  <a href={item.URL}>{item.name}</a>
                </h4>
                <span className="seeders">seeders: {item.seeders}</span>
                <span className="leechers">leechers: {item.leechers}</span>
                <MagnetURIFetcher info={item} />
              </li>
            ))}
          </ul>
        )}
      </StatusHandler>
    </>
  );
}

export default TorrentMetaList;
