import { Flex, Text } from '@chakra-ui/react';
import React from 'react';

export type JobListContainerProps = {
  children: React.ReactNode;
  isNoJobs: boolean;
};

export const JobListContainer = (props: JobListContainerProps) => {
  const { children, isNoJobs } = props;

  if (isNoJobs) {
    return (
      <Flex w="full" my={4}>
        <Text size="sm">No jobs found.</Text>
      </Flex>
    );
  }

  return <>{children}</>;
};
