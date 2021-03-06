import React from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch, shallowEqual, useSelector } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { Pivot, PivotItem, IPivotProps, mergeStyles } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { ControlsList } from './ControlsList'
import { defaultPixels, isTrue } from './Utils'

export const Tabs = React.memo<IControlProps>(({control, parentDisabled}) => {

  const ws = React.useContext(WebSocketContext);
  const dispatch = useDispatch();

  let disabled = isTrue(control.disabled) || parentDisabled;
  
  const handleChange = (item?: PivotItem, ev?: React.MouseEvent<HTMLElement>) => {

    //console.log("pivot item selected:", item.props);

    let selectedKey = item!.props.itemKey as string

    let payload: any = {}
    if (control.f) {
      // binding redirect
      const p = control.f.split('|')
      payload["i"] = p[0]
      payload[p[1]] = selectedKey
    } else {
      // unbound control
      payload["i"] = control.i
      payload["value"] = selectedKey
    }

    dispatch(changeProps([payload]));
    ws.updateControlProps([payload]);
    ws.pageEventFromWeb(control.i, 'change', control.data ? `${control.data}|${selectedKey}` : selectedKey)
  }

  const pivotClassName = mergeStyles({
    width: control.width ? defaultPixels(control.width) : undefined,
    height: control.height ? defaultPixels(control.height) : undefined,
  });  

  const pivotProps: IPivotProps = {
    className: pivotClassName,
    linkFormat: isTrue(control.solid) ? 'tabs' : undefined,
    styles: {
      root: {
        marginBottom: control.margin ? defaultPixels(control.margin) : undefined,
      }
    }
  };

  const tabControls = useSelector<any, any[]>((state: any) => {
    return (control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId]))
          .filter((tc: any) => tc.t === 'tab' && tc.visible !== "false")
          .map((tab:any) => ({
            i: tab.i,
            props: {
              itemKey: tab.key !== undefined ? tab.key : (tab.text ? tab.text : tab.i),
              headerText: tab.text ? tab.text : (tab.key ? tab.key : tab.i),
              itemIcon: tab.icon ? tab.icon : undefined,
              itemCount: tab.count !== undefined ? tab.count : undefined
            },
            controls: (tab.children !== undefined ? tab.children : tab.c.map((childId: any) => state.page.controls[childId]))
          }));
  }, shallowEqual)

  pivotProps.selectedKey = control.value !== undefined ? control.value : "";

  return <Pivot {...pivotProps} onLinkClick={handleChange}>
    {tabControls.map((tab: any) =>
    <PivotItem key={tab.i} {...tab.props}>
      <ControlsList controls={tab.controls} parentDisabled={disabled} />
    </PivotItem>)}
  </Pivot>;
})