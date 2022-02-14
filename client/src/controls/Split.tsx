import React from 'react'
import { shallowEqual, useSelector } from 'react-redux'
import { ControlsList } from './ControlsList'
import { WebSocketContext } from '../WebSocket';
import Split from 'split.js'
import { IControlProps } from './Control.types'
import { getThemeColor, isTrue } from './Utils'
import { mergeStyles, useTheme } from '@fluentui/react';

export const MySplit = React.memo<IControlProps>(({ control, parentDisabled }) => {

    //console.log("Render stack", control.i);

    let disabled = isTrue(control.disabled) || parentDisabled;
    const ws = React.useContext(WebSocketContext);
    const theme = useTheme();

    const childControls = useSelector((state: any) => {
        return control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId])
    }, shallowEqual);

    const isVertical = isTrue(control.vertical);
    console.log("isVertical", isVertical)

    const splitRef = React.useRef<HTMLDivElement>(null);

    React.useEffect(() => {

        console.log("div", splitRef.current)
        if (splitRef.current) {
            let elems: any[] = []
            for (let i = 0; i < splitRef.current.children.length; i++) {
                elems.push(splitRef.current.children[i])
            }
            console.log("elems", elems)
            if (elems.length > 0) {
                Split(elems, {
                    sizes: isVertical ? [50, 50] : [25, 50, 25],
                    gutterSize: 4,
                    direction: isVertical ? "vertical" : "horizontal",
                    onDragEnd: (sizes) => {
                        console.log("sizes:", sizes)
                    }
                })
            }
        }

        // eslint-disable-next-line
    }, []);

    const className = mergeStyles({
        display: isVertical ? undefined : "flex",
        flexDirection: isVertical ? undefined : 'row',
        height: '100%',
        ".gutter": {
            backgroundColor: control.bgcolor ? getThemeColor(theme, control.bgcolor) : undefined,
            backgroundRepeat: "no-repeat",
            backgroundPosition: "50%"
        },
        ".gutter:hover": {
            backgroundColor: control.hovercolor ? getThemeColor(theme, control.hovercolor) : getThemeColor(theme, "themeLighter"),
            transitionDelay: "0.3s"
        },
        ".gutter:active": {
            backgroundColor: control.dragcolor ? getThemeColor(theme, control.dragcolor) : getThemeColor(theme, "themeTertiary"),
            transitionDelay: "0s"
        },
        ".gutter.gutter-horizontal": {
            cursor: "col-resize"
        },
        ".gutter.gutter-vertical": {
            //backgroundImage: "url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAUAAAAeCAYAAADkftS9AAAAIklEQVQoU2M4c+bMfxAGAgYYmwGrIIiDjrELjpo5aiZeMwF+yNnOs5KSvgAAAABJRU5ErkJggg==')",
            cursor: "row-resize"
        }
    });

    return <div ref={splitRef} className={className}>
        <ControlsList controls={childControls} parentDisabled={disabled} />
    </div>
})