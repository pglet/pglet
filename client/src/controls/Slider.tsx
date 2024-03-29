import React from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { Slider, ISliderProps, ISlider } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels, getId, isTrue } from './Utils'

export const MySlider = React.memo<IControlProps>(({ control, parentDisabled }) => {

  const ws = React.useContext(WebSocketContext);
  const dispatch = useDispatch();
  const [prevValue, setPrevValue] = React.useState<number | null>(null);

  let disabled = isTrue(control.disabled) || parentDisabled;

  const handleChange = (value: number) => {

    if (prevValue === value) {
      return
    }

    const val = String(value)

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

    setPrevValue(value)
  }

  //console.log(`render Text: ${control.i}`);

  // https://developer.microsoft.com/en-us/fluentui#/controls/web/references/ifontstyles#IFontStyles

  const sliderProps: ISliderProps = {
    id: getId(control.f ? control.f : control.i),
    value: control.value ? parseInt(control.value) : undefined,
    label: control.label ? control.label : undefined,
    min: control.min ? parseInt(control.min) : undefined,
    max: control.max ? parseInt(control.max) : undefined,
    step: control.step ? parseInt(control.step) : undefined,
    showValue: isTrue(control.showvalue),
    vertical: isTrue(control.vertical),
    disabled: disabled,
    valueFormat: (value) => {
      const format = control.valueformat ? control.valueformat : '{value}';
      return format.replace('{value}', value);
    },
    styles: {
      root: {
        width: control.width !== undefined ? defaultPixels(control.width) : undefined,
        height: control.height !== undefined ? defaultPixels(control.height) : undefined,
        padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined
      }
    }
  };

  if (sliderProps.min! < 0) {
    sliderProps.originFromZero = true;
  }

  const handleFocus = () => {
    ws.pageEventFromWeb(control.i, 'focus', control.data)
  }

  const handleBlur = () => {
    ws.pageEventFromWeb(control.i, 'blur', control.data)
  }

  const ctrlRef = React.useRef<ISlider | null>(null);
  const [focused, setFocused] = React.useState<boolean>(false);

  React.useEffect(() => {
    if (isTrue(control.focused) && !focused) {
      ctrlRef.current?.focus();
      setFocused(true);
    }
  }, [control.focused, focused]);

  return <Slider
    componentRef={ctrlRef}
    {...sliderProps}
    onChange={handleChange}
    onFocus={handleFocus}
    onBlur={handleBlur}
  />
})