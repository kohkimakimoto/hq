import { extendTheme, withDefaultSize } from '@chakra-ui/react';

const theme = extendTheme(
  {
    components: {
      FormLabel: {
        baseStyle: {
          fontWeight: 'bold',
        },
      },
    },
  },
  withDefaultSize({
    size: 'md',
    components: ['Button', 'Input'],
  }),
);

export default theme;
