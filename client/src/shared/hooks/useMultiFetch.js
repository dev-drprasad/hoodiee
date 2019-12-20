import { useState, useEffect, useReducer, useMemo } from "react";
import { NS, newline, tab } from "@shared/utils";

const defaultOptions = { headers: { Accept: "application/json", "Content-Type": "application/json" } };

class ResponseError extends Error {
  constructor(message, statusCode = 0) {
    super(message);
    this.name = "ResponseError";
    this.statusCode = statusCode;
  }
}

function jsonparse(res) {
  return (
    res
      .json()
      // grafana API responses dont have `json.code`
      .then(json => (json.code === undefined ? { code: res.status, ...json } : json))
      .catch(() => {
        throw new ResponseError("Invalid JSON response from API");
      })
  );
}

function reducer(state, action) {
  switch (action.type) {
    case "SET_RESPONSE":
      const id = action.payload.id;
      const response = action.payload.response;
      return [...state.slice(0, id), response, ...state.slice(id + 1)];

    default:
      return state;
  }
}

/**
 *
 * @param {String} url
 * @param {RequestInit} options
 * @returns {[any, NS]}
 */
export default function useMultiFetch(urls, options = defaultOptions) {
  // const initValue = useMemo(() => Array(urls.length).fill([undefined, new NS("INIT")]), [urls.length]);
  // console.log("initValue :", initValue);
  const [response, dispatch] = useReducer(reducer, undefined, () =>
    Array(urls.length).fill([undefined, new NS("INIT")])
  );
  console.log("response :", response);
  useEffect(() => {
    if (urls.length) {
      for (let i = 0; i < urls.length; i++) {
        const url = urls[i];

        dispatch({ type: "SET_RESPONSE", payload: { id: i, response: [undefined, new NS("LOADING")] } });
        fetch(url, options)
          .then(jsonparse)
          .then(json => {
            if (json.code >= 300) throw new ResponseError(json.error || "Unknown Error", json.code);

            dispatch({ type: "SET_RESPONSE", payload: { id: i, response: [json.data, new NS("SUCCESS")] } });
          })
          .catch(err => {
            console.error(
              `${newline}API Error:${newline}${tab}URL: ${url}${newline}${tab}MSG: ${err.message}${newline}${tab}CODE: ${err.statusCode}`
            );
            dispatch({
              type: "SET_RESPONSE",
              payload: { id: i, response: [undefined, new NS("ERROR", err.message, err.statusCode)] },
            });
          });
      }
    }
  }, [urls, options]);
  return response;
}
