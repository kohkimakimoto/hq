export class Stats {
  public queues = 0;

  public dispatchers = 0;

  public maxWorkers = 0;

  public numWorkers = 0;

  public numJobsInQueue = 0;

  public numJobsWaiting = 0;

  public numJobsRunning = 0;

  public numStoredJobs = 0;

  public numJobsInLastMinute = 0;

  get queueUsageRate(): number {
    return Math.floor((this.numJobsInQueue / this.queues) * 100);
  }

  get maxConcurrentWorkers(): number {
    return (this.maxWorkers ? this.maxWorkers : 1) * this.dispatchers;
  }
}
