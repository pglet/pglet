import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import {
  PrimaryButton,
  DefaultButton,
  CompoundButton,
  CommandBarButton,
  IconButton,
  ActionButton,
  IButtonProps,
  IStyle } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const Button = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Button: ${control.i}`);

  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  let ButtonType = DefaultButton;
  if (control.compound === 'true') {
    ButtonType = CompoundButton
  } else if (control.commandbar === 'true') {
    ButtonType = CommandBarButton
  } else if (control.primary === 'true') {
    ButtonType = PrimaryButton
  } else if (control.action === 'true') {
    ButtonType = ActionButton
  } else if (control.icon && control.text === undefined) {
    ButtonType = IconButton
  }

  let height = control.height !== undefined ? control.height : undefined;
  if (control.commandbar === 'true' && control.height === undefined) {
    height = 40;
  }

  let buttonProps: Partial<IButtonProps> = {
    text: control.text ? control.text : control.i,
    href: control.url ? control.url : undefined,
    target: control.newwindow === 'true' ? '_blank' : undefined,
    secondaryText: control.secondarytext ? control.secondarytext : undefined,
    disabled: disabled,
    primary: control.compound === 'true' && control.primary === 'true' ? true : undefined,
    styles: {
      root: {
        width: control.width !== undefined ? control.width : undefined,
        height: height,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined   
      }
    }
  };

  if (control.icon) {
    buttonProps.iconProps = {
      iconName: control.icon
    }
  }

  const handleClick = () => {
    ws.pageEventFromWeb(control.i, 'click', control.data)
  }

  return <ButtonType onClick={handleClick} {...buttonProps} />;
})