import React from 'react';
import { Spinner, ISpinnerProps } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels } from './Utils'

export const MySpinner = React.memo<IControlProps>(({control}) => {

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