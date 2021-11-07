import { useColorMode } from '@chakra-ui/react';

export type Status = 'failure' | 'success' | 'running' | 'waiting' | 'canceled' | 'canceling' | 'unfinished' | 'unknown';

export type StatusColors = {
  [key in Status]: string;
};

export const useStatusColors = (): StatusColors => {
  const { colorMode } = useColorMode();
  if (colorMode === 'light') {
    return {
      success: 'green.500',
      failure: 'red.500',
      running: 'blue.500',
      waiting: 'gray.500',
      canceled: 'gray.500',
      canceling: 'gray.500',
      unfinished: 'gray.500',
      unknown: 'gray.500',
    } as const;
  } else {
    return {
      success: 'green.500',
      failure: 'red.500',
      running: 'blue.500',
      waiting: 'gray.500',
      canceled: 'gray.500',
      canceling: 'gray.500',
      unfinished: 'gray.500',
      unknown: 'gray.500',
    } as const;
  }
};
