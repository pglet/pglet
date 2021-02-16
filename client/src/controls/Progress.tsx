import React from 'react';
import { ProgressIndicator, IProgressIndicatorProps } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels } from './Utils'

export const Progress = React.memo<IControlProps>(({control}) => {

  const progressProps: IProgressIndicatorProps = {
    percentComplete: control.value ? parseInt(control.value) / 100 : undefined,
    label: control.label ? control.label : null,
    description: control.description ? control.description : null,
    styles: {
      root: {
        width: control.width !== undefined ? defaultPixels(control.width) : undefined,
        height: control.height !== undefined ? defaultPixels(control.height) : undefined,
        padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined
      }
    }
  };

  return <ProgressIndicator {...progressProps} />;
})