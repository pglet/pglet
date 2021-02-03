import React, { useContext, useEffect } from 'react'
import { WebSocketContext } from '../WebSocket';
import { useDispatch, useSelector, shallowEqual } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { ControlsList } from './ControlsList'
import { Dialog, DialogFooter, IDialogProps } from '@fluentui/react';
import { IControlProps, defaultPixels } from './IControlProps'

export const MyDialog = React.memo<IControlProps>(({control, parentDisabled}) => {

    let disabled = (control.disabled === 'true') || parentDisabled;

    console.log(`render dialog: ${control.i}`);

    const ws = useContext(WebSocketContext);
    const dispatch = useDispatch();
  
    const handleDismiss = () => {
  
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
      ws.pageEventFromWeb(control.i, 'dismiss', control.data)
    }

    useEffect(() => {
        //console.log('Dialog mount')
        return () => {
            const layers = document.body.getElementsByClassName("ms-Layer--fixed")
            for (let i = 0; i < layers.length; i++) {
                let layer: Element = layers[i];
                if (!layer.hasChildNodes()) {
                    //console.log('Dialog dismount', layer)
                    document.body.removeChild(layer);
                }
            }
        };
      }, [control, ws]);

    // dialog props
    const props: IDialogProps = {
        hidden: control.open !== 'true',
        minWidth: control.width ? defaultPixels(control.width) : undefined,
        maxWidth: control.maxwidth ? defaultPixels(control.maxwidth) : undefined,
        modalProps: {
            topOffsetFixed: control.fixedtop === 'true',
            isBlocking: control.blocking === 'true',
        },
        dialogContentProps: {
            type: control.largeheader === 'true' ? 1 : control.close === 'true' ? 2 : 0,
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