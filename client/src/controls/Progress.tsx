import React from 'react';
import { ProgressIndicator, IProgressIndicatorProps, IProgressIndicatorStyles } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const Progress = React.memo<IControlProps>(({control}) => {

  console.log(`render Progress: ${control.i}`);

  const progressProps: IProgressIndicatorProps = {
    percentComplete: control.value ? parseInt(control.value) / 100 : undefined,
    label: control.label ? control.label : null,
    description: control.description ? control.description : null
  };

  const progressStyles: Partial<IProgressIndicatorStyles> = {
    root: {
      width: control.width ? control.width : null
    },
  };

  return <ProgressIndicator {...progressProps} styles={progressStyles} />;
})