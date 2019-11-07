import React, { useEffect, useState } from 'react';
import { Link, useHistory } from 'react-router-dom';
import qs from 'query-string';
import { Container, Header, Icon, Table, Breadcrumb, Input, Button, Label } from 'semantic-ui-react';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import { useServices } from '../hooks/useService';
import { JobList } from '../models/JobList';
import { useQuery } from '../hooks/useQuery';

export const JobsNewScreen: React.FC = () => {
  return (
    <Container>
      <Breadcrumb>
        <Breadcrumb.Section as={Link} to={'/jobs'}>
          Jobs
        </Breadcrumb.Section>
        <Breadcrumb.Divider />
        <Breadcrumb.Section active>New Job</Breadcrumb.Section>
      </Breadcrumb>

      <div className="page-title">
        <Header as="h1" dividing>
          New Job
        </Header>
      </div>
      <div></div>
    </Container>
  );
};
