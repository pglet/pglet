import React from 'react';
import { Spinner, ISpinnerProps, SpinnerSize } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels } from './Utils'

export const MySpinner = React.memo<IControlProps>(({control}) => {

  let size: SpinnerSize = SpinnerSize.medium;

  switch (control.size ? control.size.toLowerCase() : '') {
    case 'xsmall': size = SpinnerSize.xSmall; break;
    case 'small': size = SpinnerSize.small; break;
    case 'medium': size = SpinnerSize.medium; break;
    case 'large': size = SpinnerSize.large; break;
  }

  const spinnerProps: ISpinnerProps = {
    label: control.label ? control.label : null,
    labelPosition: control.labelposition ? control.labelposition.toLowerCase() : null,
    size: size,
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