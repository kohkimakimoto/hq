import { MoonIcon, SunIcon } from '@chakra-ui/icons';
import { Text, Box, Flex, Link, Container, Button, useColorModeValue, HStack, useColorMode, Icon } from '@chakra-ui/react';
import React from 'react';
import { Outlet, Link as RouterLink } from 'react-router-dom';

import { config } from '@/lib/config';

export const Layout = () => {
  const { colorMode, toggleColorMode } = useColorMode();
  const bgColor = useColorModeValue('white', 'gray.900');
  const lineColor = useColorModeValue('gray.200', 'gray.600');

  return (
    <Flex bg={bgColor} minH="100vh" direction="column">
      {/* header */}
      <Box bg="gray.700" h={12}>
        <Container maxW="container.xl" px={6} h="full">
          <Flex direction="row" alignItems="center" justify="space-between" h="full">
            <HStack spacing={4} h="full">
              <Link as={RouterLink} to="/" h="full" display="flex" alignItems="center" _focus={{ outline: 'none' }}>
                <Icon boxSize={12} viewBox="0 0 174 80" color="white">
                  <path d="M80 80H62L62 0H80L80 80ZM18 30H26V50H18L18 80H0L0 0H18V30Z" fill="currentColor" />
                  <path d="M32 30V50L41 50V30H32Z" fill="currentColor" />
                  <path d="M47 30V50H56L56 30L47 30Z" fill="currentColor" />
                  <path
                    d="M94 9.6L112 0H131V15.7217H112V63.7217H131V80H112L94 70.5391V9.6ZM174 9.6V70.5391L156 80H137V63.7217H146.714L137 52.4522L156 57.4609V15.7217H137V0H156L174 9.6Z"
                    fill="currentColor"
                  />
                </Icon>
              </Link>
            </HStack>
            <HStack spacing={2} h="full">
              <Button
                variant="ghost"
                color="white"
                _hover={{
                  backgroundColor: 'black',
                }}
                _focus={{ outline: 'none' }}
                size="sm"
                onClick={toggleColorMode}
              >
                {colorMode === 'light' ? <SunIcon /> : <MoonIcon />}
              </Button>
            </HStack>
          </Flex>
        </Container>
      </Box>
      {/* main */}
      <Container maxW="container.xl" px={6} py={4}>
        <Outlet />
      </Container>

      <Container maxW="container.xl" px={6} py={0} mt="auto">
        <Flex
          py={4}
          alignItems="center"
          justify="center"
          borderTop="1px"
          borderTopColor={lineColor}
          direction={{ base: 'column', md: 'row' }}
        >
          <Text fontSize="12px" my={1} mx={2}>
            HQ version: {config.version}
          </Text>
          <Text fontSize="12px" my={1} mx={2}>
            Copyright (c) 2019 Kohki Makimoto
          </Text>
        </Flex>
      </Container>
    </Flex>
  );
};
