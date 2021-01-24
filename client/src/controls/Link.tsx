import React from 'react'
import { Link, ILinkProps } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const MyLink = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Text: ${control.i}`);

  // https://developer.microsoft.com/en-us/fluentui#/controls/web/references/ifontstyles#IFontStyles

  const preFont = 'SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace';

  let disabled = (control.disabled === 'true') || parentDisabled;

  const linkProps: ILinkProps = {
    href: control.url ? control.url : undefined,
    target: control.newwindow === 'true' ? '_blank' : undefined,
    disabled: disabled,    
    styles: {
      root: {
        fontWeight: control.bold === 'true' ? 'bold' : undefined,
        fontStyle: control.italic === 'true' ? 'italic' : undefined,
        textAlign: control.align !== undefined ? control.align : undefined,
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined,
        whiteSpace: control.pre === 'true' ? 'pre' : undefined,
        fontFamily: control.pre === 'true' ? preFont : undefined,
      }
    }
  };

  return <Link {...linkProps}>{control.value}</Link>;
})