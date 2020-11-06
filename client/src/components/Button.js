import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import Button from 'react-bootstrap/Button';

const button = React.memo(({control}) => {

  console.log(`render Button: ${control.i}`);

  const ws = useContext(WebSocketContext);

  const handleClick = e => {
    ws.pageEventFromWeb(control.i, 'click', control.event)
  }

  return <Button variant="primary" onClick={handleClick}>{control.text}</Button>;
})

export default button