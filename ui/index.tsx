import React from 'react';
import ReactDOM from 'react-dom';
import { App } from './src/App';
import { configureStore } from './src/store/configureStore';
import { Dispatcher } from './src/store/Dispatcher';
import { ServiceResolver } from './src/ServiceResolver';

// import style
import './style.scss';

declare global {
  const __DEV__: boolean;
  // It assumes to be embedded in a html that loads the code.
  const appConfig: any;
}

if (__DEV__) {
  console.log('This is the development mode!');
}

const resolver = ServiceResolver.defaultResolver;

// init store state.
resolver.dispatcher.commit(appConfig);

ReactDOM.render(<App resolver={resolver} />, document.getElementById('app'));
