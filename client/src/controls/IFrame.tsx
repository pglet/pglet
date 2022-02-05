import React from 'react'
import { mergeStyles, useTheme } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels, getThemeColor } from './Utils'

export const IFrame = React.memo<IControlProps>(({ control }) => {

  const theme = useTheme();

  const frameClass = mergeStyles({
    border: control.border ? control.border : 'none',
    borderWidth: control.borderwidth ? defaultPixels(control.borderwidth) : undefined,
    borderColor: control.bordercolor ? getThemeColor(theme, control.bordercolor) : undefined,
    borderStyle: control.borderstyle ? control.borderstyle : undefined,
    borderRadius: control.borderradius ? defaultPixels(control.borderradius) : undefined
  });

  const title = control.title ? control.title : control.i;

  const props: React.DetailedHTMLProps<React.IframeHTMLAttributes<HTMLIFrameElement>, HTMLIFrameElement> = {
    src: control.src ? control.src : undefined,
    title: control.title ? control.title : control.i,
    width: control.width !== undefined ? defaultPixels(control.width) : undefined,
    height: control.height !== undefined ? defaultPixels(control.height) : undefined,
  };

  return <iframe title={title} {...props} className={frameClass} />;
})