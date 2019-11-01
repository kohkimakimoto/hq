import React from 'react';
import {Container} from "semantic-ui-react";

export const NotFoundScreen: React.FC = () => {
  return (
    <Container>
      <h1>404 Page not found</h1>
      <p>HQ Web UI returned an error.</p>
    </Container>
  );
};
