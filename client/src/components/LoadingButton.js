import React from 'react';

import Button from 'react-bootstrap/Button';

class LoadingButton extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            isLoading: false
        };
      }

    simulateNetworkRequest() {
        return new Promise((resolve) => setTimeout(resolve, 2000));
    }

    handleClick = () => {
        this.setState({isLoading: true});
        this.simulateNetworkRequest().then(() => {
            this.setState({isLoading: false});
          });
    }

    render() {
        return (
            <Button
              variant="primary"
              disabled={this.state.isLoading}
              onClick={!this.state.isLoading ? this.handleClick : null}
            >
              {this.state.isLoading ? 'Loading…' : 'Click to load'}
            </Button>
          );
    }
}

export default LoadingButton;