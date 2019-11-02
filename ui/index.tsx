import React from 'react';
import ReactDOM from 'react-dom';
import { App } from './src/App';

// import style
import './style.scss';

declare var process: any;
if (process.env.NODE_ENV === 'development') {
  console.log('This is the development mode!');
}

// defined in html
declare var appConfig: any;

ReactDOM.render(<App initAppConfig={appConfig} />, document.getElementById('app'));
