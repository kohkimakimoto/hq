import React, { useEffect, useState } from 'react';
import { Container, Header, Segment, Progress, Breadcrumb } from 'semantic-ui-react';
import { Stats } from '../models/Stats';
import { useEffectDocumentTitle } from '../hooks/useEffectDocumentTitle';
import { useSelector } from 'react-redux';
import { StoreState } from '../store/State';
import { useEffectStats } from '../hooks/useEffectStats';

export const StatsScreen: React.FC = () => {
  useEffectStats();
  useEffectDocumentTitle('Stats');

  const stats = useSelector<StoreState, Stats>(state => state.stats);

  return (
    <Container>
      <Breadcrumb>
        <Breadcrumb.Section active>Stats</Breadcrumb.Section>
      </Breadcrumb>

      <div className="page-title">
        <Header as="h1" dividing>
          Stats
        </Header>
      </div>
      <div>
        <Header as="h2">Queues</Header>
        <p>{stats.queues}</p>
        <Header as="h2">Dispatchers</Header>
        <p>{stats.dispatchers}</p>
        <Header as="h2">MaxWorkers</Header>
        <p>{stats.maxWorkers}</p>
        <Header as="h2">QueueUsage</Header>
        <p>{stats.queueUsage}</p>
        <Header as="h2">NumWaitingJobs</Header>
        <p>{stats.numWaitingJobs}</p>
        <Header as="h2">NumRunningJobs</Header>
        <p>{stats.numRunningJobs}</p>
        <Header as="h2">NumWorkers</Header>
        <p>{stats.numWorkers}</p>
        <Header as="h2">NumJobs</Header>
        <p>{stats.numJobs}</p>
      </div>
    </Container>
  );
};
