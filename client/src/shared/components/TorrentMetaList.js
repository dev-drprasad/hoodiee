import React, { useMemo, useState, useEffect, useContext } from "react";
import { useFetch } from "@shared/hooks";
import { API_BASE_URL } from "@shared/consts";
import StatusHandler from "./StatusHandler";
import Spinner from "./Spinner";
import { ReactComponent as MagnetIcon } from "@shared/icons/magnet.svg";
import { ReactComponent as ArrowDownIcon } from "@shared/icons/arrow-down.svg";
import SourceMapContext from "@shared/contexts/SourceMap";
import AnchorBlank from "./AnchorBlank";

import "./TorrentMetaList.scss";

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
        <AnchorBlank href={magnetURI} title="magnet">
          <MagnetIcon />
        </AnchorBlank>
      ) : (
        <button
          className="fetch-magnet-button"
          onClick={() => setProps({ site: info.source, detailURL: info.URL })}
          title="click to fetch magnet"
        >
          <MagnetIcon />
          <ArrowDownIcon className="arrow-down-icon" />
        </button>
      )}
    </span>
  );
}

function useTorrentSearch(params) {
  const url = useMemo(() => {
    if (!params.query || !params.source) return;

    const url = new URL("/api/v1/search", API_BASE_URL);
    url.searchParams.append("query", params.query);
    url.searchParams.append("site", params.source);
    if (params.pageNo > 1) url.searchParams.append("pageNo", params.pageNo);
    return url;
  }, [params]);
  console.log("url :", url);
  return useFetch(url);
}

function TorrentMetaList({ params, setTotalPageCount }) {
  const [searchResult, status] = useTorrentSearch(params);
  const sources = useContext(SourceMapContext);

  useEffect(() => {
    if (status.isSuccess) {
      setTotalPageCount({ source: params.source, pages: searchResult.pages });
    }
  }, [status, setTotalPageCount, searchResult, params.source]);

  return (
    <>
      <div className="meta-list-header">
        <h3>{sources.find(s => s.id === params.source).name}</h3>
        <span>
          {params.pageNo} of {(searchResult || {}).pages}
        </span>
      </div>
      <StatusHandler status={status}>
        {() => (
          <ul className="meta-list">
            {searchResult.items.map(item => (
              <li className="gradient-on-hover meta-list-item">
                <h4 className="name">
                  <AnchorBlank href={item.URL}>{item.name}</AnchorBlank>
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
