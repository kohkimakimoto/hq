import { VStack, Code, Box, Text, Badge, Heading, Breadcrumb, BreadcrumbItem, BreadcrumbLink, Flex, Spinner } from '@chakra-ui/react';
import React from 'react';
import { Helmet } from 'react-helmet-async';
import { Link as RouterLink, useParams } from 'react-router-dom';

import { useGetJob } from '@/api/api';
import { useStatusColors } from '@/models/Status';

export const JobDetailPage = () => {
  const statusColors = useStatusColors();
  const { id } = useParams();
  const [{ job, validationError }] = useGetJob(id!, { useCache: false });

  if (validationError.hasError) {
    return (
      <Flex alignItems="center" justify="center" py={4} w="full">
        <Text>{validationError.message}</Text>
      </Flex>
    );
  }

  if (!job) {
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
        <BreadcrumbItem>
          <BreadcrumbLink as={RouterLink} to="/">
            Jobs
          </BreadcrumbLink>
        </BreadcrumbItem>
        <BreadcrumbItem isCurrentPage>
          <BreadcrumbLink as={RouterLink} to={`/jobs/${job.id}`}>
            {job.id}
          </BreadcrumbLink>
        </BreadcrumbItem>
      </Breadcrumb>
      <Heading
        as="h1"
        size="lg"
        pb={2}
        w="full"
        color={statusColors[job.status]}
        borderBottom="2px"
        borderBottomColor={statusColors[job.status]}
      >
        {job.id}
      </Heading>
      <Badge textColor="white" bg={statusColors[job.status]}>
        {job.status}
      </Badge>
      <VStack my="4" alignItems="flex-start" spacing={4}>
        <Box>
          <Heading as="h2" size="sm">
            Name
          </Heading>
          <Text>{job.name}</Text>
        </Box>
        <Box>
          <Heading as="h2" size="sm">
            URL
          </Heading>
          <Code>{job.url}</Code>
        </Box>
        <Box>
          <Heading as="h2" size="sm">
            Comment
          </Heading>
          <Text>{job.comment}</Text>
        </Box>
        <Box>
          <Heading as="h2" size="sm">
            Payload
          </Heading>
          <Code>{JSON.stringify(job.payload, null, '\t')}</Code>
        </Box>
        <Box>
          <Heading as="h2" size="sm">
            Headers
          </Heading>
          <Code>{JSON.stringify(job.headers, null, '\t')}</Code>
        </Box>
        <Box>
          <Heading as="h2" size="sm">
            Timeout
          </Heading>
          <Text>{job.timeout}</Text>
        </Box>
        <Box>
          <Heading as="h2" size="sm">
            CreatedAt
          </Heading>
          <Text>{job.createdAt.format('YYYY-MM-DD HH:mm:ss')}</Text>
        </Box>
        <Box>
          <Heading as="h2" size="sm">
            StartedAt
          </Heading>
          <Text>{job.startedAt ? job.startedAt.format('YYYY-MM-DD HH:mm:ss') : ''}</Text>
        </Box>
        <Box>
          <Heading as="h2" size="sm">
            FinishedAt
          </Heading>
          <Text>{job.finishedAt ? job.finishedAt.format('YYYY-MM-DD HH:mm:ss') : ''}</Text>
        </Box>
        <Box>
          <Heading as="h2" size="sm">
            Status Code
          </Heading>
          <Text>{job.statusCode}</Text>
        </Box>
        <Box>
          <Heading as="h2" size="sm">
            Output
          </Heading>
          <Code>{job.output}</Code>
        </Box>
        <Box>
          <Heading as="h2" size="sm">
            Error
          </Heading>
          <Code>{job.err}</Code>
        </Box>
      </VStack>
    </>
  );
};
