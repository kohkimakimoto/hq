import React, { useEffect, useState } from 'react';
import { Link, useHistory } from 'react-router-dom';
import qs from 'query-string';
import { Container, Header, Icon, Table, Breadcrumb, Input, Button, Label, Segment, Form } from 'semantic-ui-react';
import { Controlled as CodeMirror } from 'react-codemirror2';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import { useServices } from '../hooks/useService';
import { JobList } from '../models/JobList';
import { useQuery } from '../hooks/useQuery';
import 'codemirror/mode/javascript/javascript';
import { useSelector } from 'react-redux';
import { StoreState } from '../store/State';
import { NewJob } from '../models/NewJob';

export const JobsNewScreen: React.FC = () => {
  useDocumentTitle('New Job');

  const { handleError, dispatcher, api } = useServices();
  const history = useHistory();
  const error = useSelector<StoreState, string>(state => state.error);
  const [uploading, setUploading] = useState<boolean>(false);
  const [job, setJob] = useState<NewJob>(new NewJob());

  const handlePushJob = () => {
    if (error != '') {
      dispatcher.commit({
        error: ''
      });
    }

    setUploading(true);

    (async () => {
      const postedJob = await api.createJob({
        name: name,
        comment: job.comment,
        url: job.url,
        payload: job.payload,
        headers: job.headers,
        timeout: job.timeout
      });
      setUploading(false);
      history.push('/jobs/' + postedJob.id);
    })().catch(err => {
      handleError(err);
      setUploading(false);
    });
  };

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
      <div>
        <Form>
          <Form.Field>
            <label>URL</label>
            <Input
              placeholder="https://your-worker-app-server/example"
              value={job.url}
              onChange={(e, data) => {
                setJob(job.modify({ url: data.value }));
              }}
            />
          </Form.Field>
          <Form.Field>
            <label>Name</label>
            <Input
              placeholder="example"
              value={job.name}
              onChange={(e, data) => {
                setJob(job.modify({ name: data.value }));
              }}
            />
          </Form.Field>
          <Form.Field>
            <label>Comment</label>
            <Input
              placeholder="This is an example job!"
              value={job.comment}
              onChange={(e, data) => {
                setJob(job.modify({ comment: data.value }));
              }}
            />
          </Form.Field>
          <Form.Field>
            <label>Payload</label>
            <CodeMirror
              value={job.payloadString}
              options={{
                mode: { name: 'javascript', json: true },
                theme: 'material',
                lineNumbers: true,
                smartIndent: false
              }}
              onBeforeChange={(editor, data, value) => {
                setJob(job.modify({ payloadString: value }));
              }}
            />
          </Form.Field>
          <Form.Field>
            <label>Headers</label>
            <CodeMirror
              value={job.headersString}
              options={{
                mode: { name: 'javascript', json: true },
                theme: 'material',
                lineNumbers: true,
                smartIndent: false
              }}
              onBeforeChange={(editor, data, value) => {
                setJob(job.modify({ headersString: value }));
              }}
            />
          </Form.Field>
          <Form.Field>
            <label>Timeout</label>
            <Input
              type="number"
              placeholder="0"
              value={job.timeoutString}
              onChange={(e, data) => {
                setJob(job.modify({ timeoutString: data.value }));
              }}
            />
          </Form.Field>
        </Form>
        <div style={{ marginTop: 20 }}>
          <Button color="teal" content="Push job" onClick={handlePushJob} disabled={job.url == ''} fluid loading={uploading} />
        </div>
      </div>
    </Container>
  );
};
