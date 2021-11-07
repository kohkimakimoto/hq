import { ChakraProvider } from '@chakra-ui/react';
import { AxiosInstance } from 'axios';
import React, { FunctionComponent } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { HelmetProvider } from 'react-helmet-async';
import { RouteObject, BrowserRouter, useRoutes } from 'react-router-dom';
import { QueryParamProvider } from 'use-query-params';

import { ApiProvider } from '@/api/api';
import { Layout } from '@/components/Layout';
import { RouteAdapter } from '@/lib/RouteAdapter';
import { IndexPage } from '@/pages/IndexPage';
import { JobDetailPage } from '@/pages/job/JobDetailPage';
import { NotFoundPage } from '@/pages/NotFoundPage';
import theme from '@/theme';

const ErrorFallback = () => {
  return <div>Something went wrong.</div>;
};

export type AppProviderProps = {
  children: React.ReactNode;
  axios: AxiosInstance;
  basename: string;
};

export const AppProvider = (props: AppProviderProps) => {
  const { children, axios, basename } = props;
  return (
    <ChakraProvider theme={theme}>
      <ErrorBoundary FallbackComponent={ErrorFallback}>
        <ApiProvider axios={axios}>
          <HelmetProvider>
            <BrowserRouter basename={basename}>
              <QueryParamProvider ReactRouterRoute={RouteAdapter as unknown as FunctionComponent}>{children}</QueryParamProvider>
            </BrowserRouter>
          </HelmetProvider>
        </ApiProvider>
      </ErrorBoundary>
    </ChakraProvider>
  );
};

const routes: RouteObject[] = [
  {
    path: '/',
    element: <Layout />,
    children: [
      { index: true, element: <IndexPage /> },
      { path: '/jobs/:id', element: <JobDetailPage /> },
      { path: '*', element: <NotFoundPage /> },
    ],
  },
];

export const App = () => {
  return <>{useRoutes(routes)}</>;
};
