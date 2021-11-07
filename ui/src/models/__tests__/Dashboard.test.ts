import { plainToInstance } from 'class-transformer';

import { Dashboard } from '../Dashboard';

test('plainToInstance', () => {
  const dashboard = plainToInstance<Dashboard, any>(Dashboard, {
    stats: {
      queues: 1111,
      dispatchers: 2222,
    },
    jobList: {
      jobs: [
        {
          id: '1',
          name: 'test1',
          payload: {
            key1: 'value1',
            key2: 'value2',
          },
          timeout: 100,
          createdAt: '2021-12-18T06:31:17.977Z',
          failure: true,
        },
        {
          id: '2',
          name: 'test2',
          payload: {
            key1: 'value1',
            key2: 'value2',
          },
          timeout: 100,
          createdAt: '2021-12-18T06:31:17.977Z',
          failure: true,
        },
      ],
      count: 2,
    },
  });

  expect(dashboard.stats.queues).toBe(1111);
  expect(dashboard.stats.dispatchers).toBe(2222);

  expect(dashboard.jobList.count).toBe(2);
  expect(dashboard.jobList.jobs.length).toBe(2);

  expect(dashboard.jobList.jobs[0].id).toBe('1');
  expect(dashboard.jobList.jobs[0].name).toBe('test1');
  expect(dashboard.jobList.jobs[0].payload.key1).toBe('value1');
  expect(dashboard.jobList.jobs[0].payload.key2).toBe('value2');
  expect(dashboard.jobList.jobs[0].timeout).toBe(100);
  expect(dashboard.jobList.jobs[0].createdAt.format('YYYY-MM-DD')).toBe('2021-12-18');
  expect(dashboard.jobList.jobs[0].failure).toBe(true);

  expect(dashboard.jobList.jobs[1].id).toBe('2');
  expect(dashboard.jobList.jobs[1].name).toBe('test2');
  expect(dashboard.jobList.jobs[1].payload.key1).toBe('value1');
  expect(dashboard.jobList.jobs[1].payload.key2).toBe('value2');
  expect(dashboard.jobList.jobs[1].timeout).toBe(100);
  expect(dashboard.jobList.jobs[1].createdAt.format('YYYY-MM-DD')).toBe('2021-12-18');
  expect(dashboard.jobList.jobs[1].failure).toBe(true);
});
