import React from 'react';
import { ProgressIndicator, IProgressIndicatorProps } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const Progress = React.memo<IControlProps>(({control}) => {

  //console.log(`render Progress: ${control.i}`);

  const progressProps: IProgressIndicatorProps = {
    percentComplete: control.value ? parseInt(control.value) / 100 : undefined,
    label: control.label ? control.label : null,
    description: control.description ? control.description : null,
    styles: {
      root: {
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined
      }
    }
  };

  return <ProgressIndicator {...progressProps} />;
})