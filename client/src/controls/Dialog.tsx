import React from 'react'
import { WebSocketContext } from '../WebSocket';
import { useDispatch, useSelector, shallowEqual } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { ControlsList } from './ControlsList'
import { Dialog, DialogFooter, DialogType, IDialogProps } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels, isTrue } from './Utils'

export const MyDialog = React.memo<IControlProps>(({control, parentDisabled}) => {

    const ws = React.useContext(WebSocketContext);
    const dispatch = useDispatch();

    let disabled = isTrue(control.disabled) || parentDisabled;
  
    const handleDismiss = (ev?: React.MouseEvent<HTMLButtonElement>) => {
  
        const autoDismiss = !control.autodismiss || isTrue(control.autodismiss);

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

    const cleanupLayers = () => {
        const layers = document.body.getElementsByClassName("ms-Layer--fixed")
        for (let i = 0; i < layers.length; i++) {
            let layer: Element = layers[i];
            if (!layer.hasChildNodes()) {
                document.body.removeChild(layer);
            }
        }
    }

    let dialogType: DialogType = DialogType.normal;
    switch (control.type ? control.type.toLowerCase() : '') {
      case 'largeheader': dialogType = DialogType.largeHeader; break;
      case 'close': dialogType = DialogType.close; break;
    }

    // dialog props
    const props: IDialogProps = {
        hidden: control.open !== 'true',
        minWidth: control.width ? defaultPixels(control.width) : undefined,
        maxWidth: control.maxwidth ? defaultPixels(control.maxwidth) : undefined,
        modalProps: {
            layerProps: {
                onLayerWillUnmount: () => cleanupLayers()
            },
            topOffsetFixed: isTrue(control.fixedtop),
            isBlocking: isTrue(control.blocking),
        },
        dialogContentProps: {
            type: dialogType,
            title: control.title ? control.title : undefined,
            subText: control.subtext ? control.subtext : undefined,          
        },
        styles: {
            main: {
                height: control.height !== undefined ? defaultPixels(control.height) : undefined,
            }
        },
    };

    const bodyControls = useSelector((state: any) =>
        (control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId]))
        .filter((oc: any) => oc.t !== 'footer'), shallowEqual);

    const footerControls = useSelector((state: any) =>
        (control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId]))
        .filter((oc: any) => oc.t === 'footer')
        .map((footer: any) => footer.children !== undefined ? footer.children : footer.c.map((childId: any) => state.page.controls[childId]))
        .reduce((acc: any, footerControls: any) => ([...acc, ...footerControls])), shallowEqual);

    let key = 0;

    return <Dialog {...props} onDismiss={handleDismiss}>
        <ControlsList controls={bodyControls} parentDisabled={disabled} />
        { footerControls.length > 0 ? <DialogFooter>
            {
                footerControls.map((c:any) => <ControlsList key={key++} controls={[c]} parentDisabled={disabled} />)
            }
        </DialogFooter> : "" }
    </Dialog>
})