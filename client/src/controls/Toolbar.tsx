import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { shallowEqual, useSelector } from 'react-redux'
import {
  CommandBar,
  ICommandBarProps,
  IContextualMenuProps,
  IContextualMenuItem } from '@fluentui/react';
import { IControlProps } from './IControlProps'
import { getMenuProps } from './MenuItem'

export const Toolbar = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Button: ${control.i}`);

  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const handleMenuItemClick = (ev?: React.MouseEvent<HTMLElement> | React.KeyboardEvent<HTMLElement>, item?: IContextualMenuItem) => {
    ws.pageEventFromWeb(control.i, 'click', item!.key)
  }

  const barItems = useSelector<any, any>((state: any) =>
    getMenuProps(state, control, null, handleMenuItemClick), shallowEqual)

  let buttonProps: ICommandBarProps = {
    items: barItems.items,
    styles: {
      root: {
        paddingLeft: 0,
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined   
      }
    }
  };

  return <CommandBar {...buttonProps} />;
})