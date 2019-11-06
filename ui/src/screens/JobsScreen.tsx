import React, { useEffect, useState } from 'react';
import { Link, useHistory } from 'react-router-dom';
import qs from 'query-string';
import { Container, Header, Icon, Table, Breadcrumb, Input, Button } from 'semantic-ui-react';
import { useEffectDocumentTitle } from '../hooks/useEffectDocumentTitle';
import { useServices } from '../hooks/useService';
import { JobList } from '../models/JobList';
import { useQuery } from '../hooks/useQuery';

const limit = 100;
export const JobsScreen: React.FC = () => {
  useEffectDocumentTitle('Jobs');

  const history = useHistory();
  const query = useQuery();
  const term = query.get('term') ? query.get('term') as string : '';

  const { api, handleError } = useServices();
  const [ jobList, setJobList ] = useState(new JobList());
  const [ searchText, setSearchText ] = useState(term);

  const handleChangeSearchText = (e, { value }) => {
    setSearchText(value);
  };

  const handleClickMore = () => {
    (async () => {
      if (jobList.next) {
        const newList = await api.listJobs({
          term: term,
          reverse: true,
          limit: limit,
          begin: jobList.next,
        });

        const newJobs = jobList.jobs.concat(newList.jobs);

        setJobList(jobList.modify({
          jobs: newJobs,
          hasNext: newList.hasNext,
          next: newList.next,
          count: jobList.count + newList.count,
        }));
      }
    })().catch(err => handleError(err));
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

  useEffect(() => {
    (async () => {
      const list = await api.listJobs({
        term: term,
        reverse: true,
        limit: limit,
      });
      setJobList(list);
    })().catch(err => handleError(err));
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
          placeholder="search-term"
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
      {(() => {
        if (jobList.hasNext) {
          return (
            <div>
              <Button fluid onClick={handleClickMore}>More</Button>
            </div>
          );
        }
      })()}
    </Container>
  );
};
