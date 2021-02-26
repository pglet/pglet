import React from 'react'
import { mergeStyles } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels } from './Utils'

export const IFrame = React.memo<IControlProps>(({control}) => {

  const frameClass = mergeStyles({
    border: control.border ? control.border : 'none',
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