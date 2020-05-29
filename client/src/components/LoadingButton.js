import React, { useState } from 'react';

import Button from 'react-bootstrap/Button';

const LoadingButton = () => {

    const [isLoading, setIsLoading] = useState(false);
    const simulateNetworkRequest = () => new Promise((resolve) => setTimeout(resolve, 2000));

    const handleClick = () => {
      setIsLoading(true);
      simulateNetworkRequest().then(() => {
        setIsLoading(false);
      })
    }

    return (
      <Button
        variant="primary"
        disabled={isLoading}
        onClick={!isLoading ? handleClick : null}
      >
        {isLoading ? 'Loading…' : 'Click to load'}
      </Button>
    );
}

export default LoadingButton;