import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { PrimaryButton } from 'office-ui-fabric-react';

const button = React.memo(({control}) => {

  console.log(`render Button: ${control.i}`);

  const ws = useContext(WebSocketContext);

  const handleClick = e => {
    ws.pageEventFromWeb(control.i, 'click', control.event)
  }

  return <PrimaryButton onClick={handleClick} text={control.text} />;
})

export default button