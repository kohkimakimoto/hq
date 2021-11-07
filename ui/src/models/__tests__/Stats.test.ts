import { plainToInstance } from 'class-transformer';

import { Stats } from '../Stats';

test('plainToInstance', () => {
  const stats = plainToInstance<Stats, any>(Stats, {
    queues: 1111,
    dispatchers: 2222,
  });

  expect(stats.queues).toBe(1111);
  expect(stats.dispatchers).toBe(2222);
});
