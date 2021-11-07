// Workaround for react-router v6
// see https://github.com/pbeshai/use-query-params/issues/108#issuecomment-785209454
//     https://github.com/pbeshai/use-query-params/issues/196#issuecomment-996893750
import { Location } from 'history';
import React from 'react';
import { useLocation, useNavigate, Location as RouterLocation } from 'react-router-dom';

export const RouteAdapter: React.FunctionComponent<{
  children: React.FunctionComponent<{
    history: {
      replace(location: Location): void;
      push(location: Location): void;
    };
    location: RouterLocation;
  }>;
}> = ({ children }) => {
  const navigate = useNavigate();
  const routerLocation = useLocation();

  const adaptedHistory = React.useMemo(
    () => ({
      replace(location: Location) {
        navigate(location, { replace: true, state: location.state });
      },
      push(location: Location) {
        navigate(location, { replace: false, state: location.state });
      },
    }),
    [navigate],
  );
  if (!children) {
    return null;
  }
  return children({ history: adaptedHistory, location: routerLocation });
};
