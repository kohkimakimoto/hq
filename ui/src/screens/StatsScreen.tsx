import React, { useEffect } from 'react';
import { Container, Header } from 'semantic-ui-react';

export const StatsScreen: React.FC = () => {
  useEffect(() => {
    document.title = 'HQ | Stats';
  });

  return (
    <Container>
      <Header as="h1">Stats</Header>
    </Container>
  );
};
