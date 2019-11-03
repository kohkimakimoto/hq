import { createStore, applyMiddleware, compose, Store } from 'redux';
import { createLogger } from 'redux-logger';
import { StoreStateProperties, StoreState } from './State';

/**
 * Commit action.
 *
 * @type {string}
 */
const COMMIT = 'StoreState/COMMIT';

/**
 * Commit action creator.
 * @param {StoreState} state
 * @returns {any}
 */
function commit(state: StoreState): any {
  return {
    type: COMMIT,
    payload: state
  };
}

/**
 * reducer
 *
 * @param {StoreState} state
 * @param action
 * @returns {StoreState}
 */
export function reducer(state = new StoreState(), action: any): StoreState {
  switch (action.type) {
    case COMMIT:
      return action.payload;
    default:
      return state;
  }
}

/**
 * Dispatcher
 */
export class Dispatcher {
  private store: Store<StoreState>;

  public constructor(store: Store<StoreState>) {
    this.store = store;
  }

  public commit(props: StoreStateProperties): StoreState {
    const newState = this.store.getState().modify(props);
    this.store.dispatch(commit(newState));
    return newState;
  }
}
