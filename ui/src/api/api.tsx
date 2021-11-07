import { AxiosInstance } from 'axios';
import { makeUseAxios, UseAxios, Options as AxiosHooksOptions } from 'axios-hooks';
import { plainToInstance } from 'class-transformer';
import React, { createContext, useContext } from 'react';

import { checkValidationError } from '@/api/error';
import { useErrorHandler } from '@/lib/useErrorHandler';
import { Dashboard } from '@/models/Dashboard';
import { Job } from '@/models/Job';
import { Status } from '@/models/Status';

type ApiInfra = {
  axios: AxiosInstance;
  useAxios: UseAxios;
};

const ApiInfraContext = createContext<ApiInfra | undefined>(undefined);

type ApiProviderProps = {
  readonly axios: AxiosInstance;
  readonly children: React.ReactNode;
};

export const ApiProvider = (props: ApiProviderProps) => {
  return (
    <ApiInfraContext.Provider
      value={{
        axios: props.axios,
        useAxios: makeUseAxios({
          axios: props.axios,
        }),
      }}
    >
      {props.children}
    </ApiInfraContext.Provider>
  );
};

export const useApiInfra = () => {
  const context = useContext(ApiInfraContext);
  if (context === undefined) {
    throw new Error('useApiInfra must be used within a ApiProvider');
  }
  return context;
};

export type Options = AxiosHooksOptions;

export type UseGetDashboardParams = {
  name?: string;
  term?: string | null;
  begin?: string;
  reverse?: boolean;
  limit?: number;
  status?: Status;
};

export const useGetDashboard = (params: UseGetDashboardParams, options?: Options) => {
  const { useAxios } = useApiInfra();
  const errorHandler = useErrorHandler();

  const [values, execute] = useAxios({ method: 'GET', url: '/dashboard', params: params }, options);
  const validationError = checkValidationError(values.error);
  if (!validationError.hasError && values.error) {
    errorHandler(values.error);
  }

  const dashboard = values.data ? plainToInstance<Dashboard, any>(Dashboard, values.data) : undefined;

  return [{ dashboard, validationError, ...values }, execute] as const;
};

export const useGetJob = (id: string, options?: Options) => {
  const { useAxios } = useApiInfra();
  const errorHandler = useErrorHandler();

  const [values, execute] = useAxios({ method: 'GET', url: `/job/${id}` }, options);
  const validationError = checkValidationError(values.error);
  if (!validationError.hasError && values.error) {
    errorHandler(values.error);
  }

  const job = values.data ? plainToInstance<Job, any>(Job, values.data) : undefined;

  return [{ job, validationError, ...values }, execute] as const;
};

export const useDeleteJob = () => {
  const { useAxios } = useApiInfra();
  const errorHandler = useErrorHandler();

  const [values, execute] = useAxios({ method: 'DELETE' }, { manual: true });
  const validationError = checkValidationError(values.error);
  if (!validationError.hasError && values.error) {
    errorHandler(values.error);
  }

  return [{ validationError, ...values }, (id: string) => execute({ url: `/job/${id}` })] as const;
};

export const useRestartJob = () => {
  const { useAxios } = useApiInfra();
  const errorHandler = useErrorHandler();

  const [values, execute] = useAxios({ method: 'POST' }, { manual: true });
  const validationError = checkValidationError(values.error);
  if (!validationError.hasError && values.error) {
    errorHandler(values.error);
  }

  return [
    { validationError, ...values },
    (id: string, copy: boolean) => execute({ url: `/job/${id}/restart`, data: { copy: copy } }),
  ] as const;
};

export const useStopJob = () => {
  const { useAxios } = useApiInfra();
  const errorHandler = useErrorHandler();

  const [values, execute] = useAxios({ method: 'POST' }, { manual: true });
  const validationError = checkValidationError(values.error);
  if (!validationError.hasError && values.error) {
    errorHandler(values.error);
  }

  return [{ validationError, ...values }, (id: string) => execute({ url: `/job/${id}/stop` })] as const;
};
