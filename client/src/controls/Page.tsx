import React from 'react'
import { shallowEqual, useSelector } from 'react-redux'
import { ControlsList } from './ControlsList'
import useTitle from '../hooks/useTitle'
import { Stack, IStackProps, IStackTokens } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const Page = React.memo<IControlProps>(({control}) => {

  //console.log(`render page: ${control.i}`);

  // page title
  let title = "Pglet";
  if (control.title) {
    title = control.title
  }
  useTitle(title)

  // stack props
  const stackProps: IStackProps = {
    verticalFill: control.verticalFill ? control.verticalFill : false,
    horizontalAlign: control.horizontalalign ? control.horizontalalign : "start",
    verticalAlign: control.verticalalign ? control.verticalalign : "start",
    styles: {
      root: {
        width: control.width ? control.width : "100%",
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding ? control.padding : "10px",
        margin: control.margin !== undefined ? control.margin : undefined
      }
    },
  };

  const stackTokens: IStackTokens = {
    childrenGap: control.gap ? control.gap : 10
  }

  const childControls = useSelector((state: any) => control.c.map((childId: string) => state.page.controls[childId]), shallowEqual);

  return <Stack tokens={stackTokens} {...stackProps}>
    <ControlsList controls={childControls} />
  </Stack>
})