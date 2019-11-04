export type StatsProperties = {
  readonly queues?: number;
  readonly dispatchers?: number;
  readonly maxWorkers?: number;
  readonly queueUsage?: number;
  readonly numWaitingJobs?: number;
  readonly numRunningJobs?: number;
  readonly numWorkers?: number;
  readonly numJobs?: number;
};

export class Stats implements StatsProperties {
  readonly queues: number = 0;
  readonly dispatchers: number = 0;
  readonly maxWorkers: number = 0;
  readonly queueUsage: number = 0;
  readonly numWaitingJobs: number = 0;
  readonly numRunningJobs: number = 0;
  readonly numWorkers: number = 0;
  readonly numJobs: number = 0;

  public constructor(props?: StatsProperties) {
    props && Object.assign(this, props);
  }

  public modify(props: StatsProperties): Stats {
    return new Stats(Object.assign({}, this, props));
  }
}
