import React from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch, shallowEqual, useSelector } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { ComboBox, IComboBox, IComboBoxOption, IComboBoxProps, SelectableOptionMenuItemType } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels, getId, isFalse, isTrue } from './Utils'

export const MyComboBox = React.memo<IControlProps>(({ control, parentDisabled }) => {

  const ws = React.useContext(WebSocketContext);
  const dispatch = useDispatch();

  const disabled = isTrue(control.disabled) || parentDisabled;
  const multiSelect = isTrue(control.multiselect)

  const [selectedKeys, setSelectedKeys] = React.useState<string[]>([]);
  const [selectAll, setSelectAll] = React.useState<boolean | undefined>();

  const comboboxProps: IComboBoxProps = {
    id: getId(control.f ? control.f : control.i),
    label: control.label ? control.label : null,
    placeholder: control.placeholder ? control.placeholder : null,
    errorMessage: control.errormessage ? control.errormessage : null,
    options: [],
    multiSelect: multiSelect,
    allowFreeform: isTrue(control.allowfreeform),
    autoComplete: isFalse(control.autocomplete) ? "off" : "on",
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
      case 'divider': return SelectableOptionMenuItemType.Divider;
      case 'header': return SelectableOptionMenuItemType.Header;
      case 'selectall':
      case 'select_all': return SelectableOptionMenuItemType.SelectAll;
      default: return SelectableOptionMenuItemType.Normal;
    }
  }

  comboboxProps.options = useSelector<any, IComboBoxOption[]>((state: any) =>
    (control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId]))
      .filter((oc: any) => oc.t === 'option')
      .map((oc: any) => ({
        key: oc.key ?? oc.text,
        text: oc.text ?? oc.key,
        itemType: getItemType(oc.itemtype),
        disabled: isTrue(oc.disabled)
      })), shallowEqual);

  const handleChange = (event: React.FormEvent<IComboBox>, option?: IComboBoxOption | undefined, index?: number | undefined, value?: string | undefined) => {

    //console.log("ComboBox.change:", option);

    let selectedKey = option?.key as string;

    if (multiSelect) {
      let keys = selectedKeys;
      //console.log("Keys:", keys)

      if (option) {
        if (option?.itemType === SelectableOptionMenuItemType.SelectAll) {
          setSelectAll(option?.selected)
          keys = option?.selected ? comboboxProps.options
            .filter(option => option.itemType === SelectableOptionMenuItemType.Normal)
            .map(option => option.key as string) : [];
        } else {
          keys = option?.selected ? [...keys, option!.key as string] : keys.filter(k => k !== option!.key);
          setSelectedKeys(keys);
          setSelectAll(keys.length === comboboxProps.options
            .filter(option => option.itemType === SelectableOptionMenuItemType.Normal).length)
        }
      }

      selectedKey = keys.join(",")
    }

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

  // select keys or value
  const values: string[] = (control.value !== undefined) ? control.value.split(",").map((item: string) => {
    return item.trim();
  }) : [];
  const value = control.value ?? "";

  //console.log("--selectAll:", selectAll)
  if (selectAll) {
    values.push(comboboxProps.options
      .filter(option => option.itemType === SelectableOptionMenuItemType.SelectAll)[0].key as string)
  }

  //console.log("values:", control.i, values)
  //console.log("value:", control.i, values)
  comboboxProps.selectedKey = multiSelect ? values : value;

  const handleFocus = () => {
    ws.pageEventFromWeb(control.i, 'focus', control.data)
  }

  const handleBlur = () => {
    ws.pageEventFromWeb(control.i, 'blur', control.data)
  }

  const ctrlRef = React.useRef<IComboBox | null>(null);
  const [focused, setFocused] = React.useState<boolean>(false);

  React.useEffect(() => {

    if (multiSelect) {
      setSelectedKeys(values);
      setSelectAll(values.length === comboboxProps.options
        .filter(option => option.itemType === SelectableOptionMenuItemType.Normal).length)
    }

    // focus control
    if (isTrue(control.focused) && !focused) {
      ctrlRef.current?.focus();
      setFocused(true);
    }
    // eslint-disable-next-line
  }, [control.focused, focused]);

  return <ComboBox
    componentRef={ctrlRef}
    {...comboboxProps}
    onChange={handleChange}
    onFocus={handleFocus}
    onBlur={handleBlur}
  />
})