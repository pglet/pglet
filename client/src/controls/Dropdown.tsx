import React from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch, shallowEqual, useSelector } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { Dropdown, DropdownMenuItemType, IDropdown, IDropdownOption, IDropdownProps } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels, getId, isTrue } from './Utils'

export const MyDropdown = React.memo<IControlProps>(({ control, parentDisabled }) => {

  const ws = React.useContext(WebSocketContext);
  const dispatch = useDispatch();

  let disabled = isTrue(control.disabled) || parentDisabled;

  const handleChange = (event: React.FormEvent<HTMLDivElement>, option?: IDropdownOption, index?: number) => {

    //console.log("DROPDOWN:", option);

    let selectedKey = option!.key as string

    let payload: any = {}
    if (control.f) {
      // binding redirect
      const p = control.f.split('|')
      payload["i"] = p[0]
      payload[p[1]] = selectedKey
    } else {
      // unbound control
      payload["i"] = control.i
      payload["value"] = selectedKey
    }

    dispatch(changeProps([payload]));
    ws.updateControlProps([payload]);
    ws.pageEventFromWeb(control.i, 'change', control.data ? `${control.data}|${selectedKey}` : selectedKey)
  }

  const dropdownProps: IDropdownProps = {
    id: getId(control.f ? control.f : control.i),
    label: control.label ? control.label : null,
    placeholder: control.placeholder ? control.placeholder : null,
    errorMessage: control.errormessage ? control.errormessage : null,
    options: [],
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

  const getItemType = (t: string) => {
    switch (t ? t.toLowerCase() : '') {
      case 'divider': return DropdownMenuItemType.Divider;
      case 'header': return DropdownMenuItemType.Header;
      default: return DropdownMenuItemType.Normal;
    }
  }

  dropdownProps.options = useSelector<any, IDropdownOption[]>((state: any) =>
    (control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId]))
      .filter((oc: any) => oc.t === 'option')
      .map((oc: any) => ({
        key: oc.key ? oc.key : oc.text,
        text: oc.text ? oc.text : oc.key,
        itemType: getItemType(oc.itemtype),
        disabled: isTrue(oc.disabled)
      })), shallowEqual);

  dropdownProps.selectedKey = control.value !== undefined ? control.value : "";

  const handleFocus = () => {
    ws.pageEventFromWeb(control.i, 'focus', control.data)
  }

  const handleBlur = () => {
    ws.pageEventFromWeb(control.i, 'blur', control.data)
  }

  const ctrlRef = React.useRef<IDropdown | null>(null);
  const [focused, setFocused] = React.useState<boolean>(false);

  React.useEffect(() => {
    if (isTrue(control.focused) && !focused) {
      ctrlRef.current?.focus();
      setFocused(true);
    }
  }, [control.focused, focused]);

  return <Dropdown
    componentRef={ctrlRef}
    {...dropdownProps}
    onChange={handleChange}
    onFocus={handleFocus}
    onBlur={handleBlur}
  />
})