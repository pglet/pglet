import React from 'react'
import { shallowEqual, useSelector } from 'react-redux'
import { ControlsList } from './ControlsList'
import { WebSocketContext } from '../WebSocket';
import { Stack, IStackProps, IStackTokens, useTheme } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { getThemeColor, defaultPixels, isTrue, getId } from './Utils'

export const MyStack = React.memo<IControlProps>(({ control, parentDisabled }) => {

    //console.log("Render stack", control.i);

    const theme = useTheme();

    let disabled = isTrue(control.disabled) || parentDisabled;

    const ws = React.useContext(WebSocketContext);

    const handleKeyPress = (event: React.KeyboardEvent<HTMLElement>) => {
        if (event.code === "Enter" && ((event.target as HTMLElement).tagName === "INPUT" ||
            (event.target as HTMLElement).tagName === "TEXTAREA")) {
            ws.pageEventFromWeb(control.i, 'submit', control.data);
            event.stopPropagation();
        }
    }

    // stack props
    const stackProps: IStackProps = {
        horizontal: isTrue(control.horizontal),
        verticalFill: isTrue(control.verticalfill),
        // horizontalAlign: control.horizontalalign ? control.horizontalalign : "start",
        // verticalAlign: control.verticalalign ? control.verticalalign : "start",
        wrap: isTrue(control.wrap),
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
                backgroundColor: control.bgcolor ? getThemeColor(theme, control.bgcolor) : undefined,
                borderStyle: control.borderstyle ? control.borderstyle : undefined,
                borderWidth: control.borderwidth ? defaultPixels(control.borderwidth) : undefined,
                borderColor: control.bordercolor ? getThemeColor(theme, control.bordercolor) : undefined,
                borderRadius: control.borderradius ? defaultPixels(control.borderradius) : undefined,
                overflowX: isTrue(control.scrollx) ? "auto" : undefined,
                overflowY: isTrue(control.scrolly) ? "auto" : undefined,
                overflow: control.borderradius ? "hidden" : undefined
            }
        },
    };

    if (control.horizontalalign) {
        stackProps.horizontalAlign = control.horizontalalign;
    }

    if (control.verticalalign) {
        stackProps.verticalAlign = control.verticalalign;
    }

    if (isTrue(control.onsubmit)) {
        stackProps.onKeyPress = handleKeyPress;
    }

    if (isTrue(control.autoscroll)) {
        const id = getId(control.i);

        stackProps.id = id;

        window.requestAnimationFrame(() => {
            //console.log("window.requestAnimationFrame()")
            const div = document.getElementById(id);
            if (div != null) {
                div.scrollTop = div.scrollHeight - div.clientHeight;
            }
        });
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