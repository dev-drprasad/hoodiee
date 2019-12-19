import { useState, useEffect } from "react";
import { NS, newline, tab } from "@shared/utils";

const defaultOptions = { headers: { Accept: "application/json", "Content-Type": "application/json" } };

class ResponseError extends Error {
  constructor(message, statusCode = 0) {
    super(message);
    this.name = "ResponseError";
    this.statusCode = statusCode;
  }
}

/**
 *
 * @param {String} url
 * @param {RequestInit} options
 * @returns {[any, NS]}
 */
export default function useFetch(url, options = defaultOptions, defaultValue) {
  const [response, setResponse] = useState([defaultValue, new NS("INIT")]);

  useEffect(() => {
    if (url) {
      setResponse([defaultValue, new NS("LOADING")]);
      fetch(url, options)
        .then(res => {
          return (
            res
              .json()
              // grafana API responses dont have `json.code`
              .then(json => (json.code === undefined ? { code: res.status, ...json } : json))
              .catch(() => {
                throw new ResponseError("Invalid JSON response from API");
              })
          );
        })
        .then(json => {
          if (json.code >= 300) throw new ResponseError(json.error || "Unknown error", json.code);
          setResponse([json.data, new NS("SUCCESS", undefined)]);
        })
        .catch(err => {
          console.error(
            `${newline}API Error:${newline}${tab}URL: ${url}${newline}${tab}MSG: ${err.message}${newline}${tab}CODE: ${err.statusCode}`
          );
          setResponse([defaultValue, new NS("ERROR", err.message, err.statusCode)]);
        });
    }
  }, [url, options, defaultValue]);
  return response;
}
