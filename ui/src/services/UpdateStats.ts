import { API } from '../api/API';
import { Dispatcher } from '../store/Dispatcher';
import { HandleError } from './HandleError';

export type UpdateStats = () => Promise<void>;

export const createUpdateStats = (api: API, dispatcher: Dispatcher): UpdateStats => {
  return async () => {
    const stats = await api.stats();
    dispatcher.commit({
      stats: stats
    });
  };
};
