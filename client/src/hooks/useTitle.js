import { useEffect } from 'react';

const useTitle = function useTitle(title) {
    useEffect(() => {
      const prevTitle = document.title
      document.title = title
      return () => {
        document.title = prevTitle
      }
    })
  }

  export default useTitle