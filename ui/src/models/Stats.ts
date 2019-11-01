export interface StatsProperties {
  readonly queueSize?: number;
  readonly queueUsage?: number;
  readonly runningWorkers?: number;
}

export class Stats implements StatsProperties {
  readonly queueSize: number = 0;
  readonly queueUsage: number = 0;
  readonly runningWorkers: number = 0;

  public constructor(props?: StatsProperties) {
    props && Object.assign(this, props);
  }

  public modify(props: StatsProperties): Stats {
    return new Stats(Object.assign({}, this, props));
  }
}
