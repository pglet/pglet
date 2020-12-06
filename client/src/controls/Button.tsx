import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { PrimaryButton, DefaultButton } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const Button = React.memo<IControlProps>(({control}) => {

  //console.log(`render Button: ${control.i}`);

  const ws = useContext(WebSocketContext);

  let ButtonType = (control.primary) ? PrimaryButton : DefaultButton;

  let buttonProps = {
    text: control.text ? control.text : control.i,
    styles: {
      root: {
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined   
      }
    }
  };

  const handleClick = () => {
    ws.pageEventFromWeb(control.i, 'clicked', control.data)
  }

  return <ButtonType onClick={handleClick} {...buttonProps} />;
})