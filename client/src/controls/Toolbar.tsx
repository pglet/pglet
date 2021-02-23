import React from 'react';
import { WebSocketContext } from '../WebSocket';
import { shallowEqual, useSelector } from 'react-redux'
import { CommandBar, ICommandBarProps, useTheme } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels } from './Utils'
import { getMenuProps } from './MenuItem'

export const Toolbar = React.memo<IControlProps>(({control, parentDisabled}) => {

  const ws = React.useContext(WebSocketContext);
  const theme = useTheme();

  let disabled = (control.disabled === 'true') || parentDisabled;

  const barItems = useSelector<any, any>((state: any) =>
    getMenuProps(state, control, disabled, ws, theme), shallowEqual)

  const overflowItems = useSelector<any, any>((state: any) => {
    const overflowControls = (control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId]))
      .filter((ic: any) => ic.t === 'overflow' && ic.visible !== "false");
    if (overflowControls.length === 0) {
        return null
    }

    return getMenuProps(state, overflowControls[0], disabled, ws, theme)
  }, shallowEqual)

  const farItems = useSelector<any, any>((state: any) => {
    const farControls = (control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId]))
      .filter((ic: any) => ic.t === 'far' && ic.visible !== "false");
    if (farControls.length === 0) {
        return null
    }
    return getMenuProps(state, farControls[0], disabled, ws, theme)
  }, shallowEqual)  

  let toolbarProps: ICommandBarProps = {
    items: barItems != null ? barItems.items : [],
    overflowItems: overflowItems != null ? overflowItems.items : [],
    farItems: farItems != null ? farItems.items : [],
    styles: {
      root: {
        backgroundColor: "inherit",
        paddingLeft: 0,
        paddingRight: 0,
        width: control.width !== undefined ? defaultPixels(control.width) : undefined,
        height: control.height !== undefined ? defaultPixels(control.height) : undefined,
        padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined   
      }
    }
  };

  const darkerBkg = "linear-gradient(rgba(0,0,0,0.1),rgba(0,0,0,0.1)) !important";

  toolbarProps.overflowButtonProps = {
    styles: {
      root: {
        backgroundColor: "inherit",
        color: '#fff'
      },
      rootHovered: {
        background: darkerBkg,
      },
      rootExpanded: {
        background: darkerBkg,
      },
      rootPressed: {
        background: darkerBkg,
      },      
      menuIcon: {
        color: '#fff!important',
      },
      menuIconExpanded: {
        background: "transparent!important",
      }
    }
  }

  return <CommandBar {...toolbarProps} />;
})