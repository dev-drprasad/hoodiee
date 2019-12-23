import { NS } from "@shared/utils";

/**
 * Derive final status from given statuses
 * @param {...import('@shared/utils').NS} statuses
 */
const mergeStatuses = (...statuses) => {
  const hasData = statuses.some(s => s.hasData);

  const code = statuses.every(s => s.isSuccess)
    ? "SUCCESS" // [SUCCESS, SUCCESS, SUCCESS]
    : statuses.find(s => s.isError)
    ? "ERROR" // [SUCCESS, ERROR] or [LOADING, ERROR]
    : statuses.find(s => s.isLoading)
    ? "LOADING" // [SUCCESS, LOADING] or [LOADING, INIT] or [LOADING, LOADING]
    : statuses.find(s => s.isSuccess)
    ? "LOADING" // [SUCCESS, INIT]
    : "INIT"; // [INIT, INIT]

  return new NS(code, undefined, 0, hasData);
};

export default mergeStatuses;
