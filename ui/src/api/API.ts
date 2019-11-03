import { Client, HttpClientProvider } from './Client';
import { Stats } from '../models/Stats';

export class API {
  private client: Client;

  public constructor(provider?: HttpClientProvider) {
    this.client = new Client();
    if (provider) {
      this.registerHttpClientProvider(provider);
    }
  }

  public registerHttpClientProvider(provider: HttpClientProvider) {
    this.client.provider = provider;
  }

  public async getStats(): Promise<Stats> {
    const resp = await this.client.get('/stats');
    return new Stats(resp);
  }
}
