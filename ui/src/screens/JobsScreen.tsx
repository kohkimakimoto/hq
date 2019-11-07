import React, { useEffect, useState } from 'react';
import { Link, useHistory } from 'react-router-dom';
import qs from 'query-string';
import { Container, Header, Icon, Table, Breadcrumb, Input, Button, Modal } from 'semantic-ui-react';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import { useServices } from '../hooks/useService';
import { JobList } from '../models/JobList';
import { useQuery } from '../hooks/useQuery';
import { Job } from '../models/Job';
import {useInterval} from "../hooks/useInterval";

const limit = 300;
const pollingInterval = 5000;

export const JobsScreen: React.FC = () => {
  // hooks
  useDocumentTitle('Jobs');
  const history = useHistory();
  const query = useQuery();
  const term = query.get('term') ? (query.get('term') as string) : '';
  const { api, handleError } = useServices();

  // state
  const [jobList, setJobList] = useState(new JobList());
  const [searchText, setSearchText] = useState(term);
  // restarting
  const [restartingJob, setRestartingJob] = useState<Job | null>(null);
  const [restarting, setRestarting] = useState<boolean>(false);
  // deleting
  const [deletingJob, setDeletingJob] = useState<Job | null>(null);
  const [deleting, setDeleting] = useState<boolean>(false);

  // effect
  useEffect(() => {
    refreshList(term);
  }, [term]);

  useInterval(() => {
    // polling
    refreshList(term);
  }, pollingInterval);

  const refreshList = (term: string) => {
    (async () => {
      const list = await api.listJobs({
        term: term,
        reverse: true,
        limit: limit
      });
      setJobList(list);
    })().catch(err => handleError(err));
  };

  const handleChangeSearchText = (e, { value }) => {
    setSearchText(value);
  };

  const handleKeyDown = e => {
    if (e.keyCode === 13) {
      search(searchText);
    }
  };

  const search = (searchText: string) => {
    if (searchText != term) {
      history.push('/jobs?' + qs.stringify({ term: searchText }));
    } else {
      refreshList(searchText);
    }
  };

  const handleClickStop = (job: Job) => {

  };

  const handleClickRestart = (job: Job) => {
    setRestartingJob(job);
  };

  const handleClickDelete = (job: Job) => {
    setDeletingJob(job);
  };

  const handleRestartAsANewJob = (job: Job) => {
    setRestarting(true);

    (async () => {
      const resp = await api.restartJob(job.id, true);
      refreshList(term);
    })()
      .catch(err => {
        handleError(err);
      })
      .finally(() => {
        setRestarting(false);
        setRestartingJob(null);
      });
  };

  const handleRestart = (job: Job) => {
    setRestarting(true);

    (async () => {
      const resp = await api.restartJob(job.id, false);
      refreshList(term);
    })()
      .catch(err => {
        handleError(err);
      })
      .finally(() => {
        setRestarting(false);
        setRestartingJob(null);
      });
  };

  const handleDelete = (job: Job) => {
    setDeleting(true);

    (async () => {
      const resp = await api.deleteJob(job.id);
      refreshList(term);
    })()
      .catch(err => {
        handleError(err);
      })
      .finally(() => {
        setDeleting(false);
        setDeletingJob(null);
      });
  };

  return (
    <React.Fragment>
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
              <Icon name="search" link color="blue" onClick={() => search(searchText)} />
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
              <Table.HeaderCell colSpan="4">
                Listed {jobList.count} job(s)
                {(() => {
                  if (jobList.hasNext) {
                    return (
                      <span className="text yellow" style={{ marginLeft: 10 }}>
                        HQ has more jobs that could not be displayed here.
                      </span>
                    );
                  }
                })()}
              </Table.HeaderCell>
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
                      <Header size="tiny" as="div" to={'/jobs/' + job.id} color={job.statusColor} textAlign='right' style={{minWidth: 70}}>
                        {(() => {
                          if (job.status == 'running') {
                            return (<Icon name="spinner" loading/>);
                          }
                        })()}
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
                    <Table.Cell collapsing verticalAlign="top">
                      {(() => {
                        if (job.status == 'running' || job.status == "waiting") {
                          return (
                            <div style={{ marginBottom: 10, width: 70 }}>
                              <Button
                                basic
                                size="mini"
                                fluid
                                compact
                                content="Stop"
                                color="orange"
                                onClick={() => handleClickStop(job)}
                              />
                            </div>
                          );
                        } else {
                          return (
                            <React.Fragment>
                              <div style={{ marginBottom: 10, width: 70 }}>
                                <Button
                                  basic
                                  size="mini"
                                  fluid
                                  compact
                                  content="Restart"
                                  color="teal"
                                  onClick={() => handleClickRestart(job)}
                                />
                              </div>
                              <div>
                                <Button
                                  basic
                                  size="mini"
                                  fluid
                                  compact
                                  content="Delete"
                                  color="red"
                                  onClick={() => handleClickDelete(job)}
                                />
                              </div>
                            </React.Fragment>
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
      {(() => {
        if (restartingJob) {
          return (
            <Modal size="tiny" open={!!restartingJob} onClose={() => setRestartingJob(null)}>
              <Modal.Header>Restarting Job</Modal.Header>
              <Modal.Content>
                <p>Are you sure you want to restart the following job?</p>
                <div>
                  <Header size="tiny" color={restartingJob.statusColor}>
                    {'#' + restartingJob.id}
                  </Header>
                </div>
              </Modal.Content>
              <Modal.Actions>
                <Button
                  color="teal"
                  content="Restart as a new job"
                  onClick={() => handleRestartAsANewJob(restartingJob)}
                  loading={restarting}
                />
                <Button color="teal" content="Restart" onClick={() => handleRestart(restartingJob)} loading={restarting} />
              </Modal.Actions>
            </Modal>
          );
        }
      })()}

      {(() => {
        if (deletingJob) {
          return (
            <Modal size="tiny" open={!!deletingJob} onClose={() => setDeletingJob(null)}>
              <Modal.Header>Deleting Job</Modal.Header>
              <Modal.Content>
                <p>Are you sure you want to delete the following job?</p>
                <div>
                  <Header size="tiny" color={deletingJob.statusColor}>
                    {'#' + deletingJob.id}
                  </Header>
                </div>
              </Modal.Content>
              <Modal.Actions>
                <Button color="red" content="Delete" onClick={() => handleDelete(deletingJob)} loading={deleting} />
              </Modal.Actions>
            </Modal>
          );
        }
      })()}
    </React.Fragment>
  );
};
