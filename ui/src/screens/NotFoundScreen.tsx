import React from 'react';
import {Container, Header} from "semantic-ui-react";

export const NotFoundScreen: React.FC = () => {
  return (
    <Container>
      <Header as='h1'>404 Page not found</Header>
      <p>HQ Web UI returned an error.</p>
    </Container>
  );
};
