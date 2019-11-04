import { Stats } from '../models/Stats';

export type StoreStateProperties = {
  readonly basename?: string;
  readonly version?: string;
  readonly commitHash?: string;
  readonly error?: string;
  readonly stats?: Stats;
};

export class StoreState implements StoreStateProperties {
  readonly basename: string = '';
  readonly version: string = '';
  readonly commitHash: string = '';
  readonly error: string = '';
  readonly stats: Stats = new Stats();

  public constructor(props?: StoreStateProperties) {
    props && Object.assign(this, props);
  }

  public modify(props: StoreStateProperties): StoreState {
    return new StoreState(Object.assign({}, this, props));
  }
}
