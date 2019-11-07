import moment from 'moment';
import humanizeDuration from 'humanize-duration';
import { SemanticCOLORS } from 'semantic-ui-react';

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
  readonly status?: 'failure' | 'success' | 'running' | 'waiting' | 'canceled' | 'canceling' | 'unfinished' | 'unknown';
};

export class Job implements JobProperties {
  readonly id: string = '';
  readonly name: string = '';
  readonly comment: string = '';
  readonly url: string = '';
  readonly payload: any = {};
  readonly headers: any = {};
  readonly timeout: number = 0;
  readonly createdAt: string = '';
  readonly startedAt: string | null = null;
  readonly finishedAt: string | null = null;
  readonly failure: boolean = false;
  readonly success: boolean = false;
  readonly canceled: boolean = false;
  readonly statusCode: number | null = null;
  readonly err: string = '';
  readonly output: string = '';
  readonly waiting: boolean = false;
  readonly running: boolean = false;
  readonly status: 'failure' | 'success' | 'running' | 'waiting' | 'canceled' | 'canceling' | 'unfinished' | 'unknown' =
    'unknown';

  public constructor(props?: JobProperties) {
    props && Object.assign(this, props);
  }

  public modify(props: JobProperties): Job {
    return new Job(Object.assign({}, this, props));
  }

  get createdAtFromNow(): string {
    return moment(this.createdAt)
      .locale('en')
      .fromNow();
  }

  get startedAtFromNow(): string {
    if (!this.startedAt) {
      return '';
    }

    return moment(this.startedAt)
      .locale('en')
      .fromNow();
  }

  get duration(): string {
    if (!this.finishedAt || !this.startedAt) {
      return '';
    }

    const startedAt = moment(this.startedAt);
    const finishedAt = moment(this.finishedAt);

    const diff = finishedAt.diff(startedAt);
    return humanizeDuration(diff);
  }

  get statusColor(): SemanticCOLORS {
    switch (this.status) {
      case 'failure':
        return 'red';
      case 'success':
        return 'green';
      case 'running':
        return 'blue';
      case 'waiting':
        return 'black';
      case 'canceled':
        return 'grey';
      case 'canceling':
        return 'grey';
      case 'unfinished':
        return 'grey';
      case 'unknown':
        return 'yellow';
      default:
        return 'grey';
    }
  }
}
