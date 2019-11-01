import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';
import { errorParser } from './Error';

/**
 * HttpClientProvider
 */
export type HttpClientProvider = () => AxiosInstance;

/**
 * Client is an API client to make requests to the myRule API.
 * The myRule APIs only use POST request.
 */
export class Client {
  public provider: HttpClientProvider;

  public constructor() {
    this.provider = function() {
      return axios.create();
    };
  }

  public registerHttpClientProvider(provider: HttpClientProvider) {
    this.provider = provider;
  }

  public put(url: string, data?: any, config?: AxiosRequestConfig): Promise<any> {
    return this.provider()
      .put(url, data, config)
      .then(resp => {
        return resp.data;
      })
      .catch(err => {
        return errorParser(err);
      });
  }

  public post(url: string, data?: any, config?: AxiosRequestConfig): Promise<any> {
    return this.provider()
      .post(url, data, config)
      .then(resp => {
        return resp.data;
      })
      .catch(err => {
        return errorParser(err);
      });
  }

  public delete(url: any, params: any = null, config: AxiosRequestConfig = {}): Promise<any> {
    if (params) {
      config.params = params;
    }

    return this.provider()
      .delete(url, config)
      .then(resp => {
        return resp.data;
      })
      .catch(err => {
        return errorParser(err);
      });
  }

  public get(url: any, params: any = null, config: AxiosRequestConfig = {}): Promise<any> {
    if (params) {
      config.params = params;
    }

    return this.provider()
      .get(url, config)
      .then(resp => {
        return resp.data;
      })
      .catch(err => {
        return errorParser(err);
      });
  }
}
