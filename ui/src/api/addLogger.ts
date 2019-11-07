import axios, { AxiosInstance, AxiosPromise, AxiosRequestConfig } from 'axios';

/**
 * API logger
 * @param instance
 */
export function addLogger(instance: AxiosInstance) {
  instance.interceptors.request.use(function(config) {
    if ((console as any).group) {
      (console as any).group(
        '%c request %c' + config.method!.toUpperCase() + ' ' + config.baseURL + config.url,
        'color: gray; font-weight: lighter;',
        'font-weight: bold'
      );
      console.log(' headers', config.headers);
      console.log(' data   ', config.data);
      console.log(' config ', config);
      (console as any).groupEnd();
    } else {
      console.log(config.method!.toUpperCase() + ' ' + config.baseURL + config.url);
      console.log(' headers ', config.headers);
      console.log(' data  ', config.data);
      console.log(' config', config);
    }

    return config;
  });

  instance.interceptors.response.use(
    function(response) {
      const { data, status, headers, config } = response;
      if ((console as any).group) {
        (console as any).group(
          '%c response %c' + status + ' ' + config.method!.toUpperCase() + ' ' + config.url,
          'color: gray; font-weight: lighter;',
          'font-weight: bold'
        );
        console.log(' headers', headers);
        console.log(' data   ', data);
        console.log(' config ', config);
        (console as any).groupEnd();
      } else {
        console.log('Response');
        console.log('data  ', response.data);
      }
      return response;
    },
    function(error) {
      if (error.response) {
        const response = error.response;
        const { data, status, headers, config } = response;
        if ((console as any).group) {
          (console as any).group(
            '%c response %c' + status + ' ' + config.method!.toUpperCase() + ' ' + config.url,
            'color: red; font-weight: lighter;',
            'font-weight: bold'
          );
          console.log(' error  ', error);
          console.log(' headers', headers);
          console.log(' data   ', data);
          console.log(' config ', config);
          (console as any).groupEnd();
        } else {
          console.log(error.name + ' ' + error.message);
          console.log('status ', status);
          console.log('data   ', data);
        }
      } else {
        if ((console as any).group) {
          console.log('error ', error);
        } else {
          console.log('error ', error);
        }
      }

      throw error;
    }
  );
}
