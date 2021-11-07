import { AxiosInstance } from 'axios';

/**
 * Axios logger
 *
 * @param {AxiosInstance} instance
 */
export function axiosLogger(instance: AxiosInstance) {
  instance.interceptors.request.use(function (config) {
    if ((console as any).group) {
      (console as any).group(
        '%caxios request %c' + config.method?.toUpperCase() + ' ' + config.baseURL + config.url,
        'color: blue; font-weight: lighter;',
        'font-weight: bold',
      );
      console.log(' headers', config.headers);
      console.log(' data   ', config.data);
      console.log(' config ', config);
      (console as any).groupEnd();
    } else {
      console.log(config.method?.toUpperCase() + ' ' + config.baseURL + config.url);
      console.log(' headers ', config.headers);
      console.log(' data  ', config.data);
      console.log(' config', config);
    }

    return config;
  });

  instance.interceptors.response.use(
    function (response) {
      const { data, status, headers, config } = response;
      if ((console as any).group) {
        (console as any).group(
          '%caxios response %c' + status + ' ' + config.method?.toUpperCase() + ' ' + config.baseURL + config.url,
          'color: blue; font-weight: lighter;',
          'font-weight: bold',
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
    function (error) {
      if (error.response) {
        const response = error.response;
        const { data, status, headers, config } = response;
        if ((console as any).group) {
          (console as any).group(
            '%caxios response %c' + status + ' ' + config.method?.toUpperCase() + ' ' + config.baseURL + config.url,
            'color: red; font-weight: lighter;',
            'font-weight: bold',
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
          (console as any).group('%caxios no response', 'color: red; font-weight: lighter;');
          console.log(' error  ', error);
          (console as any).groupEnd();
        } else {
          console.log('error  ', error);
        }
      }

      throw error;
    },
  );
}
