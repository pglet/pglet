import React, { useContext } from 'react'
import { WebSocketContext } from '../WebSocket';
import { useDispatch, useSelector, shallowEqual } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { ControlsList } from './ControlsList'
import { Panel, IPanelProps, PanelType } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels } from './Utils'

export const MyPanel = React.memo<IControlProps>(({ control, parentDisabled }) => {

    let disabled = (control.disabled === 'true') || parentDisabled;

    //console.log(`render dialog: ${control.i}`);

    const ws = useContext(WebSocketContext);
    const dispatch = useDispatch();

    const cleanupLayers = () => {
        const layers = document.body.getElementsByClassName("ms-Layer--fixed")
        for (let i = 0; i < layers.length; i++) {
            let layer: Element = layers[i];
            if (!layer.hasChildNodes()) {
                document.body.removeChild(layer);
            }
        }
    }

    const handleDismiss = (ev?: React.SyntheticEvent<HTMLElement> | KeyboardEvent) => {

        const autoDismiss = !control.autodismiss || control.autodismiss === 'true';

        if (autoDismiss) {
            const val = "false"

            let payload: any = {}
            if (control.f) {
                // binding redirect
                const p = control.f.split('|')
                payload["i"] = p[0]
                payload[p[1]] = val
            } else {
                // unbound control
                payload["i"] = control.i
                payload["open"] = val
            }
    
            dispatch(changeProps([payload]));
            ws.updateControlProps([payload]);
        }

        ws.pageEventFromWeb(control.i, 'dismiss', control.data)

        if (!autoDismiss) {
            ev?.preventDefault();
            return
        }
    }

    // dialog props
    const props: IPanelProps = {
        isOpen: control.open === 'true',
        isLightDismiss: control.lightdismiss === 'true',
        isBlocking: control.blocking !== 'false',
        headerText: control.title ? control.title : undefined,
        layerProps: {
            onLayerWillUnmount: () => cleanupLayers()
        },   
    };

    switch (control.type ? control.type.toLowerCase() : '') {
        case 'small': props.type = PanelType.smallFixedFar; break;
        case 'smallleft': props.type = PanelType.smallFixedNear; break;
        case 'medium': props.type = PanelType.medium; break;
        case 'large': props.type = PanelType.large; break;
        case 'largefixed': props.type = PanelType.largeFixed; break;
        case 'extralarge': props.type = PanelType.extraLarge; break;
        case 'fluid': props.type = PanelType.smallFluid; break;
        case 'custom': props.type = PanelType.custom; break;
        case 'customleft': props.type = PanelType.customNear; break;
        default: props.type = PanelType.smallFixedFar; break;
    }

    if (props.type === PanelType.custom || props.type === PanelType.customNear) {
        props.customWidth = control.width ? defaultPixels(control.width) : undefined
    }

    const bodyControls = useSelector((state: any) =>
        (control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId]))
            .filter((oc: any) => oc.t !== 'footer'), shallowEqual);

    const footerControls = useSelector((state: any) =>
        (control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId]))
            .filter((oc: any) => oc.t === 'footer')
            .map((footer: any) => footer.children !== undefined ? footer.children : footer.c.map((childId: any) => state.page.controls[childId]))
            .reduce((acc: any, footerControls: any) => ([...acc, ...footerControls])), shallowEqual);

    const onRenderFooterContent = React.useCallback(
        () => (<ControlsList controls={footerControls} parentDisabled={disabled} />),
        [footerControls, disabled]);

    if (footerControls.length > 0) {
        props.onRenderFooterContent = onRenderFooterContent
        props.isFooterAtBottom = true;
    }

    return <Panel {...props} onDismiss={handleDismiss}>
        <ControlsList controls={bodyControls} parentDisabled={disabled} />
    </Panel>
})