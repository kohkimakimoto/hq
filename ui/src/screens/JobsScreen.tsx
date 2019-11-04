import React, {useEffect, useState} from 'react';
import { Container, Header, Table } from 'semantic-ui-react';
import { useEffectDocumentTitle } from '../hooks/useEffectDocumentTitle';
import {useServices} from "../hooks/useService";
import {JobList} from "../models/JobList";

export const JobsScreen: React.FC = () => {
  useEffectDocumentTitle('Jobs');

  const { api, handleError } = useServices();
  const [ jobList, setJobList] = useState(new JobList());

  useEffect(() => {
    (async () => {
      try {
        const listJobs = await api.listJobs();
        setJobList(listJobs);
      } catch (err) {
        handleError(err);
      }
    })();
  }, []);

  return (
    <Container>
      <div className="page-title">
        <Header as="h1" dividing>
          Jobs
        </Header>
      </div>

      <Table basic='very'>
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell>ID</Table.HeaderCell>
            <Table.HeaderCell>Name</Table.HeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {(() => {
            return jobList.jobs.map((job, index, array) => {
              return (
                <Table.Row key={'job_' + job.id}>
                  <Table.Cell>{job.id}</Table.Cell>
                  <Table.Cell>{job.name}</Table.Cell>
                </Table.Row>
              );
            });
          })()}
        </Table.Body>
      </Table>
    </Container>
  );
};
