import React, { useEffect } from 'react';
import { Container, Header } from 'semantic-ui-react';
import { useEffectDocumentTitle } from '../hooks/useEffectDocumentTitle';

export const NotFoundScreen: React.FC = () => {
  useEffectDocumentTitle('Not Found');

  return (
    <Container>
      <div className="page-title">
        <Header as="h1" dividing>
          404 Page not found
        </Header>
      </div>
      <p>HQ returned an error.</p>
    </Container>
  );
};
