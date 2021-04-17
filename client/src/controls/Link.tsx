import React from 'react'
import { Link, ILinkProps } from '@fluentui/react';
import { shallowEqual, useSelector } from 'react-redux'
import { ControlsList } from './ControlsList'
import { WebSocketContext } from '../WebSocket';
import { IControlProps } from './Control.types'
import { defaultPixels, getId, isTrue } from './Utils'

export const MyLink = React.memo<IControlProps>(({ control, parentDisabled }) => {

  const ws = React.useContext(WebSocketContext);

  let disabled = isTrue(control.disabled) || parentDisabled;

  const handleClick = () => {
    ws.pageEventFromWeb(control.i, 'click', control.data)
  }

  const linkProps: ILinkProps = {
    id: getId(control.f ? control.f : control.i),
    href: control.url ? control.url : undefined,
    target: isTrue(control.newwindow) ? '_blank' : undefined,
    title: control.title ? control.title : undefined,
    onClick: handleClick,
    disabled: disabled,
    styles: {
      root: {
        fontSize: control.size ? defaultPixels(control.size) : '14px',
        fontWeight: isTrue(control.bold) ? 'bold' : undefined,
        fontStyle: isTrue(control.italic) ? 'italic' : undefined,
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
    : isTrue(control.pre) ? <pre>{control.value}</pre> : control.value}</Link>;
})