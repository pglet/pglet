import React from 'react'
import { DatePicker, getId, IDatePickerProps } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels, isTrue, parseDate } from './Utils'
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux';
import { changeProps } from '../slices/pageSlice';

export const MyDatePicker = React.memo<IControlProps>(({control, parentDisabled}) => {

  const ws = React.useContext(WebSocketContext);
  const dispatch = useDispatch();
  let disabled = isTrue(control.disabled) || parentDisabled;

  const handleSelectDate = (date: Date | null | undefined) => {

    let newValue = date?.toString();
    if (newValue === undefined) {
      newValue = "";
    }
    
    let payload: any = {}
    if (control.f) {
      // binding redirect
      const p = control.f.split('|')
      payload["i"] = p[0]
      payload[p[1]] = newValue
    } else {
      // unbound control
      payload["i"] = control.i
      payload["value"] = newValue
    }

    dispatch(changeProps([payload]));
    ws.updateControlProps([payload]);
    ws.pageEventFromWeb(control.i, 'change', control.data ? `${control.data}|${newValue!}` : newValue!)
  }  
  
  const pickerProps: IDatePickerProps = {
    id: getId(control.f ? control.f : control.i),
    value: control.value ? parseDate(control.value) : undefined,
    label: control.label ? control.label : null,
    placeholder: control.placeholder ? control.placeholder : null,
    isRequired: isTrue(control.required),
    allowTextInput: isTrue(control.allowtextinput),
    borderless: isTrue(control.borderless),
    underlined: isTrue(control.underlined),
    disabled: disabled,
    styles: {
      root: {
        width: control.width !== undefined ? defaultPixels(control.width) : undefined
      }
    }    
  };

  return <DatePicker {...pickerProps} onSelectDate={handleSelectDate} />;
})