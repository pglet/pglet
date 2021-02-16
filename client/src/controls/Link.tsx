import React, { useContext } from 'react'
import { Link, ILinkProps } from '@fluentui/react';
import { WebSocketContext } from '../WebSocket';
import { IControlProps } from './Control.types'
import { defaultPixels } from './Utils'

export const MyLink = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Text: ${control.i}`);

  // https://developer.microsoft.com/en-us/fluentui#/controls/web/references/ifontstyles#IFontStyles

  const preFont = 'SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace';

  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const linkProps: ILinkProps = {
    href: control.url ? control.url : undefined,
    target: control.newwindow === 'true' ? '_blank' : undefined,
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
        whiteSpace: control.pre === 'true' ? 'pre' : undefined,
        fontFamily: control.pre === 'true' ? preFont : undefined,
      }
    }
  };

  const handleClick = () => {
    ws.pageEventFromWeb(control.i, 'click', control.data)
  }

  if (!linkProps.href) {
    linkProps.onClick = handleClick;
  }

  return <Link {...linkProps}>{control.value}</Link>;
})