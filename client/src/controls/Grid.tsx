import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch, shallowEqual, useSelector } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { DetailsList, IDetailsListProps, IColumn, Selection } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const Grid = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Dropdown: ${control.i}`);
  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const dispatch = useDispatch();
  
  // const handleChange = (event: React.FormEvent<HTMLDivElement>, option?: IDropdownOption, index?: number) => {

  //   //console.log("DROPDOWN:", option);

  //   let selectedKey = option!.key as string

  //   const payload = [
  //     {
  //       i: control.i,
  //       "value": selectedKey
  //     }
  //   ];

  //   dispatch(changeProps(payload));
  //   ws.updateControlProps(payload);
  //   ws.pageEventFromWeb(control.i, 'change', selectedKey)
  // }

  console.log("GRID - START REDNER");

  const _onItemInvoked = (item: any) => {
    alert(`Item invoked: ${item.name}`);
  }

  const _selection = new Selection({
    onSelectionChanged: () => {
      console.log(_selection.getSelection());
    },
  });  

  const gridConfig = useSelector<any, any>((state: any) => {
    const columns: IColumn[] = control.c.map((childId: any) => state.page.controls[childId])
      .filter((c: any) => c.t === 'columns').map((columns: any) =>
        columns.c.map((childId: any) => state.page.controls[childId]))
        .reduce((acc: any, columns: any) => ([...acc, ...columns])).map((cc: any) => {
          return {
            key: cc.i,
            name: cc.name,
            iconName: cc.iconname,
            isIconOnly: cc.icononly === 'true',
            fieldName: cc.fieldname ? cc.fieldname.toLowerCase() : undefined,
            isResizable: cc.resizable === 'true',
            //isSorted: false,
            minWidth: cc.minwidth ? parseInt(cc.minwidth) : undefined,
            maxWidth: cc.maxwidth ? parseInt(cc.maxwidth) : undefined,
            //onColumnClick: _onColumnClick
          }
        });

  const items = control.c.map((childId: any) => state.page.controls[childId])
  .filter((c: any) => c.t === 'items').map((items: any) =>
    items.c.map((childId: any) => state.page.controls[childId]))
    .reduce((acc: any, items: any) => ([...acc, ...items]));        

    return {
      columns,
      items
    }
  }, shallowEqual);

  const gridProps: IDetailsListProps = {
    columns: gridConfig.columns,
    items: gridConfig.items,
    compact: false,
    isHeaderVisible: true,
    onItemInvoked: _onItemInvoked,
    selection: _selection,
    selectionPreservedOnEmptyClick: true,
    //disabled: disabled,
    styles: {
      root: {
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined   
      }
    }
  };  

  console.log("RENDER:", gridConfig);

  // if (control.value) {
  //   dropdownProps.defaultSelectedKey = control.value;
  // }

  return <DetailsList {...gridProps} />;
})