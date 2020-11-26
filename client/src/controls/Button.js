import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { PrimaryButton, DefaultButton } from 'office-ui-fabric-react';

const button = React.memo(({control}) => {

  console.log(`render Button: ${control.i}`);

  const ws = useContext(WebSocketContext);

  let ButtonType = (control.primary) ? PrimaryButton : DefaultButton;

  let buttonProps = {
    text: control.text ? control.text : control.i
  };

  const handleClick = e => {
    ws.pageEventFromWeb(control.i, 'click', control.event)
  }

  return <ButtonType onClick={handleClick} {...buttonProps} />;
})

export default button