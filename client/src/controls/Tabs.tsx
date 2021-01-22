import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch, shallowEqual, useSelector } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { Pivot, PivotItem, IPivotProps } from '@fluentui/react';
import { IControlProps } from './IControlProps'
import { ControlsList } from './ControlsList'

export const Tabs = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Dropdown: ${control.i}`);
  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const dispatch = useDispatch();
  
  const handleChange = (item?: PivotItem, ev?: React.MouseEvent<HTMLElement>) => {

    //console.log("pivot item selected:", item.props);

    let selectedKey = item!.props.itemKey as string

    const payload = [
      {
        i: control.i,
        "value": selectedKey
      }
    ];

    dispatch(changeProps(payload));
    ws.updateControlProps(payload);
    ws.pageEventFromWeb(control.i, 'change', selectedKey)
  }

  const pivotProps: IPivotProps = {
  };

  const tabControls = useSelector<any, any[]>((state: any) => {
    return control.c.map((childId: any) =>
          state.page.controls[childId]).filter((oc: any) => oc.t === 'tab')
          .map((tab:any) => ({
            i: tab.i,
            props: {
              itemKey: tab.key,
              headerText: tab.text ? tab.text : tab.key
            },
            controls: tab.c.map((childId: any) => state.page.controls[childId])
          }));
  }, shallowEqual)

  if (control.value) {
    pivotProps.selectedKey = control.value;
  }

  return <Pivot {...pivotProps} onLinkClick={handleChange}>
    {tabControls.map(tab =>
    <PivotItem key={tab.i} {...tab.props}>
      <ControlsList controls={tab.controls} parentDisabled={disabled} />
    </PivotItem>)}
  </Pivot>;
})