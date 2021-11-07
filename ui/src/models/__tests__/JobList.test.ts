import { plainToInstance } from 'class-transformer';

import { JobList } from '../JobList';

test('plainToInstance', () => {
  const jobList = plainToInstance<JobList, any>(JobList, {
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
  });

  expect(jobList.count).toBe(2);
  expect(jobList.jobs.length).toBe(2);

  expect(jobList.jobs[0].id).toBe('1');
  expect(jobList.jobs[0].name).toBe('test1');
  expect(jobList.jobs[0].payload.key1).toBe('value1');
  expect(jobList.jobs[0].payload.key2).toBe('value2');
  expect(jobList.jobs[0].timeout).toBe(100);
  expect(jobList.jobs[0].createdAt.format('YYYY-MM-DD')).toBe('2021-12-18');
  expect(jobList.jobs[0].failure).toBe(true);

  expect(jobList.jobs[1].id).toBe('2');
  expect(jobList.jobs[1].name).toBe('test2');
  expect(jobList.jobs[1].payload.key1).toBe('value1');
  expect(jobList.jobs[1].payload.key2).toBe('value2');
  expect(jobList.jobs[1].timeout).toBe(100);
  expect(jobList.jobs[1].createdAt.format('YYYY-MM-DD')).toBe('2021-12-18');
  expect(jobList.jobs[1].failure).toBe(true);
});
