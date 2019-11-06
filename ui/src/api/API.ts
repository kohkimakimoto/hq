import { Client } from './Client';
import { Stats } from '../models/Stats';
import { AxiosInstance } from 'axios';
import { JobList } from '../models/JobList';
import { Job } from '../models/Job';

export class API {
  private client: Client;

  public constructor(provider?: () => AxiosInstance) {
    this.client = new Client();
    if (provider) {
      this.registerHttpClientProvider(provider);
    }
  }

  public registerHttpClientProvider(provider: () => AxiosInstance) {
    this.client.httpClientProvider = provider;
  }

  public async stats(): Promise<Stats> {
    const resp = await this.client.get('/stats');
    return new Stats(resp);
  }

  public async listJobs(data: {
    readonly name?: string;
    readonly term?: string;
    readonly begin?: string;
    readonly reverse?: boolean;
    readonly limit?: number;
    readonly status?: string;
  }): Promise<JobList> {
    const resp = await this.client.get('/job', data);

    return new JobList({
      jobs: resp.jobs.map(value => {
        return new Job(value);
      }),
      hasNext: resp.hasNext,
      next: resp.next,
      count: resp.count
    });
  }

  public async getJob(id: string): Promise<Job> {
    const resp = await this.client.get('/job/' + id);
    return new Job(resp);
  }
}
