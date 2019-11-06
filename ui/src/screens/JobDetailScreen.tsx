import React, { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import { Breadcrumb, Container, Header, Icon, Table, Label, Segment } from 'semantic-ui-react';
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
                <Header as="h2">ID</Header>
                <Segment>
                  <p>{'#' + job.id}</p>
                </Segment>
                <Header as="h2">Name</Header>
                <Segment>
                  <p>{job.name}</p>
                </Segment>
                <Header as="h2">URL</Header>
                <Segment>
                  <p>{job.url}</p>
                </Segment>
                <Header as="h2">Comment</Header>
                <Segment>
                  <p>{job.comment}</p>
                </Segment>
                <Header as="h2">Payload</Header>
                <Segment>
                  <pre>
                    <code>{JSON.stringify(job.payload, null, '\t')}</code>
                  </pre>
                </Segment>
                <Header as="h2">Headers</Header>
                <Segment>
                  <pre>
                    <code>{JSON.stringify(job.headers, null, '\t')}</code>
                  </pre>
                </Segment>
                <Header as="h2">Timeout</Header>
                <Segment>
                  <p>{job.timeout}</p>
                </Segment>
                <Header as="h2">CreatedAt</Header>
                <Segment>
                  <p>{job.createdAt}</p>
                </Segment>
                <Header as="h2">StartedAt</Header>
                <Segment>
                  <p>{job.startedAt}</p>
                </Segment>
                <Header as="h2">FinishedAt</Header>
                <Segment>
                  <p>{job.finishedAt}</p>
                </Segment>
                <Header as="h2">StatusCode</Header>
                <Segment>
                  <p>{job.statusCode}</p>
                </Segment>
                <Header as="h2">Output</Header>
                <Segment>
                  <pre>
                    <code>{job.output}</code>
                  </pre>
                </Segment>
                <Header as="h2">Error</Header>
                <Segment>
                  <pre>
                    <code>{job.err}</code>
                  </pre>
                </Segment>
                <Header as="h2">Status</Header>
                <Segment>
                  <p className={'text ' + job.statusColor}>{job.status}</p>
                </Segment>
              </div>
            </React.Fragment>
          );
        }
      })()}
    </Container>
  );
};
