import React from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { Checkbox, ICheckbox, ICheckboxProps } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels, getId, isTrue } from './Utils'

export const MyCheckbox = React.memo<IControlProps>(({ control, parentDisabled }) => {

  const ws = React.useContext(WebSocketContext);
  const dispatch = useDispatch();

  let disabled = isTrue(control.disabled) || parentDisabled;

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
    id: getId(control.f ? control.f : control.i),
    checked: isTrue(control.value),
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

  const ctrlRef = React.useRef<ICheckbox | null>(null);
  const [focused, setFocused] = React.useState<boolean>(false);

  React.useEffect(() => {
    if (isTrue(control.focused) && !focused) {
      ctrlRef.current?.focus();
      setFocused(true);
    }
  }, [control.focused, focused]);

  return <Checkbox componentRef={ctrlRef} {...checkboxProps} onChange={handleChange} />;
})