import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { Checkbox, ICheckboxProps } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const MyCheckbox = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Checkbox: ${control.i}`);

  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const dispatch = useDispatch();
  
  const handleChange = (event?: React.FormEvent<HTMLElement | HTMLInputElement>, checked?: boolean) => {

    if (checked !== undefined) {
      const payload = [
        {
          i: control.i,
          "value": checked.toString()
        }
      ];
  
      dispatch(changeProps(payload));
      ws.updateControlProps(payload);
      ws.pageEventFromWeb(control.i, 'change', checked.toString())
    }
  }

  const checkboxProps: ICheckboxProps = {
    checked: control.value === "true",
    label: control.label ? control.label : null,
    boxSide: control.boxside ? control.boxside : 'start',
    disabled: disabled,
    styles: {
      root: {
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined   
      }
    }
  };

  return <Checkbox {...checkboxProps} onChange={handleChange} />;
})