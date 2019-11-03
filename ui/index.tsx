import React from 'react';
import ReactDOM from 'react-dom';
import { App } from './src/App';
import { configureStore } from "./src/store/configureStore";
import {Dispatcher} from "./src/store/Dispatcher";

// import style
import './style.scss';

declare global {
  const __DEV__: boolean;
}

if (__DEV__) {
  console.log('This is the development mode!');
}

const store = configureStore();
const dispatcher = new Dispatcher(store);

// initialize store state.
declare var appConfig: any;
dispatcher.commit(appConfig);

ReactDOM.render(<App store={store} />, document.getElementById('app'));
