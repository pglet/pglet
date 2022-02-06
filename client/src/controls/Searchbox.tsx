import React from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { SearchBox, ISearchBoxProps, useTheme, ISearchBox } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { getThemeColor, defaultPixels, getId, isTrue } from './Utils'

export const Searchbox = React.memo<IControlProps>(({ control, parentDisabled }) => {

  //console.log("Render Searchbox", control.i);

  const ws = React.useContext(WebSocketContext);
  const dispatch = useDispatch();
  const theme = useTheme();

  let disabled = isTrue(control.disabled) || parentDisabled;

  let _lastChangeValue: string | undefined;

  const handleChange = (event?: React.ChangeEvent<HTMLInputElement>, newValue?: string) => {

    if (newValue === _lastChangeValue) {
      _lastChangeValue = undefined;
      return;
    }
    _lastChangeValue = newValue;

    //console.log("Searchbox handleChange:", newValue);

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

  const handleClear = (ev?: any) => {
    ws.pageEventFromWeb(control.i, 'clear', control.data)
  }

  const handleSearch = (newValue: any) => {
    ws.pageEventFromWeb(control.i, 'search', control.data ? `${control.data}|${newValue!}` : newValue!)
  }

  // https://stackoverflow.com/questions/56696136/how-to-change-iconbutton-color

  const props: ISearchBoxProps = {
    id: getId(control.f ? control.f : control.i),
    value: control.value ? control.value : "",
    placeholder: control.placeholder ? control.placeholder : null,
    underlined: isTrue(control.underlined),
    disabled: disabled,
    styles: {
      icon: {
        color: control.iconcolor !== undefined ? getThemeColor(theme, control.iconcolor) : undefined,
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

  const ctrlRef = React.useRef<ISearchBox | null>(null);
  const [focused, setFocused] = React.useState<boolean>(false);

  React.useEffect(() => {
    if (isTrue(control.focused) && !focused) {
      ctrlRef.current?.focus();
      setFocused(true);
    }
  }, [control.focused, focused]);

  return <SearchBox componentRef={ctrlRef} {...props}
    onChange={handleChange}
    onClear={handleClear}
    onSearch={handleSearch} />;
})