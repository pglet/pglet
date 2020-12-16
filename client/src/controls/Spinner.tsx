import React from 'react';
import { Spinner, ISpinnerProps } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const MySpinner = React.memo<IControlProps>(({control}) => {

  //console.log(`render Progress: ${control.i}`);

  const spinnerProps: ISpinnerProps = {
    label: control.label ? control.label : null,
    labelPosition: control.labelposition ? control.labelposition : null,
    styles: {
      root: {
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined
      }
    }
  };

  return <Spinner {...spinnerProps} />;
})