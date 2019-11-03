import React, { useEffect } from 'react';
import { Container, Header } from 'semantic-ui-react';
import {useServices} from "../services";

export const NotFoundScreen: React.FC = () => {
  useEffect(() => {
    document.title = 'HQ | Not Found';
  });

  return (
    <Container>
      <Header as="h1">404 Page not found</Header>
      <p>HQ returned an error.</p>
    </Container>
  );
};
