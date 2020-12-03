import React from 'react'
import { shallowEqual, useSelector } from 'react-redux'
import ControlsList from './ControlsList'
import { Stack, IStackProps, IStackTokens } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const MyStack = React.memo<IControlProps>(({control}) => {

    console.log(`render stack: ${control.i}`);

    // stack props
    const stackProps: IStackProps = {
        horizontal: control.horizontal ? control.horizontal : false,
        verticalFill: control.verticalFill ? control.verticalFill : false,
        horizontalAlign: control.horizontalalign ? control.horizontalalign : "start",
        verticalAlign: control.verticalalign ? control.verticalalign : "start",
        styles: {
            root: {
                width: control.width ? control.width : "100%"
            }
        },
    };

    const stackTokens: IStackTokens = {
        childrenGap: control.gap ? control.gap : 10
    }

    const childControls = useSelector((state: any) => control.c.map((childId: any) => state.page.controls[childId]), shallowEqual);

    return <Stack tokens={stackTokens} {...stackProps}>
        <ControlsList controls={childControls} />
    </Stack>
})

export default MyStack