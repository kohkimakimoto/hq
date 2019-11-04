import {Job} from "./Job";

export type JobListProperties = {
  readonly jobs?: Job[];
  readonly hasNext?: boolean;
  readonly next?: string | null;
  readonly count?: number;
};

export class JobList implements JobListProperties {
  readonly jobs: Job[] = [];
  readonly hasNext: boolean = false;
  readonly next: string | null = null;
  readonly count: number = 0;

  public constructor(props?: JobListProperties) {
    props && Object.assign(this, props);
  }

  public modify(props: JobListProperties): JobList {
    return new JobList(Object.assign({}, this, props));
  }
}
