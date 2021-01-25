import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { Toggle, IToggleProps } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const MyToggle = React.memo<IControlProps>(({control, parentDisabled}) => {

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

  const toggleProps: IToggleProps = {
    checked: control.value === "true",
    inlineLabel: control.inline === "true",
    label: control.label ? control.label : undefined,
    onText: control.ontext ? control.ontext : undefined,
    offText: control.offtext ? control.offtext : undefined,
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

  return <Toggle {...toggleProps} onChange={handleChange} />;
})