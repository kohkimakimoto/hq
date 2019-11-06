import React, { useEffect, useState } from 'react';
import { Link, useHistory } from 'react-router-dom';
import qs from 'query-string';
import { Container, Header, Icon, Table, Breadcrumb, Input, Button } from 'semantic-ui-react';
import { useEffectDocumentTitle } from '../hooks/useEffectDocumentTitle';
import { useServices } from '../hooks/useService';
import { JobList } from '../models/JobList';
import { useQuery } from '../hooks/useQuery';

export const JobsScreen: React.FC = () => {
  useEffectDocumentTitle('Jobs');

  const history = useHistory();
  const query = useQuery();

  let term = query.get('term');

  const { api, handleError } = useServices();
  const [jobList, setJobList] = useState(new JobList());
  const [searchText, setSearchText] = useState(term ? term : '');

  const handleChangeSearchText = (e, { value }) => {
    setSearchText(value);
  };

  const handleKeyDown = e => {
    if (e.keyCode === 13) {
      history.push(
        '/jobs?' +
          qs.stringify({
            term: searchText
          })
      );
    }
  };

  const loadJobs = () => {
    (async () => {
      const listJobs = await api.listJobs({
        term: term === null ? undefined : term,
        reverse: true,
        limit: 100
      });
      setJobList(listJobs);
    })().catch(err => handleError(err));
  };

  useEffect(() => {
    loadJobs();
  }, [term]);

  return (
    <Container>
      <Breadcrumb>
        <Breadcrumb.Section active>Jobs</Breadcrumb.Section>
      </Breadcrumb>

      <div className="page-title">
        <Header as="h1" dividing>
          Jobs
        </Header>
      </div>

      <div>
        <Input
          placeholder="name:default status:failure search-term"
          icon={
            <Icon
              name="search"
              link
              color="blue"
              onClick={() => {
                history.push(
                  '/jobs?' +
                    qs.stringify({
                      term: searchText
                    })
                );
              }}
            />
          }
          fluid
          value={searchText}
          onChange={handleChangeSearchText}
          onKeyDown={handleKeyDown}
        />
      </div>

      <Table basic="very">
        <Table.Header>
          <Table.Row>
            <Table.HeaderCell colSpan="3"></Table.HeaderCell>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {(() => {
            return jobList.jobs.map((job, index, array) => {
              return (
                <Table.Row key={'job_' + job.id}>
                  <Table.Cell verticalAlign="top">
                    <div>
                      <Header size="tiny" as={Link} to={'/jobs/' + job.id} color={job.statusColor}>
                        {'#' + job.id}
                      </Header>
                    </div>
                    <div>{job.name}</div>
                    <div>{job.url}</div>
                    <div className="text grey">{job.comment}</div>
                  </Table.Cell>
                  <Table.Cell collapsing verticalAlign="top">
                    <Header size="tiny" as="span" to={'/jobs/' + job.id} color={job.statusColor}>
                      {job.status}
                    </Header>
                  </Table.Cell>
                  <Table.Cell collapsing verticalAlign="top">
                    <div>
                      <Icon color="grey" name="calendar outline" />
                      {job.createdAtFromNow}
                    </div>
                    {(() => {
                      if (job.startedAtFromNow != '') {
                        return (
                          <div>
                            <Icon color="grey" name="calendar check" />
                            {job.startedAtFromNow}
                          </div>
                        );
                      }
                    })()}
                    {(() => {
                      if (job.duration != '') {
                        return (
                          <div>
                            <Icon color="grey" name="clock outline" />
                            {job.duration}
                          </div>
                        );
                      }
                    })()}
                  </Table.Cell>
                </Table.Row>
              );
            });
          })()}
        </Table.Body>
      </Table>
    </Container>
  );
};
