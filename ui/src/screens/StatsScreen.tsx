import React, { useEffect } from 'react';
import { Container, Header } from 'semantic-ui-react';
import { useServices } from '../services';

export const StatsScreen: React.FC = () => {
  useEffect(() => {
    document.title = 'HQ | Stats';
  });

  const { api, errorHandler } = useServices();

  useEffect(() => {
    (async () => {
      try {
        const stat = await api.getStats();
      } catch (err) {
        errorHandler.handle(err);
      }
    })();
  });

  return (
    <Container>
      <Header as="h1">Stats</Header>
    </Container>
  );
};
