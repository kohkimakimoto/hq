import React, { useEffect } from 'react';
import { Container, Header } from 'semantic-ui-react';

export const JobsScreen: React.FC = () => {
  useEffect(() => {
    document.title = 'HQ | Jobs';
  });

  return (
    <Container>
      <Header as="h1">Jobs</Header>
    </Container>
  );
};
