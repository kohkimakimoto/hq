import React from 'react';
import { Helmet } from 'react-helmet-async';

export const NotFoundPage = () => {
  return (
    <>
      <Helmet>
        <title>Not Found | HQ</title>
      </Helmet>
      <div>
        <h1>404 Page not found</h1>
      </div>
    </>
  );
};
