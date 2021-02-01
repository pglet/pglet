import React from 'react'
import { FontIcon, mergeStyles } from '@fluentui/react';
import { IControlProps, defaultPixels } from './IControlProps'

export const Icon = React.memo<IControlProps>(({control}) => {

  //console.log(`render Text: ${control.i}`);

  // https://developer.microsoft.com/en-us/fluentui#/controls/web/icon
  // https://developer.microsoft.com/en-us/fluentui#/styles/web/icons#fabric-react

  const iconClass = mergeStyles({
    color: control.color ? control.color : undefined,
    fontSize: control.size ? defaultPixels(control.size) : undefined,
    height: control.size ? defaultPixels(control.size) : undefined,
    width: control.size ? defaultPixels(control.size) : undefined,
  });

  return <FontIcon iconName={control.name} className={iconClass} />;
})