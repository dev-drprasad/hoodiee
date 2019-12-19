export const assert = (condition, message) => {
  if (!condition) throw Error(message || "Assertion failed");
};

export default class NS {
  /**
   * Create a point.
   * @param {'INIT'|'LOADING'|'SUCCESS'|'ERROR'} status
   * @param {string|null} message
   */
  constructor(status, message, statusCode = 0, hasData = true) {
    this.code = status;
    this.message = message;
    this.statusCode = statusCode;
    this.hasData = hasData;
  }

  get isLoading() {
    return this.code === "LOADING";
  }

  get isError() {
    return this.code === "ERROR";
  }

  get isSuccess() {
    return this.code === "SUCCESS";
  }

  toString() {
    return this.code;
  }
}
