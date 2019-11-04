import React from 'react';
import { configureStore } from './store/configureStore';
import { ServiceResolver } from './ServiceResolver';

// Init default resolver.
const store = configureStore();
ServiceResolver.defaultResolver = new ServiceResolver(configureStore());

export const ServiceContext = React.createContext<ServiceResolver>(ServiceResolver.defaultResolver);
