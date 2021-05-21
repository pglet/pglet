import React from 'react'
import { IControlProps } from './Control.types'

export const Html = React.memo<IControlProps>(({control}) => {

  const content = control.value ? control.value : "";
  return <div dangerouslySetInnerHTML={{ __html: content }} />;
})