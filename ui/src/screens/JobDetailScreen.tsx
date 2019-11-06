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
        <Header as="h2">URL</Header>
        <p>{job.url}</p>
        <Header as="h2">Comment</Header>
        <p>{job.comment}</p>
        <Header as="h2">Payload</Header>
        <pre><code>{JSON.stringify(job.payload)}</code></pre>
        <Header as="h2">Headers</Header>
        <pre><code>{JSON.stringify(job.headers)}</code></pre>
        <Header as="h2">Timeout</Header>
        <p>{job.timeout}</p>
        <Header as="h2">CreatedAt</Header>
        <p>{job.createdAt}</p>
        <Header as="h2">StartedAt</Header>
        <p>{job.startedAt}</p>
        <Header as="h2">FinishedAt</Header>
        <p>{job.finishedAt}</p>
        <Header as="h2">StatusCode</Header>
        <p>{job.statusCode}</p>
        <Header as="h2">Output</Header>
        <pre><code>{job.output}</code></pre>
        <Header as="h2">Error</Header>
        <pre><code>{job.err}</code></pre>
        <Header as="h2">Status</Header>
        <p>{job.status}</p>
      </div>
    </Container>
  );
};
