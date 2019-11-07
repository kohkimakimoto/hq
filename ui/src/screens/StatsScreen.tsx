import React, { useEffect, useState } from 'react';
import { Container, Header, Segment, Progress, Breadcrumb } from 'semantic-ui-react';
import { Stats } from '../models/Stats';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import { useSelector } from 'react-redux';
import { StoreState } from '../store/State';
import { useUpdateStats } from '../hooks/useUpdateStats';

export const StatsScreen: React.FC = () => {
  useUpdateStats();
  useDocumentTitle('Stats');

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
        <Segment vertical>
          <Header as="h2" size="medium">
            Queues
          </Header>
          <p>{stats.queues}</p>
        </Segment>
        <Segment vertical>
          <Header as="h2" size="medium">
            Dispatchers
          </Header>
          <p>{stats.dispatchers}</p>
        </Segment>
        <Segment vertical>
          <Header as="h2" size="medium">
            MaxWorkers
          </Header>
          <p>{stats.maxWorkers}</p>
        </Segment>
        <Segment vertical>
          <Header as="h2" size="medium">
            QueueUsage
          </Header>
          <p>{stats.queueUsage}</p>
        </Segment>
        <Segment vertical>
          <Header as="h2" size="medium">
            NumWaitingJobs
          </Header>
          <p>{stats.numWaitingJobs}</p>
        </Segment>
        <Segment vertical>
          <Header as="h2" size="medium">
            NumRunningJobs
          </Header>
          <p>{stats.numRunningJobs}</p>
        </Segment>
        <Segment vertical>
          <Header as="h2" size="medium">
            NumWorkers
          </Header>
          <p>{stats.numWorkers}</p>
        </Segment>
        <Segment vertical>
          <Header as="h2" size="medium">
            NumJobs
          </Header>
          <p>{stats.numJobs}</p>
        </Segment>
      </div>
    </Container>
  );
};
