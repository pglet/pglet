import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch, shallowEqual, useSelector } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { ChoiceGroup, IChoiceGroupOption, IChoiceGroupProps } from '@fluentui/react';
import { IControlProps, defaultPixels } from './IControlProps'

export const MyChoiceGroup = React.memo<IControlProps>(({ control, parentDisabled }) => {

  //console.log(`render Dropdown: ${control.i}`);
  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const dispatch = useDispatch();

  const handleChange = (ev?: React.FormEvent<HTMLElement | HTMLInputElement>, option?: IChoiceGroupOption) => {

    //console.log("DROPDOWN:", option);

    let selectedKey = option!.key as string

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

  const choiceProps: IChoiceGroupProps = {
    label: control.label ? control.label : null,
    options: [],
    disabled: disabled,
    styles: {
      root: {
        width: control.width !== undefined ? defaultPixels(control.width) : undefined,
        height: control.height !== undefined ? defaultPixels(control.height) : undefined,
        padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined
      }
    }
  };

  choiceProps.options = useSelector<any, IChoiceGroupOption[]>((state: any) =>
    (control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId]))
    .filter((oc: any) => oc.t === 'option')
    .map((oc: any) => {
      let option: any = {
        key: oc.key,
        text: oc.text ? oc.text : oc.key,
      }
      if (oc.icon) {
        option.iconProps = {
          iconName: oc.icon
        }
      }
      return option;
    }), shallowEqual);

  if (control.value) {
    choiceProps.defaultSelectedKey = control.value;
  }

  return <ChoiceGroup {...choiceProps} onChange={handleChange} />;
})