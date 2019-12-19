import React, { useMemo, useState } from "react";
import "./TorrentMetaList.scss";
import { useFetch } from "@shared/hooks";
import { API_BASE_URL } from "@shared/consts";
import StatusHandler from "./StatusHandler";

function useDetail({ site, detailURL }) {
  const url = useMemo(() => {
    if (!site || !detailURL) return;
    const url = new URL("/api/v1/detail", API_BASE_URL);
    url.searchParams.append("url", detailURL);
    url.searchParams.append("site", site);
    return url.toString();
  }, [site, detailURL]);
  return useFetch(url);
}

function MagnetURIFetcher({ site, detailURL }) {
  const [props, setProps] = useState({});
  const [info = {}, infoStatus] = useDetail(props);
  if (infoStatus.code === "INIT") return <button onClick={() => setProps({ site, detailURL })}>fetch magnet</button>;
  return (
    <StatusHandler status={infoStatus}>
      {() => {
        return info.magnetURI ? (
          <a className="magnet-uri" href={info.magnetURI}>
            magnet
          </a>
        ) : (
          "no magnet found"
        );
      }}
    </StatusHandler>
  );
}

function TorrentMetaList({ items, site }) {
  return (
    <ul className="meta-list">
      {items.map(item => (
        <li className="meta-list-item">
          <h4 className="name">
            <a href={item.URL}>{item.name}</a>
          </h4>
          <span className="seeders">seeders: {item.seeders}</span>
          <span className="leechers">leechers: {item.leechers}</span>
          {item.magnetURI ? (
            <a className="magnet-uri" href={item.magnetURI}>
              magnet
            </a>
          ) : (
            <MagnetURIFetcher site={site} detailURL={item.URL} />
          )}
        </li>
      ))}
    </ul>
  );
}

export default TorrentMetaList;
