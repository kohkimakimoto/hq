import { Transform, Type } from 'class-transformer';
import dayjs, { Dayjs } from 'dayjs';

import { Status } from './Status';

export class Job {
  public id = '';

  public name = '';

  public comment = '';

  public url = '';

  public payload: any = {};

  public headers: any = {};

  public timeout = 0;

  @Type(() => Date)
  @Transform(({ value }) => dayjs(value), { toClassOnly: true })
  public createdAt: Dayjs = dayjs();

  @Type(() => Date)
  @Transform(({ value }) => (value ? dayjs(value) : null), { toClassOnly: true })
  public startedAt: Dayjs | null = null;

  @Type(() => Date)
  @Transform(({ value }) => (value ? dayjs(value) : null), { toClassOnly: true })
  public finishedAt: Dayjs | null = null;

  public failure = false;

  public success = false;

  public canceled = false;

  public statusCode: number | null = null;

  public err = '';

  public output = '';

  public waiting = false;

  public running = false;

  public status: Status = 'unknown';

  get duration(): string {
    if (!this.finishedAt || !this.startedAt) {
      return '';
    }

    const diff = this.finishedAt.diff(this.startedAt, 'seconds');
    return dayjs.duration(diff, 'seconds').humanize();
  }
}
