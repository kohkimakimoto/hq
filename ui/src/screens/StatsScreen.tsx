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
          <Header as="h2" size="small">
            Queues
          </Header>
          <p>{stats.queues}</p>
        </Segment>
        <Segment vertical>
          <Header as="h2" size="small">
            Dispatchers
          </Header>
          <p>{stats.dispatchers}</p>
        </Segment>
        <Segment vertical>
          <Header as="h2" size="small">
            MaxWorkers
          </Header>
          <p>{stats.maxWorkers}</p>
        </Segment>
        <Segment vertical>
          <Header as="h2" size="small">
            QueueUsage
          </Header>
          <p>{stats.queueUsage}</p>
        </Segment>
        <Segment vertical>
          <Header as="h2" size="small">
            NumWaitingJobs
          </Header>
          <p>{stats.numWaitingJobs}</p>
        </Segment>
        <Segment vertical>
          <Header as="h2" size="small">
            NumRunningJobs
          </Header>
          <p>{stats.numRunningJobs}</p>
        </Segment>
        <Segment vertical>
          <Header as="h2" size="small">
            NumWorkers
          </Header>
          <p>{stats.numWorkers}</p>
        </Segment>
        <Segment vertical>
          <Header as="h2" size="small">
            NumJobs
          </Header>
          <p>{stats.numJobs}</p>
        </Segment>
      </div>
    </Container>
  );
};
