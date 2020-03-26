import axios, { AxiosInstance, AxiosRequestConfig } from 'axios';
import { errorParser } from './Error';

/**
 * Client is an API client to make requests to the API.
 */
export class Client {
  public httpClientProvider: () => AxiosInstance;

  public constructor() {
    this.httpClientProvider = function() {
      return axios.create();
    };
  }

  public registerHttpClientProvider(provider: () => AxiosInstance) {
    this.httpClientProvider = provider;
  }

  public put(url: string, data?: any, config?: AxiosRequestConfig): Promise<any> {
    return this.httpClientProvider()
      .put(url, data, config)
      .then(resp => {
        return resp.data;
      })
      .catch(err => {
        return errorParser(err);
      });
  }

  public post(url: string, data?: any, config?: AxiosRequestConfig): Promise<any> {
    return this.httpClientProvider()
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

    return this.httpClientProvider()
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

    return this.httpClientProvider()
      .get(url, config)
      .then(resp => {
        return resp.data;
      })
      .catch(err => {
        return errorParser(err);
      });
  }
}
