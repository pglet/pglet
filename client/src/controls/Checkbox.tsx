import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { Checkbox, ICheckboxProps } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels } from './Utils'

export const MyCheckbox = React.memo<IControlProps>(({control, parentDisabled}) => {

  const ws = useContext(WebSocketContext);
  const dispatch = useDispatch();

  let disabled = (control.disabled === 'true') || parentDisabled;

  const handleChange = (event?: React.FormEvent<HTMLElement | HTMLInputElement>, checked?: boolean) => {

    if (checked !== undefined) {

      const val = checked.toString();

      let payload: any = {}
      if (control.f) {
        // binding redirect
        const p = control.f.split('|')
        payload["i"] = p[0]
        payload[p[1]] = val
      } else {
        // unbound control
        payload["i"] = control.i
        payload["value"] = val
      }
  
      dispatch(changeProps([payload]));
      ws.updateControlProps([payload]);
      ws.pageEventFromWeb(control.i, 'change', control.data ? `${control.data}|${val}` : val)
    }
  }

  const checkboxProps: ICheckboxProps = {
    checked: control.value === "true",
    label: control.label ? control.label : null,
    boxSide: control.boxside ? control.boxside : 'start',
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

  return <Checkbox {...checkboxProps} onChange={handleChange} />;
})