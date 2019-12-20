import React, { useMemo, useState } from "react";
import "./TorrentMetaList.scss";
import { useFetch } from "@shared/hooks";
import { API_BASE_URL } from "@shared/consts";
import StatusHandler from "./StatusHandler";
import Spinner from "./Spinner";

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

function TorrentMetaList({ searchResult }) {
  console.log("sources :", searchResult);
  return searchResult.map(([source, [items = [], status]]) => (
    <>
      <h4>{source}</h4>
      <StatusHandler status={status}>
        {() => (
          <ul className="meta-list">
            {items.map(item => (
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
  ));
}

export default TorrentMetaList;
