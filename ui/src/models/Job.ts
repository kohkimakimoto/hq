export type JobProperties = {
  readonly id?: string;
  readonly name?: string;
  readonly comment?: string;
  readonly url?: string;
  readonly payload?: any;
  readonly headers?: any;
  readonly timeout?: number;
  readonly createdAt?: string;
  readonly startedAt?: string | null;
  readonly finishedAt?: string | null;
  readonly failure?: boolean;
  readonly success?: boolean;
  readonly canceled?: boolean;
  readonly statusCode?: number | null;
  readonly err?: string;
  readonly output?: string;
  readonly waiting?: boolean;
  readonly running?: boolean;
  readonly status?: string;
};

export class Job implements JobProperties {
  readonly id: string = "";
  readonly name: string = "";
  readonly comment: string = "";
  readonly url: string = "";
  readonly payload: any = {};
  readonly headers: any = {};
  readonly timeout: number = 0;
  readonly createdAt: string = "";
  readonly startedAt: string | null = null;
  readonly finishedAt: string | null = null;
  readonly failure: boolean = false;
  readonly success: boolean = false;
  readonly canceled: boolean = false;
  readonly statusCode: number | null = null;
  readonly err: string = "";
  readonly output: string = "";
  readonly waiting: boolean = false;
  readonly running: boolean = false;
  readonly status: string = "";

  public constructor(props?: JobProperties) {
    props && Object.assign(this, props);
  }

  public modify(props: JobProperties): Job {
    return new Job(Object.assign({}, this, props));
  }
}
