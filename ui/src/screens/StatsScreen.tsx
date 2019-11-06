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
        <Segment>
          <p>{stats.queues}</p>
        </Segment>
        <Header as="h2">Dispatchers</Header>
        <Segment>
          <p>{stats.dispatchers}</p>
        </Segment>
        <Header as="h2">MaxWorkers</Header>
        <Segment>
          <p>{stats.maxWorkers}</p>
        </Segment>
        <Header as="h2">QueueUsage</Header>
        <Segment>
          <p>{stats.queueUsage}</p>
        </Segment>
        <Header as="h2">NumWaitingJobs</Header>
        <Segment>
          <p>{stats.numWaitingJobs}</p>
        </Segment>
        <Header as="h2">NumRunningJobs</Header>
        <Segment>
          <p>{stats.numRunningJobs}</p>
        </Segment>
        <Header as="h2">NumWorkers</Header>
        <Segment>
          <p>{stats.numWorkers}</p>
        </Segment>
        <Header as="h2">NumJobs</Header>
        <Segment>
          <p>{stats.numJobs}</p>
        </Segment>
      </div>
    </Container>
  );
};
