import { Client, HttpClientProvider } from './Client';

export class API {
  private client: Client;

  public constructor() {
    this.client = new Client();
  }

  public registerHttpClientProvider(provider: HttpClientProvider) {
    this.client.provider = provider;
  }
}
