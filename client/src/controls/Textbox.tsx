import React from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { TextField, ITextFieldProps, useTheme } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels, getId, getThemeColor, isTrue } from './Utils'

export const Textbox = React.memo<IControlProps>(({control, parentDisabled}) => {

  const ws = React.useContext(WebSocketContext);
  const dispatch = useDispatch();
  const theme = useTheme();

  let disabled = isTrue(control.disabled) || parentDisabled;
  
  const handleChange = (event: React.FormEvent<HTMLInputElement | HTMLTextAreaElement>, newValue?: string) => {

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

    if (isTrue(control.onchange)) {
      ws.pageEventFromWeb(control.i, 'change', control.data ? `${control.data}|${newValue!}` : newValue!)
    }
  }

  const textFieldProps: ITextFieldProps = {
    id: getId(control.f ? control.f : control.i),
    value: control.value ? control.value : "",
    label: control.label ? control.label : null,
    placeholder: control.placeholder ? control.placeholder : null,
    errorMessage: control.errormessage ? control.errormessage : null,
    description: control.description ? control.description : null,
    multiline: isTrue(control.multiline),
    type: isTrue(control.password) ? "password" : undefined,
    canRevealPassword: isTrue(control.password),
    required: isTrue(control.required),
    readOnly: isTrue(control.readonly),
    autoAdjustHeight: isTrue(control.autoadjustheight),
    underlined: isTrue(control.underlined),
    borderless: isTrue(control.borderless),
    prefix: control.prefix ? control.prefix : undefined,
    suffix: control.suffix ? control.suffix : undefined,
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

  if (control.icon) {
    textFieldProps.iconProps = {
      iconName: control.icon
    }
    if (control.iconcolor !== undefined) {
      textFieldProps.iconProps!.styles = {
          root: {
              color: getThemeColor(theme, control.iconcolor)
          }
      }
    }    
  }

  return <TextField {...textFieldProps} onChange={handleChange} />;
})