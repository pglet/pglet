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
    linkFormat: control.solid === 'true' ? 'tabs' : undefined,
    styles: {
      root: {
        marginBottom: control.margin ? control.margin : undefined,
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