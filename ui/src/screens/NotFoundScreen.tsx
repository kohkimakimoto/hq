import React, { useEffect } from 'react';
import { Breadcrumb, Container, Header } from 'semantic-ui-react';
import { useDocumentTitle } from '../hooks/useDocumentTitle';
import { Link } from 'react-router-dom';

export const NotFoundScreen: React.FC = () => {
  useDocumentTitle('Not Found');

  return (
    <Container>
      <Breadcrumb>
        <Breadcrumb.Section active>Not Found</Breadcrumb.Section>
      </Breadcrumb>

      <div className="page-title">
        <Header as="h1" dividing>
          404 Page not found
        </Header>
      </div>
      <p>HQ returned an error.</p>
    </Container>
  );
};
