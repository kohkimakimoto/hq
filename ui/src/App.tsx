import React, { useContext, useEffect } from 'react';
import { BrowserRouter as Router, Switch, Route, Link, Redirect, useLocation } from 'react-router-dom';
import {
  Container,
  Divider,
  Dropdown,
  Grid,
  Header,
  Image,
  List,
  Menu,
  Segment,
  Message,
  SemanticCOLORS,
  Icon
} from 'semantic-ui-react';
import { JobsScreen } from './screens/JobsScreen';
import { NotFoundScreen } from './screens/NotFoundScreen';
import { StatsScreen } from './screens/StatsScreen';
import { Store } from 'redux';
import { Provider as StoreProvider, useSelector } from 'react-redux';
import { StoreState } from './store/State';
import { ServiceResolver } from './ServiceResolver';
import { ServiceContext } from './ServiceContext';
import { JobDetail } from './screens/JobDetailScreen';
import { useServices } from './hooks/useService';
import { JobsNewScreen } from './screens/JobsNewScreen';
import { usePrevious } from './hooks/usePrevious';

const Navbar: React.FC<{}> = () => {
  const location = useLocation();
  const pathname = location.pathname;

  return (
    <Menu
      color="teal"
      inverted
      borderless
      style={{
        borderRadius: 0,
        boxShadow: 'none',
        marginBottom: 20
      }}
    >
      <Container>
        <Menu.Item name="HQ" as={Link} to="/" header />
        <Menu.Item name="Jobs" as={Link} to="/jobs" active={/\/jobs/.test(pathname)} />
        <Menu.Item name="Stats" as={Link} to="/stats" active={/\/stats/.test(pathname)} />
        <Menu.Menu position="right">
          <Menu.Item as={Link} to="/jobs/new" icon={<Icon name="plus" size="large" />} />
          <Menu.Item
            as="a"
            target="_blank"
            href="https://github.com/kohkimakimoto/hq"
            icon={<Icon name="github" size="large" />}
          />
        </Menu.Menu>
      </Container>
    </Menu>
  );
};

const MessageArea: React.FC = () => {
  const location = useLocation();
  const prevLocation: any = usePrevious(location);
  const { dispatcher } = useServices();
  const error = useSelector<StoreState, string>(state => state.error);

  useEffect(() => {
    // if the url is changed, it clear error message.
    if (prevLocation && location.key !== prevLocation.key && error !== '') {
      dispatcher.commit({ error: '' });
    }
  });

  if (error == '') {
    return null;
  }

  const handleDismiss = () => {
    dispatcher.commit({ error: '' });
  };

  return (
    <Container textAlign="center" style={{ marginBottom: 20 }}>
      <Message color="red" onDismiss={handleDismiss}>
        {error}
      </Message>
    </Container>
  );
};

const Footer: React.FC = () => {
  const version = useSelector<StoreState, string>(state => state.version);

  return (
    <Container textAlign="center" style={{ marginTop: 40, marginBottom: 40 }}>
      <Divider />
      <List horizontal divided size="small">
        <List.Item>HQ Web UI version {version}</List.Item>
      </List>
    </Container>
  );
};

const Main: React.FC = () => {
  const basename = useSelector<StoreState, string>(state => state.basename);

  return (
    <Router basename={basename}>
      <Navbar />
      <MessageArea />
      <Switch>
        <Route exact path="/">
          <Redirect to="/jobs" />
        </Route>
        <Route exact path="/jobs">
          <JobsScreen />
        </Route>
        <Route exact path="/jobs/new">
          <JobsNewScreen />
        </Route>
        <Route exact path="/jobs/:id">
          <JobDetail />
        </Route>
        <Route exact path="/stats">
          <StatsScreen />
        </Route>
        <Route path="*">
          <NotFoundScreen />
        </Route>
      </Switch>
      <Footer />
    </Router>
  );
};

export const App: React.FC<{ resolver: ServiceResolver }> = ({ resolver }) => {
  return (
    <ServiceContext.Provider value={resolver}>
      <StoreProvider store={resolver.store}>
        <Main />
      </StoreProvider>
    </ServiceContext.Provider>
  );
};
