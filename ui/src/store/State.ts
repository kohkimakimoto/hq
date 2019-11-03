export type StoreStateProperties = {
  readonly basename?: string;
  readonly version?: string;
  readonly commitHash?: string;
};

export class StoreState implements StoreStateProperties {
  readonly basename: string = '';
  readonly version: string = '';
  readonly commitHash: string = '';

  public constructor(props?: StoreStateProperties) {
    props && Object.assign(this, props);
  }

  public modify(props: StoreStateProperties): StoreState {
    return new StoreState(Object.assign({}, this, props));
  }
}
