import React, { useContext, useEffect } from 'react';
import { BrowserRouter as Router, Switch, Route, Link, Redirect, useLocation } from 'react-router-dom';
import { Container, Divider, Dropdown, Grid, Header, Image, List, Menu, Segment, Icon } from 'semantic-ui-react';
import { JobsScreen } from './screens/JobsScreen';
import { NotFoundScreen } from './screens/NotFoundScreen';
import { StatsScreen } from './screens/StatsScreen';
import {Store } from "redux";
import {Provider as StoreProvider, useSelector } from 'react-redux';
import {StoreState} from "./store/State";

const Navbar: React.FC<{}> = () => {
  const location = useLocation();

  return (
    <Menu
      color="teal"
      inverted
      borderless
      style={{
        borderRadius: 0,
        boxShadow: 'none',
        marginBottom: 40
      }}
    >
      <Container>
        <Menu.Item name="HQ" as={Link} to="/" header />
        <Menu.Item name="Jobs" as={Link} to="/jobs" active={location.pathname == '/jobs'} />
        <Menu.Item name="Stats" as={Link} to="/stats" active={location.pathname == '/stats'} />
        <Menu.Menu position="right">
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

const Footer: React.FC = () => {
  const version = useSelector<StoreState, string>(state => state.version);

  return (
    <Container textAlign="center" style={{ marginTop: 40 }}>
      <Divider />
      <List horizontal divided size="small">
        <List.Item>version {version}</List.Item>
      </List>
    </Container>
  );
};

const Main: React.FC = () => {
  const basename = useSelector<StoreState, string>(state => state.basename);

  return (
    <Router basename={basename}>
      <Navbar />
      <Switch>
        <Route exact path="/">
          <Redirect to="/jobs" />
        </Route>
        <Route path="/jobs">
          <JobsScreen />
        </Route>
        <Route path="/stats">
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

export const App: React.FC<{ store: Store<StoreState> }> = props => {
  return (
    <StoreProvider store={props.store}>
      <Main />
    </StoreProvider>
  );
};
