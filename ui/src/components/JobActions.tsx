import {
  IconButton,
  HStack,
  Text,
  Button,
  useDisclosure,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalCloseButton,
  ModalBody,
  ModalFooter,
} from '@chakra-ui/react';
import React from 'react';
import { BsFillTrashFill, BsFillStopCircleFill } from 'react-icons/bs';
import { MdRestartAlt } from 'react-icons/md';

import { useDeleteJob, useRestartJob, useStopJob } from '@/api/api';
import { Job } from '@/models/Job';
import { useStatusColors } from '@/models/Status';

type Props = {
  job: Job;
  onStop: (job: Job) => void;
  onRestart: (job: Job) => void;
  onRestartAsNew: (job: Job) => void;
  onDelete: (job: Job) => void;
};

export const JobActions = ({ job, onStop, onRestart, onRestartAsNew, onDelete }: Props) => {
  const stopModalDisclosure = useDisclosure();
  const restartModalDisclosure = useDisclosure();
  const deleteModalDisclosure = useDisclosure();
  const statusColors = useStatusColors();

  const [{ loading: deleting }, executeDelete] = useDeleteJob();
  const [{ loading: stopping }, executeStop] = useStopJob();
  const [{ loading: restarting }, executeRestart] = useRestartJob();

  const handleRestart = async (job: Job) => {
    await executeRestart(job.id, false);
    onRestart(job);
    restartModalDisclosure.onClose();
  };

  const handleRestartAsNew = async (job: Job) => {
    await executeRestart(job.id, true);
    onRestartAsNew(job);
    restartModalDisclosure.onClose();
  };

  const handleDelete = async (job: Job) => {
    await executeDelete(job.id);
    onDelete(job);
    deleteModalDisclosure.onClose();
  };

  const handleStop = async (job: Job) => {
    await executeStop(job.id);
    onStop(job);
    stopModalDisclosure.onClose();
  };

  return (
    <>
      {(() => {
        if (job.status == 'running' || job.status == 'waiting') {
          return (
            <HStack spacing={2}>
              <IconButton
                onClick={stopModalDisclosure.onOpen}
                size="xs"
                aria-label="Stop"
                colorScheme="orange"
                variant="outline"
                icon={<BsFillStopCircleFill />}
              />
            </HStack>
          );
        } else {
          return (
            <HStack spacing={2}>
              <IconButton
                onClick={restartModalDisclosure.onOpen}
                size="xs"
                aria-label="Restart"
                colorScheme="teal"
                variant="outline"
                icon={<MdRestartAlt />}
              />
              <IconButton
                onClick={deleteModalDisclosure.onOpen}
                size="xs"
                aria-label="Delete"
                colorScheme="red"
                variant="outline"
                icon={<BsFillTrashFill />}
              />
            </HStack>
          );
        }
      })()}
      <Modal isOpen={stopModalDisclosure.isOpen} onClose={stopModalDisclosure.onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Stopping Job</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Text>Are you sure you want to stop the following job?</Text>
            <Text color={statusColors[job.status]} fontWeight="bold">
              {job.id}
            </Text>
          </ModalBody>
          <ModalFooter>
            <Button isLoading={stopping} size="sm" colorScheme="orange" onClick={() => handleStop(job)}>
              Stop
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>

      <Modal isOpen={restartModalDisclosure.isOpen} onClose={restartModalDisclosure.onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Restarting Job</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Text>Are you sure you want to restart the following job?</Text>
            <Text color={statusColors[job.status]} fontWeight="bold">
              {job.id}
            </Text>
          </ModalBody>
          <ModalFooter>
            <Button isLoading={restarting} size="sm" colorScheme="teal" mr={4} onClick={() => handleRestartAsNew(job)}>
              Restart as a new job
            </Button>
            <Button isLoading={restarting} size="sm" colorScheme="teal" onClick={() => handleRestart(job)}>
              Restart
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>

      <Modal isOpen={deleteModalDisclosure.isOpen} onClose={deleteModalDisclosure.onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Deleting Job</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Text>Are you sure you want to delete the following job?</Text>
            <Text color={statusColors[job.status]} fontWeight="bold">
              {job.id}
            </Text>
          </ModalBody>
          <ModalFooter>
            <Button isLoading={deleting} size="sm" colorScheme="red" onClick={() => handleDelete(job)}>
              Delete
            </Button>
          </ModalFooter>
        </ModalContent>
      </Modal>
    </>
  );
};
