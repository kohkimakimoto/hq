import axios from 'axios';

export const createAxiosInstance = (basename: string, timeout = 10000) => {
  return axios.create({
    baseURL: `${basename}/internal`,
    timeout,
  });
};
