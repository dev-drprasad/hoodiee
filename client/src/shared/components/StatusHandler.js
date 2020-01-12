import React from "react";
import Spinner from "./Spinner";

/**
 * props of StatusHandler
 * @typedef {Object} StatusHandlerProps
 * @property {import('@shared/utils').NS} status -
 * @property {Boolean} ignoreError
 * @property {Boolean} hasData - Represents whenther data is empty or not
 * @property {String} height - Height of container
 * @property {React.ReactNode} children - Children of StatusHandler
 */

/**
 * Generic component to handle network requests
 * @param {StatusHandlerProps}
 */
export default function StatusHandler({
  status,
  styles = {},
  data = null,
  hasData: overriddenHasData,
  height,
  ignoreError = false,
  children = null,
}) {
  const hasData = overriddenHasData === undefined ? status.hasData : overriddenHasData;
  if (status.isSuccess && hasData) {
    return children(data);
  }

  if (status.isError && ignoreError) {
    return children(data);
  }

  let fallback;
  if (status.isError) {
    if (ignoreError) {
      return children(data);
    } else if (status.statusCode === 404) {
      fallback = "Resource not found.";
    } else {
      fallback = "Oops! Something went wrong. Please try reloading again.";
    }
  } else if (status.isLoading) {
    fallback = <Spinner />;
  } else if (status.isSuccess && !hasData) {
    fallback = "No data available.";
  }

  return (
    <div
      className="status-handler-fallback"
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        alignSelf: "center",
        ...(height && { height }),
      }}
    >
      {fallback}
    </div>
  );
}
