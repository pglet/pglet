import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { TextField, ITextFieldProps } from '@fluentui/react';
import { IControlProps, defaultPixels } from './IControlProps'

export const Textbox = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Textbox: ${control.i}`);
  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const dispatch = useDispatch();
  
  const handleChange = (event: React.FormEvent<HTMLInputElement | HTMLTextAreaElement>, newValue?: string) => {

    const payload = [
      {
        i: control.i,
        "value": newValue
      }
    ];

    dispatch(changeProps(payload));
    ws.updateControlProps(payload);

    if (control.onchange === 'true') {
      ws.pageEventFromWeb(control.i, 'change', newValue!)
    }
  }

  const textFieldProps: ITextFieldProps = {
    value: control.value ? control.value : "",
    label: control.label ? control.label : null,
    placeholder: control.placeholder ? control.placeholder : null,
    errorMessage: control.errormessage ? control.errormessage : null,
    description: control.description ? control.description : null,
    multiline: control.multiline ? true : false,
    type: control.password ? "password" : undefined,
    canRevealPassword: control.password ? true : undefined,
    required: control.required ? true : undefined,
    disabled: disabled,
    styles: {
      root: {
        width: control.width !== undefined ? defaultPixels(control.width) : undefined,
        height: control.height !== undefined ? defaultPixels(control.height) : undefined,
        padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined  
      },
      field: {
        textAlign: control.align !== undefined ? control.align : undefined,
      },
    }
  };

  return <TextField {...textFieldProps} onChange={handleChange} />;
})