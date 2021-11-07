import { SearchIcon } from '@chakra-ui/icons';
import {
  VStack,
  Input,
  InputGroup,
  InputLeftElement,
  SimpleGrid,
  Stat,
  StatLabel,
  StatNumber,
  StatHelpText,
  Text,
  Stack,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Link,
  Spinner,
  Flex,
  useColorModeValue,
  HStack,
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  useToast,
} from '@chakra-ui/react';
import React from 'react';
import { Helmet } from 'react-helmet-async';
import { Link as RouterLink } from 'react-router-dom';
import { useQueryParam, StringParam } from 'use-query-params';

import { useGetDashboard } from '@/api/api';
import { JobActions } from '@/components/JobActions';
import { JobListContainer } from '@/components/JobListContainer';
import { useInterval } from '@/lib/useInterval';
import { useStatusColors } from '@/models/Status';

export const IndexPage = () => {
  const lineColor = useColorModeValue('gray.200', 'gray.600');
  const searchIconColor = useColorModeValue('gray.200', 'gray.600');
  const warningTextColor = useColorModeValue('yellow.500', 'yellow.100');
  const tableHeaderColor = useColorModeValue('gray.50', 'gray.800');
  const tableTextColor = useColorModeValue('gray.800', 'gray.50');
  const toast = useToast();

  const statusColors = useStatusColors();
  const [term, setTerm] = useQueryParam('term', StringParam);
  const [{ dashboard }, executeGetDashboard] = useGetDashboard({ term: term, reverse: true });

  useInterval(executeGetDashboard, 2000);

  if (!dashboard) {
    return (
      <Flex alignItems="center" justify="center" py={4} w="full">
        <Spinner size={'md'} />
      </Flex>
    );
  }

  return (
    <>
      <Helmet>
        <title>Dashboard | HQ</title>
      </Helmet>
      <Breadcrumb mb={4}>
        <BreadcrumbItem isCurrentPage>
          <BreadcrumbLink as={RouterLink} to="/">
            Jobs
          </BreadcrumbLink>
        </BreadcrumbItem>
      </Breadcrumb>
      <SimpleGrid columns={{ base: 1, md: 2, lg: 4 }} spacingX={{ base: 0, md: 2 }} spacingY={2}>
        <Stat borderWidth="1px" borderColor={lineColor} rounded="sm" p={2}>
          <StatLabel>Jobs Running</StatLabel>
          <StatNumber>{dashboard.stats.numJobsRunning}</StatNumber>
          <StatHelpText>Max Concurrent Workers: {dashboard.stats.maxConcurrentWorkers}</StatHelpText>
        </Stat>
        <Stat borderWidth="1px" borderColor={lineColor} rounded="sm" p={2}>
          <StatLabel>Jobs in Queue</StatLabel>
          <StatNumber>{dashboard.stats.numJobsInQueue}</StatNumber>
          <StatHelpText>
            Queue Usage: {dashboard.stats.queueUsageRate}% ({dashboard.stats.numJobsInQueue}/{dashboard.stats.queues})
          </StatHelpText>
        </Stat>
        <Stat borderWidth="1px" borderColor={lineColor} rounded="sm" p={2}>
          <StatLabel>Total Stored Jobs</StatLabel>
          <StatNumber>{dashboard.stats.numStoredJobs}</StatNumber>
          <StatHelpText>&nbsp;</StatHelpText>
        </Stat>
        <Stat borderWidth="1px" borderColor={lineColor} rounded="sm" p={2}>
          <StatLabel>New Jobs in Last Minute</StatLabel>
          <StatNumber>{dashboard.stats.numJobsInLastMinute}</StatNumber>
          <StatHelpText>&nbsp;</StatHelpText>
        </Stat>
      </SimpleGrid>
      <VStack mt={10} w="full">
        <InputGroup>
          <InputLeftElement pointerEvents="none">
            <SearchIcon color={searchIconColor} />
          </InputLeftElement>
          <Input placeholder="search-term" size="md" w="full" onChange={(e) => setTerm(e.target.value)} value={term == null ? '' : term} />
        </InputGroup>
      </VStack>

      <Stack mt={4} w="full">
        <Text size="xs">
          Listed <span style={{ fontWeight: 'bold' }}>{dashboard.jobList.count}</span> job(s)
          <Text display={dashboard.jobList.hasNext ? 'inline' : 'none'} as="span" size="xs" ml={2} color={warningTextColor}>
            HQ has more jobs that could not be displayed here.
          </Text>
        </Text>
      </Stack>
      <JobListContainer isNoJobs={dashboard.jobList.count === 0}>
        <Table size="sm" mt={4} w="full">
          <Thead bgColor={tableHeaderColor}>
            <Tr>
              <Th>ID</Th>
              <Th>Name</Th>
              <Th>URL</Th>
              <Th>Created</Th>
              <Th>Started</Th>
              <Th>Finished</Th>
              <Th>Duration</Th>
              <Th>Status</Th>
              <Th>Actions</Th>
            </Tr>
          </Thead>
          <Tbody>
            {dashboard.jobList.jobs.map((job) => {
              return (
                <Tr key={`job-${job.id}`} textColor={tableTextColor}>
                  <Td w={1}>
                    <Link
                      fontWeight="bold"
                      color={statusColors[job.status]}
                      as={RouterLink}
                      to={`/jobs/${job.id}`}
                      h="full"
                      display="flex"
                      alignItems="center"
                      _focus={{ outline: 'none' }}
                    >
                      {job.id}
                    </Link>
                  </Td>
                  <Td>{job.name}</Td>
                  <Td>{job.url}</Td>
                  <Td>{job.createdAt.format('YYYY-MM-DD HH:mm:ss')}</Td>
                  <Td>{job.startedAt ? job.startedAt.format('YYYY-MM-DD HH:mm:ss') : ''}</Td>
                  <Td>{job.finishedAt ? job.finishedAt.format('YYYY-MM-DD HH:mm:ss') : ''}</Td>
                  <Td>{job.duration}</Td>
                  <Td>
                    <HStack spacing={2}>
                      {(() => {
                        if (job.status == 'running') {
                          return <Spinner size="xs" color={statusColors[job.status]} />;
                        }
                      })()}
                      <Text fontWeight="bold" color={statusColors[job.status]}>
                        {job.status}
                      </Text>
                    </HStack>
                  </Td>
                  <Td>
                    <JobActions
                      job={job}
                      onStop={(job) => {
                        toast({
                          title: 'Job stopped',
                          description: `We stopped the job: ${job.id}`,
                          status: 'success',
                          position: 'top-right',
                          duration: 3000,
                          isClosable: true,
                        });
                        executeGetDashboard();
                      }}
                      onRestartAsNew={(job) => {
                        toast({
                          title: 'Job restarted',
                          description: `We restarted the job: ${job.id}`,
                          status: 'success',
                          position: 'top-right',
                          duration: 3000,
                          isClosable: true,
                        });
                        executeGetDashboard();
                      }}
                      onRestart={(job) => {
                        toast({
                          title: 'Job restarted',
                          description: `We restarted the job: ${job.id}`,
                          status: 'success',
                          position: 'top-right',
                          duration: 3000,
                          isClosable: true,
                        });
                        executeGetDashboard();
                      }}
                      onDelete={(job) => {
                        toast({
                          title: 'Job deleted',
                          description: `We deleted the job: ${job.id}`,
                          status: 'success',
                          position: 'top-right',
                          duration: 3000,
                          isClosable: true,
                        });
                        executeGetDashboard();
                      }}
                    />
                  </Td>
                </Tr>
              );
            })}
          </Tbody>
        </Table>
      </JobListContainer>
    </>
  );
};
