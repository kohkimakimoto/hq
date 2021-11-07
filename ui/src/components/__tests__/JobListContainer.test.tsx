import React from 'react';

import { JobListContainer } from '@/components/JobListContainer';
import { render, screen } from '@/test/test-utils';

test('render with no jobs', async () => {
  render(<JobListContainer isNoJobs={true}>should not render</JobListContainer>);
  // screen.debug();
  expect(screen.getByText('No jobs found.')).toBeInTheDocument();
});

test('render with jobs', async () => {
  render(<JobListContainer isNoJobs={false}>should render</JobListContainer>);
  expect(screen.getByText('should render')).toBeInTheDocument();
});
