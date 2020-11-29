import React from 'react'
import { shallowEqual, useSelector } from 'react-redux'
import ControlsList from './ControlsList'
import { Stack } from 'office-ui-fabric-react/lib/Stack';

const MyStack = React.memo(({ control }) => {

    console.log(`render stack: ${control.i}`);

    // stack props
    const stackProps = {
        horizontal: control.horizontal ? control.horizontal : false,
        verticalFill: true,
        horizontalAlign: control.horizontalalign ? control.horizontalalign : "start",
        verticalAlign: control.verticalalign ? control.verticalalign : "start",
        childrenGap: control.gap ? control.gap : 10,
        styles: {
            root: {
                width: control.width ? control.width : "100%"
            }
        },
    };

    const childControls = useSelector(state => control.c.map(childId => state.page.controls[childId]), shallowEqual);

    return <Stack {...stackProps}>
        <ControlsList controls={childControls} />
    </Stack>
})

export default MyStack