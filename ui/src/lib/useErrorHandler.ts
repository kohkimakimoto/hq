import { useToast, ToastOptions } from '@chakra-ui/react';

export type ErrorHandler = (error: Error | null | undefined) => void;

export type UseErrorHandlerOptions = {
  toastId?: string;
  duration?: ToastOptions['duration'];
};

export const useErrorHandler = (options?: UseErrorHandlerOptions): ErrorHandler => {
  const toast = useToast();

  return (error) => {
    const toastId = options?.toastId ?? 'global-error';
    if (error && !toast.isActive(toastId)) {
      toast({
        id: toastId,
        title: 'Error',
        description: error.message,
        status: 'error',
        position: 'top-right',
        duration: options?.duration ?? null,
        isClosable: true,
      });
    }
  };
};
