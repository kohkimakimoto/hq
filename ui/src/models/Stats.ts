export type StatsProperties = {
  readonly version?: string;
  readonly commitHash?: string;
  readonly serverId?: number;
  readonly queues?: number;
  readonly dispatchers?: number;
  readonly maxWorkers?: number;
  readonly shutdownTimeout?: number;
  readonly jobLifetime?: number;
  readonly JobLifetimeStr?: string;
  readonly jobListDefaultLimit?: number;
  readonly queueMax?: number;
  readonly queueUsage?: number;
  readonly numWaitingJobs?: number;
  readonly numRunningJobs?: number;
  readonly numWorkers?: number;
  readonly numJobs?: number;
};

export class Stats implements StatsProperties {
  readonly version: string = '';
  readonly commitHash: string = '';
  readonly serverId: number = 0;
  readonly queues: number = 0;
  readonly dispatchers: number = 0;
  readonly maxWorkers: number = 0;
  readonly shutdownTimeout: number = 0;
  readonly jobLifetime: number = 0;
  readonly JobLifetimeStr: string = '';
  readonly jobListDefaultLimit: number = 0;
  readonly queueMax: number = 0;
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
