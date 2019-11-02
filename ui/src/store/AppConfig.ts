import React, { Provider, useContext } from 'react';

export type AppConfigProperties = {
  basename?: string;
  version?: string;
  commitHash?: string;
};

export class AppConfig implements AppConfigProperties {
  basename: string = '';
  version: string = '';
  commitHash: string = '';

  public constructor(props?: AppConfigProperties) {
    props && Object.assign(this, props);
  }

  public modify(props: AppConfigProperties): AppConfig {
    return new AppConfig(Object.assign({}, this, props));
  }
}

const Context = React.createContext<AppConfig>(new AppConfig());

export const AppConfigProvider = Context.Provider;

export const useAppConfig = (): AppConfig => {
  return useContext<AppConfig>(Context);
};
