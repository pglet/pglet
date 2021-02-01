import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { SearchBox, ISearchBoxProps } from '@fluentui/react';
import { IControlProps, defaultPixels } from './IControlProps'

export const Searchbox = React.memo<IControlProps>(({ control, parentDisabled }) => {

  //console.log(`render Textbox: ${control.i}`);
  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const dispatch = useDispatch();

  const handleChange = (event?: React.ChangeEvent<HTMLInputElement>, newValue?: string) => {

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

  const handleClear = (ev?: any) => {
    ws.pageEventFromWeb(control.i, 'clear', "")
  }

  const handleEscape = (ev?: any) => {
    ws.pageEventFromWeb(control.i, 'escape', "")
  }

  const handleSearch = (newValue: any) => {
    ws.pageEventFromWeb(control.i, 'search', newValue)
  }  

  // https://stackoverflow.com/questions/56696136/how-to-change-iconbutton-color

  const props: ISearchBoxProps = {
    value: control.value ? control.value : "",
    placeholder: control.placeholder ? control.placeholder : null,
    underlined: control.underlined === 'true',
    disabled: disabled,
    styles: {
      icon: {
        color: control.iconcolor !== undefined ? control.iconcolor : undefined,
      },
      root: {
        width: control.width !== undefined ? defaultPixels(control.width) : undefined,
        height: control.height !== undefined ? defaultPixels(control.height) : undefined,
        padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined
      },
    }
  };

  if (control.icon) {
    props.iconProps = {
      iconName: control.icon
    }
  }

  return <SearchBox {...props}
    onChange={handleChange}
    onClear={handleClear}
    onEscape={handleEscape}
    onSearch={handleSearch} />;
})