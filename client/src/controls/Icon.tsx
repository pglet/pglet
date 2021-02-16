import React from 'react'
import { FontIcon, mergeStyles, useTheme } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { getThemeColor, defaultPixels } from './Utils'

export const Icon = React.memo<IControlProps>(({control}) => {

  // https://developer.microsoft.com/en-us/fluentui#/controls/web/icon
  // https://developer.microsoft.com/en-us/fluentui#/styles/web/icons#fabric-react

  const theme = useTheme();

  const iconClass = mergeStyles({
    color: control.color ? getThemeColor(theme, control.color) : undefined,
    fontSize: control.size ? defaultPixels(control.size) : undefined,
    height: control.size ? defaultPixels(control.size) : undefined,
    width: control.size ? defaultPixels(control.size) : undefined,
  });

  return <FontIcon iconName={control.name} className={iconClass} />;
})