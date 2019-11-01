import React from 'react';
import { BrowserRouter as Router, Switch, Route, Link, Redirect } from 'react-router-dom';
import { Container, Divider, Dropdown, Grid, Header, Image, List, Menu, Segment } from 'semantic-ui-react';
import { JobsScreen } from './screens/JobsScreen';
import { NotFoundScreen } from './screens/NotFoundScreen';

interface PropsInterface {
  basename?: string;
}

export const App: React.FC<PropsInterface> = props => {
  return (
    <Router basename={props.basename}>
      <React.Fragment>
        <Menu color='teal' inverted borderless style={{
          borderRadius: 0,
          boxShadow: 'none',
        }}>
          <Container>
            <Menu.Item as={Link} to="/" header>
              HQ
            </Menu.Item>
            <Menu.Item as={Link} to="/jobs">
              Jobs
            </Menu.Item>
            <Menu.Item as={Link} to="/stats">
              Stats
            </Menu.Item>
          </Container>
        </Menu>
        <Switch>
          <Route path="/jobs">
            <JobsScreen/>
          </Route>
          <Route path="/stats">
            <JobsScreen/>
          </Route>
          <Route path="/">
            <Redirect to="/jobs"/>
          </Route>
          <Route path="*">
            <NotFoundScreen/>
          </Route>
        </Switch>
        <Divider/>
      </React.Fragment>
    </Router>
  );
};

App.defaultProps = {
  basename: '/ui'
};
