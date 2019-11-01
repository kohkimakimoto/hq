/**
 * ErrorValidationFailed
 */
export class ErrorValidationFailed {
  /**
   * message
   */
  public message: string;

  /**
   * constructor
   *
   * @param {string} message
   */
  public constructor(message: string) {
    this.message = message;
  }
}

export function errorParser(err): Promise<any> {
  if (err.response) {
    // received response
    const { data, status } = err.response;

    if (data instanceof Object) {
      // received a JSON body.
      if (status === 422 && data.errors !== 'undefined') {
        // validation failed
        return Promise.reject(new ErrorValidationFailed(data.message));
      } else if (typeof data.message !== 'undefined') {
        // has custom error message.
        return Promise.reject(new Error(data.message));
      } else {
        // no custom error message.
        return Promise.reject(new Error('Error with status ' + status));
      }
    }

    return Promise.reject(new Error('Error with status ' + status));
  }

  return Promise.reject(err);
}
