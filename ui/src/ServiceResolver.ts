import React, { useContext } from 'react';
import { Store } from 'redux';
import { StoreState } from './store/State';
import { Dispatcher } from './store/Dispatcher';
import { configureStore } from './store/configureStore';
import { API } from './api/API';
import { createHTTPClientProvider } from './api/createHTTPClientProvider';
import { createHandleError, HandleError } from './services/HandleError';
import { createUpdateStats, UpdateStats } from './services/UpdateStats';

/**
 * ServiceResolver
 */
export class ServiceResolver {
  public store: Store<StoreState>;

  public static defaultResolver: ServiceResolver;

  public constructor(store: Store<StoreState>) {
    this.store = store;

    ServiceResolver.defaultResolver = this;
  }

  get dispatcher(): Dispatcher {
    return new Dispatcher(this.store);
  }

  get handleError(): HandleError {
    return createHandleError(this.dispatcher);
  }

  get updateStats(): UpdateStats {
    return createUpdateStats(this.api, this.dispatcher);
  }

  get api(): API {
    return new API(createHTTPClientProvider(this.store.getState().basename));
  }
}
