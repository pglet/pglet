import React from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { Toggle, IToggleProps, IToggle } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels, getId, isTrue } from './Utils'

export const MyToggle = React.memo<IControlProps>(({ control, parentDisabled }) => {

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

  const toggleProps: IToggleProps = {
    id: getId(control.f ? control.f : control.i),
    checked: isTrue(control.value),
    inlineLabel: isTrue(control.inline),
    label: control.label ? control.label : undefined,
    onText: control.ontext ? control.ontext : undefined,
    offText: control.offtext ? control.offtext : undefined,
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

  const handleFocus = () => {
    ws.pageEventFromWeb(control.i, 'focus', control.data)
  }

  const handleBlur = () => {
    ws.pageEventFromWeb(control.i, 'blur', control.data)
  }

  const ctrlRef = React.useRef<IToggle | null>(null);
  const [focused, setFocused] = React.useState<boolean>(false);

  React.useEffect(() => {
    if (isTrue(control.focused) && !focused) {
      ctrlRef.current?.focus();
      setFocused(true);
    }
  }, [control.focused, focused]);

  return <Toggle
    componentRef={ctrlRef}
    {...toggleProps}
    onChange={handleChange}
    onFocus={handleFocus}
    onBlur={handleBlur}
  />
})