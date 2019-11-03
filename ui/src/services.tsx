import React, { useContext } from 'react';
import { Store } from 'redux';
import { StoreState } from './store/State';
import { ErrorHandler } from './ErrorHandler';
import { Dispatcher } from './store/Dispatcher';
import { configureStore } from './store/configureStore';

const store = configureStore();

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

  get error() {
    return new ErrorHandler(this.dispatcher);
  }
}

// Init default resolver.
ServiceResolver.defaultResolver = new ServiceResolver(configureStore());

const context = React.createContext<ServiceResolver>(ServiceResolver.defaultResolver);

export const ServiceProvider = context.Provider;

export const useServices = () => {
  return useContext<ServiceResolver>(context);
};
