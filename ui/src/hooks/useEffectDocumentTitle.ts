import { useEffect } from 'react';

export function useEffectDocumentTitle(title) {
  useEffect(() => {
    document.title = 'HQ | ' + title;
  });
}
