import React from 'react';
import { Spinner, ISpinnerProps } from '@fluentui/react';
import { IControlProps, defaultPixels } from './IControlProps'

export const MySpinner = React.memo<IControlProps>(({control}) => {

  //console.log(`render Progress: ${control.i}`);

  const spinnerProps: ISpinnerProps = {
    label: control.label ? control.label : null,
    labelPosition: control.labelposition ? control.labelposition : null,
    styles: {
      root: {
        width: control.width !== undefined ? defaultPixels(control.width) : undefined,
        height: control.height !== undefined ? defaultPixels(control.height) : undefined,
        padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined
      }
    }
  };

  return <Spinner {...spinnerProps} />;
})