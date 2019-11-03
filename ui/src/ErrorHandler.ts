import { ErrorValidationFailed } from './api/Error';
import { Dispatcher } from './store/Dispatcher';

// ErrorHandler
export class ErrorHandler {
  private readonly dispatcher: Dispatcher;

  public constructor(dispatcher: Dispatcher) {
    this.dispatcher = dispatcher;
  }

  public handle(err: any, outputConsole = true): void {
    if (outputConsole) {
      console.log(err);
    }

    if (err instanceof ErrorValidationFailed) {
    } else if (typeof err === 'string') {
    } else {
    }
  }
}
