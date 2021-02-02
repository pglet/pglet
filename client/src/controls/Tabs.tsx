import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch, shallowEqual, useSelector } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { Pivot, PivotItem, IPivotProps, mergeStyles } from '@fluentui/react';
import { IControlProps, defaultPixels } from './IControlProps'
import { ControlsList } from './ControlsList'

export const Tabs = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Dropdown: ${control.i}`);
  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const dispatch = useDispatch();
  
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
    linkFormat: control.solid === 'true' ? 'tabs' : undefined,
    styles: {
      root: {
        marginBottom: control.margin ? defaultPixels(control.margin) : undefined,
      }
    }
  };

  const tabControls = useSelector<any, any[]>((state: any) => {
    return control.c.map((childId: any) =>
          state.page.controls[childId])
          .filter((tc: any) => tc.t === 'tab' && tc.visible !== "false")
          .map((tab:any) => ({
            i: tab.i,
            props: {
              itemKey: tab.key ? tab.key : tab.i,
              headerText: tab.text ? tab.text : (tab.key ? tab.key : tab.i),
              itemIcon: tab.icon ? tab.icon : undefined,
              itemCount: tab.count !== undefined ? tab.count : undefined
            },
            controls: tab.c.map((childId: any) => state.page.controls[childId])
          }));
  }, shallowEqual)

  if (control.value) {
    pivotProps.selectedKey = control.value;
  }

  return <Pivot {...pivotProps} onLinkClick={handleChange}>
    {tabControls.map((tab: any) =>
    <PivotItem key={tab.i} {...tab.props}>
      <ControlsList controls={tab.controls} parentDisabled={disabled} />
    </PivotItem>)}
  </Pivot>;
})