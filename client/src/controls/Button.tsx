import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { PrimaryButton, DefaultButton, IButtonProps } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const Button = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Button: ${control.i}`);

  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  let ButtonType = (control.primary) ? PrimaryButton : DefaultButton;

  let buttonProps: IButtonProps = {
    text: control.text ? control.text : control.i,
    disabled: disabled,
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
    ws.pageEventFromWeb(control.i, 'click', control.data)
  }

  return <ButtonType onClick={handleClick} {...buttonProps} />;
})