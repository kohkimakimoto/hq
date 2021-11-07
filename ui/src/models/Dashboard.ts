import { Type } from 'class-transformer';

import { JobList } from './JobList';
import { Stats } from './Stats';

export class Dashboard {
  @Type(() => Stats)
  public stats: Stats = new Stats();

  @Type(() => JobList)
  public jobList: JobList = new JobList();
}
