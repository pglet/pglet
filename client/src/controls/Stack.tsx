import React from 'react'
import { shallowEqual, useSelector } from 'react-redux'
import { ControlsList } from './ControlsList'
import { Stack, IStackProps, IStackTokens } from '@fluentui/react';
import { IControlProps, defaultPixels } from './IControlProps'

export const MyStack = React.memo<IControlProps>(({control, parentDisabled}) => {

    //console.log(`render stack: ${control.i}`);

    let disabled = (control.disabled === 'true') || parentDisabled;

    // stack props
    const stackProps: IStackProps = {
        horizontal: control.horizontal ? control.horizontal : false,
        verticalFill: control.verticalfill ? control.verticalfill : false,
        // horizontalAlign: control.horizontalalign ? control.horizontalalign : "start",
        // verticalAlign: control.verticalalign ? control.verticalalign : "start",
        styles: {
            root: {
                width: control.width ? defaultPixels(control.width) : undefined,
                minWidth: control.minwidth ? defaultPixels(control.minwidth) : undefined,
                maxWidth: control.maxwidth ? defaultPixels(control.maxwidth) : undefined,
                height: control.height !== undefined ? defaultPixels(control.height) : undefined,
                minHeight: control.minheight !== undefined ? defaultPixels(control.minheight) : undefined,
                maxHeight: control.maxheight !== undefined ? defaultPixels(control.maxheight) : undefined,
                padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
                margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined,
                backgroundColor: control.bgcolor ? control.bgcolor : undefined,
                border: control.border ? control.border : undefined,
                borderRadius: control.borderradius ? defaultPixels(control.borderradius) : undefined,
                borderLeft: control.borderleft ? control.borderleft : undefined,
                borderRight: control.borderright ? control.borderright : undefined,
                borderTop: control.bordertop ? control.bordertop : undefined,
                borderBottom: control.borderbottom ? control.borderbottom : undefined,                
            }
        },
    };

    if (control.horizontalalign) {
        stackProps.horizontalAlign = control.horizontalalign;
    }

    if (control.verticalalign) {
        stackProps.verticalAlign = control.verticalalign;
    }    

    const stackTokens: IStackTokens = {
        childrenGap: control.gap ? control.gap : 10
    }

    const childControls = useSelector((state: any) => {
        return control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId])
    }, shallowEqual);

    return <Stack tokens={stackTokens} {...stackProps}>
        <ControlsList controls={childControls} parentDisabled={disabled} />
    </Stack>
})