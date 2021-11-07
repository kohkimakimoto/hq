import { render as rtlRender, RenderOptions } from '@testing-library/react';
import axios from 'axios';
import React, { ReactElement } from 'react';

import { AppProvider, AppProviderProps } from '@/App';

type WrapperProps = {
  children: React.ReactNode;
};

// eslint-disable-next-line import/export
export const render = (ui: ReactElement, props?: Partial<AppProviderProps>, options?: Omit<RenderOptions, 'wrapper'>) => {
  const axiosInstance = props?.axios ?? axios.create();
  const basename = props?.basename ?? '';

  const wrapper = ({ children }: WrapperProps) => {
    return (
      <AppProvider axios={axiosInstance} basename={basename}>
        {children}
      </AppProvider>
    );
  };

  return rtlRender(ui, { wrapper, ...options });
};

// eslint-disable-next-line import/export
export * from '@testing-library/react';
