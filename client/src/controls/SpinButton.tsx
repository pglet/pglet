import React from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { SpinButton, ISpinButtonProps, Position, ISpinButton } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels, getId, isTrue } from './Utils'

export const MySpinButton = React.memo<IControlProps>(({ control, parentDisabled }) => {

  const ws = React.useContext(WebSocketContext);
  const dispatch = useDispatch();

  let disabled = isTrue(control.disabled) || parentDisabled;

  const handleChange = (event: React.SyntheticEvent<HTMLElement>, newValue?: string) => {
    //console.log(newValue);

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

  //console.log(`render Text: ${control.i}`);

  // https://developer.microsoft.com/en-us/fluentui#/controls/web/references/ifontstyles#IFontStyles

  let labelPosition: Position | undefined = undefined;
  switch (control.labelposition ? control.labelposition.toLowerCase() : '') {
    case 'left': labelPosition = Position.start; break;
    case 'top': labelPosition = Position.top; break;
    case 'right': labelPosition = Position.end; break;
    case 'bottom': labelPosition = Position.bottom; break;
  }

  const props: ISpinButtonProps = {
    id: getId(control.f ? control.f : control.i),
    defaultValue: control.value ? control.value : undefined,
    label: control.label ? control.label : undefined,
    labelPosition: labelPosition,
    min: control.min ? parseInt(control.min) : undefined,
    max: control.max ? parseInt(control.max) : undefined,
    step: control.step ? parseFloat(control.step) : undefined,
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

  if (control.icon) {
    props.iconProps = {
      iconName: control.icon
    }
  }

  const handleFocus = () => {
    ws.pageEventFromWeb(control.i, 'focus', control.data)
  }

  const handleBlur = () => {
    ws.pageEventFromWeb(control.i, 'blur', control.data)
  }

  const ctrlRef = React.useRef<ISpinButton | null>(null);
  const [focused, setFocused] = React.useState<boolean>(false);

  React.useEffect(() => {
    if (isTrue(control.focused) && !focused) {
      ctrlRef.current?.focus();
      setFocused(true);
    }
  }, [control.focused, focused]);

  return <SpinButton
    componentRef={ctrlRef}
    {...props}
    onChange={handleChange}
    onFocus={handleFocus}
    onBlur={handleBlur}
  />
})