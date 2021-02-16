import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { shallowEqual, useSelector } from 'react-redux'
import {
  useTheme,
  PrimaryButton,
  DefaultButton,
  CompoundButton,
  CommandBarButton,
  IconButton,
  ActionButton,
  IButtonProps,
  IContextualMenuProps } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { getThemeColor, defaultPixels } from './Utils'
import { getMenuProps } from './MenuItem'

export const Button = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Button: ${control.i}`);

  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const theme = useTheme();

  let ButtonType = DefaultButton;
  if (control.compound === 'true') {
    ButtonType = CompoundButton
  } else if (control.toolbar === 'true') {
    ButtonType = CommandBarButton
  } else if (control.primary === 'true') {
    ButtonType = PrimaryButton
  } else if (control.action === 'true') {
    ButtonType = ActionButton
  } else if (control.icon && control.text === undefined) {
    ButtonType = IconButton
  }

  let height = control.height !== undefined ? control.height : undefined;
  if (control.toolbar === 'true' && control.height === undefined) {
    height = 40;
  }

  const menuProps = useSelector<any, IContextualMenuProps | undefined>((state: any) =>
    getMenuProps(state, control, disabled, ws, theme), shallowEqual)

  let buttonProps: Partial<IButtonProps> = {
    text: control.text ? control.text : control.i,
    href: control.url ? control.url : undefined,
    title: control.title ? control.title : undefined,
    target: control.newwindow === 'true' ? '_blank' : undefined,
    secondaryText: control.secondarytext ? control.secondarytext : undefined,
    disabled: disabled,
    primary: control.compound === 'true' && control.primary === 'true' ? true : undefined,
    split: control.split === 'true' ? true : undefined,
    menuProps: menuProps,
    styles: {    
      root: {
        width: control.width !== undefined ? defaultPixels(control.width) : undefined,
        height: defaultPixels(height),
        padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined
      },    
    }
  };

  // https://stackoverflow.com/questions/62532550/how-can-i-change-the-hover-style-of-a-primarybutton-in-fluent-ui

  if (control.icon) {
    buttonProps.iconProps = {
      iconName: control.icon
    }

    if (control.iconcolor !== undefined) {

      const iconColor = getThemeColor(theme, control.iconcolor);

      buttonProps.styles!.icon = {
        color: iconColor
      }
      buttonProps.styles!.rootHovered = {
        '.ms-Button-icon': {
          color: iconColor
        }
      };
      buttonProps.styles!.rootPressed = {
        '.ms-Button-icon': {
          color: iconColor
        }
      }
    }
  }

  const handleClick = () => {
    ws.pageEventFromWeb(control.i, 'click', control.data)
  }

  return <ButtonType onClick={handleClick} {...buttonProps} />;
})