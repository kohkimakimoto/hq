import React, { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { Breadcrumb, Container, Header, Icon, Table, Label, Segment } from 'semantic-ui-react';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import { useServices } from '../hooks/useService';
import { Job } from '../models/Job';

export const JobDetail: React.FC = () => {
  const { id } = useParams();

  useDocumentTitle('Jobs');

  const { api, handleError } = useServices();
  const [job, setJob] = useState(new Job());

  useEffect(() => {
    (async () => {
      const job = await api.getJob(id!);
      setJob(job);
    })().catch(err => {
      handleError(err);
    });
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
      {(() => {
        if (job.id !== '') {
          return (
            <React.Fragment>
              <div className="page-title">
                <Header as="h1" dividing color={job.statusColor}>
                  #{job.id}
                </Header>
                <Label color={job.statusColor}>{job.status}</Label>
              </div>

              <div>
                <Segment vertical>
                  <Header as="h2" size="small">
                    Name
                  </Header>
                  <p>{job.name}</p>
                </Segment>
                <Segment vertical>
                  <Header as="h2" size="small">
                    URL
                  </Header>
                  <pre>
                    <code>{job.url}</code>
                  </pre>
                </Segment>
                <Segment vertical>
                  <Header as="h2" size="small">
                    Comment
                  </Header>
                  <p>{job.comment}</p>
                </Segment>
                <Segment vertical>
                  <Header as="h2" size="small">
                    Payload
                  </Header>
                  <pre>
                    <code>{JSON.stringify(job.payload, null, '\t')}</code>
                  </pre>
                </Segment>
                <Segment vertical>
                  <Header as="h2" size="small">
                    Headers
                  </Header>
                  <pre>
                    <code>{JSON.stringify(job.headers, null, '\t')}</code>
                  </pre>
                </Segment>
                <Segment vertical>
                  <Header as="h2" size="small">
                    Timeout
                  </Header>
                  <p>{job.timeout}</p>
                </Segment>
                <Segment vertical>
                  <Header as="h2" size="small">
                    CreatedAt
                  </Header>
                  <p>{job.createdAt}</p>
                </Segment>
                <Segment vertical>
                  <Header as="h2" size="small">
                    StartedAt
                  </Header>
                  <p>{job.startedAt}</p>
                </Segment>
                <Segment vertical>
                  <Header as="h2" size="small">
                    FinishedAt
                  </Header>
                  <p>{job.finishedAt}</p>
                </Segment>
                <Segment vertical>
                  <Header as="h2" size="small">
                    StatusCode
                  </Header>
                  <p>{job.statusCode}</p>
                </Segment>
                <Segment vertical>
                  <Header as="h2" size="small">
                    Output
                  </Header>
                  <pre>
                    <code>{job.output}</code>
                  </pre>
                </Segment>
                <Segment vertical>
                  <Header as="h2" size="small">
                    Error
                  </Header>
                  <pre>
                    <code>{job.err}</code>
                  </pre>
                </Segment>
              </div>
            </React.Fragment>
          );
        }
      })()}
    </Container>
  );
};
