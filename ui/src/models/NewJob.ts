export type NewJobProperties = {
  readonly name?: string;
  readonly comment?: string;
  readonly url?: string;
  readonly payloadString?: string;
  readonly headersString?: string;
  readonly timeoutString?: string;
};

export class NewJob implements NewJobProperties {
  readonly name: string = '';
  readonly comment: string = '';
  readonly url: string = '';
  readonly payloadString: string = '';
  readonly headersString: string = '';
  readonly timeoutString: string = '';

  public constructor(props?: NewJobProperties) {
    props && Object.assign(this, props);
  }

  public modify(props: NewJobProperties): NewJob {
    return new NewJob(Object.assign({}, this, props));
  }

  get payload(): any {
    if (this.payloadString == '') {
      return null;
    }

    try {
      return JSON.parse(this.payloadString);
    } catch (err) {
      throw new Error('payload has error: ' + err);
    }
  }

  get headers(): any {
    if (this.headersString == '') {
      return null;
    }

    try {
      return JSON.parse(this.headersString);
    } catch (err) {
      throw new Error('headers has error: ' + err);
    }
  }

  get timeout(): number {
    return Number(this.timeoutString);
  }
}
