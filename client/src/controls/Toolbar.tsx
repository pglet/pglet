import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { shallowEqual, useSelector } from 'react-redux'
import { CommandBar, ICommandBarProps } from '@fluentui/react';
import { IControlProps } from './IControlProps'
import { getMenuProps } from './MenuItem'

export const Toolbar = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Button: ${control.i}`);

  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const barItems = useSelector<any, any>((state: any) =>
    getMenuProps(state, control, disabled, ws), shallowEqual)

  const overflowItems = useSelector<any, any>((state: any) => {
    const overflowControls = control.c.map((childId: any) =>
        state.page.controls[childId]).filter((ic: any) => ic.t === 'overflow' && ic.visible !== "false");
    if (overflowControls.length === 0) {
        return null
    }

    return getMenuProps(state, overflowControls[0], disabled, ws)
  }, shallowEqual)

  const farItems = useSelector<any, any>((state: any) => {
    const farControls = control.c.map((childId: any) =>
        state.page.controls[childId]).filter((ic: any) => ic.t === 'far' && ic.visible !== "false");
    if (farControls.length === 0) {
        return null
    }
    return getMenuProps(state, farControls[0], disabled, ws)
  }, shallowEqual)  

  let toolbarProps: ICommandBarProps = {
    items: barItems != null ? barItems.items : [],
    overflowItems: overflowItems != null ? overflowItems.items : [],
    farItems: farItems != null ? farItems.items : [],
    styles: {
      root: {
        paddingLeft: 0,
        paddingRight: 0,
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined   
      }
    }
  };

  return <CommandBar {...toolbarProps} />;
})