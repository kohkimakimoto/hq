import { AxiosError } from 'axios';

export type ValidationError = {
  hasError: boolean;
  message: string;
};

export const checkValidationError = (error: AxiosError | null): ValidationError => {
  if (!error || !error.response) {
    return { hasError: false, message: '' };
  }

  const { data, status } = error.response;
  if (status === 422 && data.error !== 'undefined') {
    return { hasError: true, message: data.error };
  }

  return { hasError: false, message: '' };
};
