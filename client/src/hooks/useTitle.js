import React from 'react';

const useTitle = function useTitle(title) {
    React.useEffect(() => {
      const prevTitle = document.title
      document.title = title
      return () => {
        document.title = prevTitle
      }
    })
  }

  export default useTitle