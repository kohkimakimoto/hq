import '@/initGlobal';
import React from 'react';
import { render } from 'react-dom';

import { createAxiosInstance } from '@/api/createAxiosInstance';
import { AppProvider, App } from '@/App';
import { axiosLogger } from '@/lib/axiosLogger';
import { initConfig } from '@/lib/config';

declare global {
  const __DEV__: boolean;
  // It assumes to be embedded in a html that loads the code.
  const appConfig: {
    readonly basename: string;
    readonly version: string;
    readonly commitHash: string;
  };
}

const config = initConfig(appConfig);
const axios = createAxiosInstance(config.basename);

if (__DEV__) {
  console.log('This is the development mode!');
  axiosLogger(axios);
}

render(
  <AppProvider axios={axios} basename={config.basename}>
    <App />
  </AppProvider>,
  document.getElementById('app'),
);
