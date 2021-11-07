import { plainToInstance } from 'class-transformer';

import { Job } from '../Job';

test('plainToInstance', () => {
  const job = plainToInstance<Job, any>(Job, {
    id: '1',
    name: 'test',
    payload: {
      key1: 'value1',
      key2: 'value2',
    },
    timeout: 100,
    createdAt: '2021-12-18T06:31:17.977Z',
    failure: true,
  });

  expect(job.id).toBe('1');
  expect(job.name).toBe('test');
  expect(job.payload.key1).toBe('value1');
  expect(job.payload.key2).toBe('value2');
  expect(job.timeout).toBe(100);
  expect(job.createdAt.format('YYYY-MM-DD')).toBe('2021-12-18');
  expect(job.failure).toBe(true);
});
