import { Type } from 'class-transformer';

import { Job } from './Job';

export class JobList {
  @Type(() => Job)
  public jobs: Job[] = [];

  public hasNext = false;

  public next: string | null = null;

  public count = 0;
}
