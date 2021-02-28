import React from 'react'
import { Callout, DirectionalHint, ICalloutProps } from '@fluentui/react';
import { shallowEqual, useSelector, useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { ControlsList } from './ControlsList'
import { WebSocketContext } from '../WebSocket';
import { IControlProps } from './Control.types'
import { defaultPixels, getId, parseNumber } from './Utils'

export const MyCallout = React.memo<IControlProps>(({ control, parentDisabled }) => {

  const ws = React.useContext(WebSocketContext);

  let disabled = (control.disabled === 'true') || parentDisabled;
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
      payload["visible"] = val
    }

    dispatch(changeProps([payload]));
    ws.updateControlProps([payload]);
    ws.pageEventFromWeb(control.i, 'dismiss', control.data)
  }

  let position: DirectionalHint = DirectionalHint.bottomAutoEdge;
  switch (control.position ? control.position.toLowerCase() : '') {
    case 'topleft': position = DirectionalHint.topLeftEdge; break;
    case 'topcenter': position = DirectionalHint.topCenter; break;
    case 'topright': position = DirectionalHint.topRightEdge; break;
    case 'topauto': position = DirectionalHint.topAutoEdge; break;
    case 'bottomleft': position = DirectionalHint.bottomLeftEdge; break;
    case 'bottomcenter': position = DirectionalHint.bottomCenter; break;
    case 'bottomright': position = DirectionalHint.bottomRightEdge; break;
    case 'bottomauto': position = DirectionalHint.bottomAutoEdge; break;
    case 'lefttop': position = DirectionalHint.leftTopEdge; break;
    case 'leftcenter': position = DirectionalHint.leftCenter; break;
    case 'leftbottom': position = DirectionalHint.leftBottomEdge; break;
    case 'righttop': position = DirectionalHint.rightTopEdge; break;
    case 'rightcenter': position = DirectionalHint.rightCenter; break;
    case 'rightbottom': position = DirectionalHint.rightBottomEdge; break;
  }  

  const props: ICalloutProps = {
    gapSpace: control.gap !== undefined ? parseNumber(control.gap) : undefined,
    beakWidth: control.beakwidth !== undefined ? parseNumber(control.beakwidth) : undefined,
    minPagePadding: control.pagepadding !== undefined ? parseNumber(control.pagepadding) : undefined,
    setInitialFocus: control.focus === "true",
    coverTarget: control.cover === "true",
    isBeakVisible: control.beak === "false" ? false : true,
    directionalHint: position,
    onDismiss: handleDismiss,
    styles: {
      root: {
        width: control.width !== undefined ? defaultPixels(control.width) : undefined,
        height: control.height !== undefined ? defaultPixels(control.height) : undefined,
        padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined,
      }
    }
  };

  if (control.target) {
    const targetId = getId(control.target);
    if (document.getElementById(targetId)) {
      // ID
      props.target = `#${targetId}`;
    } else {
      // ClassName
      props.target = `.${targetId}`;
    }
  }

  const childControls = useSelector((state: any) => {
    return control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId])
  }, shallowEqual);

  return <Callout {...props}>
    <ControlsList controls={childControls} parentDisabled={disabled} />
  </Callout>;
})