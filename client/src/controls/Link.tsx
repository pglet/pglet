import React, { useContext } from 'react'
import { Link, ILinkProps } from '@fluentui/react';
import { shallowEqual, useSelector } from 'react-redux'
import { ControlsList } from './ControlsList'
import { WebSocketContext } from '../WebSocket';
import { IControlProps } from './Control.types'
import { defaultPixels } from './Utils'

export const MyLink = React.memo<IControlProps>(({ control, parentDisabled }) => {

  const ws = useContext(WebSocketContext);

  let disabled = (control.disabled === 'true') || parentDisabled;

  const handleClick = () => {
    ws.pageEventFromWeb(control.i, 'click', control.data)
  }

  const linkProps: ILinkProps = {
    href: control.url ? control.url : undefined,
    target: control.newwindow === 'true' ? '_blank' : undefined,
    title: control.title ? control.title : undefined,
    onClick: handleClick,
    disabled: disabled,
    styles: {
      root: {
        fontSize: control.size ? defaultPixels(control.size) : '14px',
        fontWeight: control.bold === 'true' ? 'bold' : undefined,
        fontStyle: control.italic === 'true' ? 'italic' : undefined,
        textAlign: control.align !== undefined ? control.align : undefined,
        width: control.width !== undefined ? defaultPixels(control.width) : undefined,
        height: control.height !== undefined ? defaultPixels(control.height) : undefined,
        padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined,
      }
    }
  };

  const childControls = useSelector((state: any) => {
    return control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId])
  }, shallowEqual);

  return <Link {...linkProps}>{childControls.length > 0 ?
    <ControlsList controls={childControls} parentDisabled={disabled} />
    : control.pre === "true" ? <pre>{control.value}</pre> : control.value}</Link>;
})