import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch, shallowEqual, useSelector } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { Dropdown, IDropdownOption, IDropdownProps } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const MyDropdown = React.memo<IControlProps>(({control}) => {

  //console.log(`render Dropdown: ${control.i}`);

  const ws = useContext(WebSocketContext);

  const dispatch = useDispatch();
  
  const handleChange = (event: React.FormEvent<HTMLDivElement>, option?: IDropdownOption, index?: number) => {

    //console.log("DROPDOWN:", option);

    const payload = [
      {
        i: control.i,
        "value": option!.key
      }
    ];

    dispatch(changeProps(payload));
    ws.updateControlProps(payload);
  }

  const dropdownProps: IDropdownProps = {
    label: control.label ? control.label : null,
    placeholder: control.placeholder ? control.placeholder : null,
    errorMessage: control.errormessage ? control.errormessage : null,
    options: [],
    styles: {
      root: {
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined   
      }
    }
  };

  dropdownProps.options = useSelector<any, IDropdownOption[]>((state: any) => control.c.map((childId: any) =>
    state.page.controls[childId])
      .filter((oc: any) => oc.t === 'option')
      .map((oc: any) => ({ key: oc.key, text: oc.text})), shallowEqual);  

  if (control.value) {
    dropdownProps.defaultSelectedKey = control.value;
  }

  return <Dropdown {...dropdownProps} onChange={handleChange} />;
})