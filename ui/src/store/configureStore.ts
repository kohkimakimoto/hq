import { applyMiddleware, compose, createStore, Store } from 'redux';
import { reducer } from './Dispatcher';
import { StoreState } from './State';
import { createLogger } from 'redux-logger';

export function configureStore(): Store<StoreState> {
  // middleware.
  let middleware = [];
  if (__DEV__) {
    middleware.push(createLogger() as never);
  }

  // store
  return createStore(reducer, compose(applyMiddleware(...middleware)));
}
