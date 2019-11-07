import { useEffect } from 'react';
import { useServices } from './useService';

export function useUpdateStats() {
  const { updateStats, handleError } = useServices();

  useEffect(() => {
    updateStats().catch(err => handleError(err));
  }, []);
}
