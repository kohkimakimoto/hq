export type Config = {
  readonly basename: string;
  readonly version: string;
  readonly commitHash: string;
};

export let config: Config = {
  basename: '',
  version: '',
  commitHash: '',
};

export const initConfig = (c: Config): Config => {
  config = c;
  return config;
};
