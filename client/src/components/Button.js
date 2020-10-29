import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';

const Button = React.memo(({control}) => {

  console.log(`render Button: ${control.i}`);

  const ws = useContext(WebSocketContext);

  const handleClick = e => {
    ws.pageEventFromWeb(control.i, 'click', control.event)
  }

  return <button type="button" className="btn btn-primary" onClick={handleClick}>{control.text}</button>;
})

export default Button