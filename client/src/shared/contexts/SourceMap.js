import { createContext } from "react";

const SourceMapContext = createContext([]);
export default SourceMapContext;

export const SourceMapProvider = SourceMapContext.Provider;
export const SourceMapConsumer = SourceMapContext.Consumer;
