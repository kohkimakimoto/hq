import { useContext } from 'react';
import { ServiceContext } from '../ServiceContext';
import { ServiceResolver } from '../ServiceResolver';

export const useServices = () => {
  return useContext<ServiceResolver>(ServiceContext);
};
