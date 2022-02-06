import React from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch, shallowEqual, useSelector } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { ChoiceGroup, IChoiceGroup, IChoiceGroupOption, IChoiceGroupProps, useTheme } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { getThemeColor, defaultPixels, getId, isTrue } from './Utils'

export const MyChoiceGroup = React.memo<IControlProps>(({ control, parentDisabled }) => {

  const ws = React.useContext(WebSocketContext);
  const dispatch = useDispatch();
  const theme = useTheme();

  let disabled = isTrue(control.disabled) || parentDisabled;

  const handleChange = (ev?: React.FormEvent<HTMLElement | HTMLInputElement>, option?: IChoiceGroupOption) => {

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

  const choiceProps: IChoiceGroupProps = {
    id: getId(control.f ? control.f : control.i),
    label: control.label ? control.label : null,
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

  choiceProps.options = useSelector<any, IChoiceGroupOption[]>((state: any) =>
    (control.children !== undefined ? control.children : control.c.map((childId: any) => state.page.controls[childId]))
      .filter((oc: any) => oc.t === 'option')
      .map((oc: any) => {
        let option: any = {
          key: oc.key ? oc.key : oc.text,
          text: oc.text ? oc.text : oc.key,
        }
        if (oc.icon) {
          option.iconProps = {
            iconName: oc.icon
          }

          if (oc.iconcolor !== undefined) {
            option.iconProps!.styles = {
              root: {
                color: getThemeColor(theme, oc.iconcolor)
              }
            }
          }
        }
        return option;
      }), shallowEqual);

  choiceProps.selectedKey = control.value !== undefined ? control.value : "";

  const ctrlRef = React.useRef<IChoiceGroup | null>(null);
  const [focused, setFocused] = React.useState<boolean>(false);

  React.useEffect(() => {
    if (isTrue(control.focused) && !focused) {
      ctrlRef.current?.focus();
      setFocused(true);
    }
  }, [control.focused, focused]);

  return <ChoiceGroup componentRef={ctrlRef} {...choiceProps} onChange={handleChange} />;
})