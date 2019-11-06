import React, { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { Breadcrumb, Container, Header, Icon, Table } from 'semantic-ui-react';
import { useEffectDocumentTitle } from '../hooks/useEffectDocumentTitle';
import { useServices } from '../hooks/useService';
import { Job } from '../models/Job';

export const JobDetail: React.FC = () => {
  const { id } = useParams();

  useEffectDocumentTitle('Jobs');

  const { api, handleError } = useServices();
  const [job, setJob] = useState(new Job());

  useEffect(() => {
    (async () => {
      try {
        const job = await api.getJob(id!);
        setJob(job);
      } catch (err) {
        handleError(err);
      }
    })();
  }, []);

  return (
    <Container>
      <Breadcrumb>
        <Breadcrumb.Section as={Link} to={'/jobs'}>
          Jobs
        </Breadcrumb.Section>
        <Breadcrumb.Divider />
        <Breadcrumb.Section active>#{id}</Breadcrumb.Section>
      </Breadcrumb>

      <div className="page-title">
        <Header as="h1" dividing>
          #{job.id}
        </Header>
      </div>
      <div>
        <Header as="h2">ID</Header>
        <p>{job.id}</p>
        <Header as="h2">Name</Header>
        <p>{job.name}</p>
      </div>
    </Container>
  );
};
