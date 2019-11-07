import { ErrorValidationFailed } from '../api/Error';
import { Dispatcher } from '../store/Dispatcher';

export type HandleError = (err: any) => void;

export const createHandleError = (dispatcher: Dispatcher): HandleError => {
  return (err: any) => {
    if (err instanceof ErrorValidationFailed) {
      dispatcher.commit({
        error: err.message
      });
    } else if (typeof err === 'string') {
      dispatcher.commit({
        error: err
      });
    } else if (err.hasOwnProperty('message')) {
      dispatcher.commit({
        error: err.message
      });
    }

    console.log(err);
  };
};
