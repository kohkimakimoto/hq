import React from 'react';
import ReactDOM from 'react-dom';
import { App } from './src/App';

declare var process: any;
if (process.env.NODE_ENV === 'development') {
  console.log('This is the development mode!');
}

ReactDOM.render(<App />, document.getElementById('app'));
