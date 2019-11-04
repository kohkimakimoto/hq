import axios, { AxiosInstance } from 'axios';
import { addLogger } from './addLogger';

export function createHTTPClientProvider(basename: string): () => AxiosInstance {
  return () => {
    const headers = {};

    const client = axios.create({
      baseURL: basename + '/internal',
      timeout: 10000,
      headers: headers
    });

    if (__DEV__) {
      addLogger(client);
    }

    return client;
  };
}
