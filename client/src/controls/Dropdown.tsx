import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch, shallowEqual, useSelector } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { Dropdown, IDropdownOption, IDropdownProps } from '@fluentui/react';
import { IControlProps, defaultPixels } from './IControlProps'

export const MyDropdown = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Dropdown: ${control.i}`);
  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const dispatch = useDispatch();
  
  const handleChange = (event: React.FormEvent<HTMLDivElement>, option?: IDropdownOption, index?: number) => {

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

  const dropdownProps: IDropdownProps = {
    label: control.label ? control.label : null,
    placeholder: control.placeholder ? control.placeholder : null,
    errorMessage: control.errormessage ? control.errormessage : null,
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

  dropdownProps.options = useSelector<any, IDropdownOption[]>((state: any) =>
    (control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId]))
      .filter((oc: any) => oc.t === 'option')
      .map((oc: any) => ({
        key: oc.key ? oc.key : oc.text,
        text: oc.text ? oc.text : oc.key
      })), shallowEqual);  

  if (control.value) {
    dropdownProps.defaultSelectedKey = control.value;
  }

  return <Dropdown {...dropdownProps} onChange={handleChange} />;
})