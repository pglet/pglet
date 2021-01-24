import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { Slider, ISliderProps } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const MySlider = React.memo<IControlProps>(({control, parentDisabled}) => {

  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const dispatch = useDispatch();

  const handleChange = (value: number) => {

    const payload = [
      {
        i: control.i,
        "value": value
      }
    ];

    dispatch(changeProps(payload));
    ws.updateControlProps(payload);
    ws.pageEventFromWeb(control.i, 'change', String(value))
  }

  //console.log(`render Text: ${control.i}`);

  // https://developer.microsoft.com/en-us/fluentui#/controls/web/references/ifontstyles#IFontStyles

  const sliderProps: ISliderProps = {
    value: control.value ? parseInt(control.value) : undefined,
    label: control.label ? control.label : undefined,
    min: control.min ? parseInt(control.min) : undefined,
    max: control.max ? parseInt(control.max) : undefined,
    step: control.step ? parseInt(control.step) : undefined,
    showValue: control.showvalue === 'true',
    vertical: control.vertical === 'true',
    disabled: disabled,
    valueFormat: (value) => {
      const format = control.valueformat ? control.valueformat : '{value}';
      return format.replace('{value}', value);
    },
    styles: {
      root: {
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined
      }
    }
  };

  if (sliderProps.min! < 0) {
    sliderProps.originFromZero = true;
  }

  return <Slider {...sliderProps} onChange={handleChange}></Slider>;
})