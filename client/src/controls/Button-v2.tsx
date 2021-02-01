import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { shallowEqual, useSelector } from 'react-redux'
import {
  // PrimaryButton,
  // DefaultButton,
  // CompoundButton,
  // CommandBarButton,
  // IconButton,
  // ActionButton,
  // IButtonProps,
  IContextualMenuProps, FontIcon } from '@fluentui/react';
import { Button, MenuButton, SplitButton, ButtonProps } from '@fluentui/react-button';  
import { IControlProps } from './IControlProps'
import { getMenuProps } from './MenuItem'

export const MyButton = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Button: ${control.i}`);

  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  let ButtonTag = Button;
  // if (control.compound === 'true') {
  //   ButtonTag = CompoundButton
  // } else if (control.toolbar === 'true') {
  //   ButtonTag = CommandBarButton
  // } else if (control.primary === 'true') {
  //   ButtonTag = PrimaryButton
  // } else if (control.action === 'true') {
  //   ButtonTag = ActionButton
  // } else if (control.icon && control.text === undefined) {
  //   ButtonTag = IconButton
  // }

  let height = control.height !== undefined ? control.height : undefined;
  if (control.toolbar === 'true' && control.height === undefined) {
    height = 40;
  }

  // const menuProps = useSelector<any, IContextualMenuProps | undefined>((state: any) =>
  //   getMenuProps(state, control, disabled, ws), shallowEqual)

  const menuProps = {
    items: [
      { key: 'a', text: 'text a'}
    ]
  }
  
  if (menuProps != null && menuProps.items.length > 0) {
    ButtonTag = control.split === 'true' ? SplitButton : MenuButton;
    console.log(menuProps!);
  }



  let buttonProps: ButtonProps = {
    content: control.text ? control.text : control.i,
    primary: control.primary === 'true',
    href: control.url ? control.url : undefined,
    title: control.title ? control.title : undefined,
    target: control.newwindow === 'true' ? '_blank' : undefined,
    secondaryText: control.secondarytext ? control.secondarytext : undefined,
    disabled: disabled,
    split: control.split === 'true' ? true : undefined,
    iconOnly: (control.icon && control.text === undefined),
    ghost: control.ghost === 'true' || (control.icon && control.text === undefined),
    menu: menuProps,
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
    buttonProps.icon = <FontIcon iconName={control.icon} />
  }

  const handleClick = () => {
    ws.pageEventFromWeb(control.i, 'click', control.data)
  }

  return <ButtonTag onClick={handleClick} {...buttonProps} />;
})