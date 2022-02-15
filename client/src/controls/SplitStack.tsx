import React from 'react'
import { shallowEqual, useSelector } from 'react-redux'
import { ControlsList } from './ControlsList'
import { WebSocketContext } from '../WebSocket';
import Split from 'split.js'
import { IControlProps } from './Control.types'
import { getThemeColor, isTrue, parseNumber } from './Utils'
import { mergeStyles, useTheme } from '@fluentui/react';

export const SplitStack = React.memo<IControlProps>(({ control, parentDisabled }) => {

    //console.log("Render stack", control.i);

    let disabled = isTrue(control.disabled) || parentDisabled;
    const ws = React.useContext(WebSocketContext);
    const theme = useTheme();

    const childControls = useSelector((state: any) => {
        return control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId])
    }, shallowEqual);

    const isHorizontal = isTrue(control.horizontal);
    console.log("isVertical", isHorizontal)
    console.log("child controls:", childControls)

    const splitRef = React.useRef<HTMLDivElement>(null);

    React.useEffect(() => {

        console.log("div", splitRef.current)
        if (splitRef.current && childControls.length === splitRef.current.children.length) {

            // go through child elements and calculate sizes
            let elems: any[] = []
            let sizes: number[] = []
            let minSizes: number[] = []
            let maxSizes: number[] = []
            let totalSize = 100;
            let unsetSizes = 0;

            const containerSize = isHorizontal ? splitRef.current.offsetWidth : splitRef.current.offsetHeight;
            console.log("containerSize", containerSize);

            for (let i = 0; i < splitRef.current.children.length; i++) {
                const elem: any = splitRef.current.children[i]
                elems.push(elem)

                const childControl = childControls[i]
                let sizeStr = isHorizontal ? childControl.width : childControl.height
                let minSizeStr = isHorizontal ? childControl.minwidth : childControl.minheight
                let maxSizeStr = isHorizontal ? childControl.maxwidth : childControl.maxheight

                // size
                if (sizeStr === undefined || sizeStr === "") {
                    sizes.push(Infinity)
                    unsetSizes++
                } else if (sizeStr.indexOf('%') != -1) {
                    sizes.push(parseNumber(sizeStr.trim().slice(0, -1)))
                } else {
                    sizes.push(parseNumber(sizeStr) / containerSize * 100)
                }

                if (sizes[i] !== Infinity) {
                    totalSize -= sizes[i]
                }

                // min size
                minSizes.push(minSizeStr ? parseNumber(minSizeStr) : 0)

                // max size
                maxSizes.push(maxSizeStr ? parseNumber(maxSizeStr) : Infinity)
            }

            // calculate the size of the rest of controls
            for (let i = 0; i < sizes.length; i++) {
                if (sizes[i] === Infinity) {
                    sizes[i] = totalSize / unsetSizes;
                }
            }

            console.log("elems", elems)
            console.log("sizes", sizes)
            console.log("minSizes", minSizes)
            console.log("maxSizes", maxSizes)

            if (elems.length > 0) {
                Split(elems, {
                    sizes: sizes,
                    minSize: minSizes,
                    maxSize: maxSizes,
                    gutterSize: control.guttersize ? parseNumber(control.guttersize) : 4,
                    direction: isHorizontal ? "horizontal" : "vertical",
                    onDragEnd: (sizes) => {
                        console.log("sizes:", sizes)
                    }
                })
            }
        }

        // eslint-disable-next-line
    }, []);

    const className = mergeStyles({
        display: isHorizontal ? "flex" : undefined,
        flexDirection: isHorizontal ? "row" : undefined,
        height: control.height ? control.height : undefined,
        width: control.width ? control.width : undefined,
        ".gutter": {
            backgroundColor: control.guttercolor ? getThemeColor(theme, control.guttercolor) : undefined,
            backgroundRepeat: "no-repeat",
            backgroundPosition: "50%"
        },
        ".gutter:hover": {
            backgroundColor: control.gutterhovercolor ? getThemeColor(theme, control.gutterhovercolor) : getThemeColor(theme, "themeLighter"),
            transitionDelay: "0.3s"
        },
        ".gutter:active": {
            backgroundColor: control.gutterdragcolor ? getThemeColor(theme, control.gutterdragcolor) : getThemeColor(theme, "themeTertiary"),
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