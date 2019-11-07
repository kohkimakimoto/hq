import React, { useState, useEffect, useRef } from 'react';

// see https://overreacted.io/making-setinterval-declarative-with-react-hooks/
export function useInterval(callback: Function, delay: number) {
  const savedCallback = useRef<Function>();

  // Remember the latest callback.
  useEffect(() => {
    savedCallback.current = callback;
  });

  // Set up the interval.
  useEffect(() => {
    function tick() {
      savedCallback.current && savedCallback.current();
    }
    if (delay !== null) {
      let id = setInterval(tick, delay);
      return () => clearInterval(id);
    }
  }, [delay]);
}
