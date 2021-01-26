import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch, shallowEqual, useSelector } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { DetailsList, IDetailsListProps, IColumn, SelectionMode, Selection } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const Grid = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Dropdown: ${control.i}`);
  //let disabled = (control.disabled === 'true') || parentDisabled;

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

  let columns: IColumn[] = [];
  let items = null;

  const _onColumnClick = (ev: React.MouseEvent<HTMLElement>, column: IColumn): void => {
    //console.log(column)
    //console.log("DROPDOWN:", option);

    if ((column as any).onClick) {
      ws.pageEventFromWeb(control.i, 'click', control.data)
      return
    }

    if ((column as any).isSortable === undefined || (column as any).isSortable === 'false') {
      return;
    }

    let payload: any = [];

    columns.forEach(currCol => {
      let pc: any = {
        i: currCol.key
      };

      if (currCol.key === column.key) {
        pc.sorted = column.isSortedDescending! ? "asc" : "desc";
      } else {
        pc.sorted = 'false';
      }
      payload.push(pc);
    })

    //console.log(payload);

    dispatch(changeProps(payload));
    ws.updateControlProps(payload);
  }

  const _onItemInvoked = (item: any) => {
    alert(`Item invoked: ${item.name}`);
  }

  const _selection = new Selection({
    onSelectionChanged: () => {
      console.log(_selection.getSelection());
    },
  });

  columns = useSelector<any, IColumn[]>((state: any) => {
    return control.c.map((childId: any) => state.page.controls[childId])
      .filter((c: any) => c.t === 'columns').map((columns: any) =>
        columns.c.map((childId: any) => state.page.controls[childId]))
        .reduce((acc: any, columns: any) => ([...acc, ...columns])).map((cc: any) => {
          return {
            key: cc.i,
            name: cc.name,
            iconName: cc.icon,
            isIconOnly: cc.icononly === 'true',
            fieldName: cc.fieldname ? cc.fieldname.toLowerCase() : undefined,
            isResizable: cc.resizable === 'true',
            isSortable: cc.sortable,
            isSorted: cc.sorted === 'true' || cc.sorted === 'asc' || cc.sorted === 'desc',
            isSortedDescending: cc.sorted === 'desc',
            minWidth: cc.minwidth ? parseInt(cc.minwidth) : undefined,
            maxWidth: cc.maxwidth ? parseInt(cc.maxwidth) : undefined,
            onClick: cc.onclick === 'true',
            onColumnClick: _onColumnClick
          }
        });
  }, shallowEqual);

  items = useSelector<any, any>((state: any) => {
    return control.c.map((childId: any) => state.page.controls[childId])
    .filter((c: any) => c.t === 'items').map((items: any) =>
      items.c.map((childId: any) => state.page.controls[childId]))
      .reduce((acc: any, items: any) => ([...acc, ...items]));
  }, shallowEqual);

  // sort items
  const sortColumns = columns.filter(c => c.isSorted);
  if (sortColumns.length > 0) {
    const sortColumn = sortColumns[0];
    const key = sortColumn.fieldName!;
    items = items.slice(0).sort((a: any, b: any) => {
      if ((sortColumn as any).isSortable === 'number') {
        return (sortColumn.isSortedDescending ? parseFloat(a[key]) < parseFloat(b[key]) : parseFloat(a[key]) > parseFloat(b[key])) ? 1 : -1;
      } else {
        return (sortColumn.isSortedDescending ? a[key] < b[key] : a[key] > b[key]) ? 1 : -1;
      }
    })
  }

  const gridProps: IDetailsListProps = {
    columns: columns,
    items: items,
    compact: false,
    isHeaderVisible: true,
    onItemInvoked: _onItemInvoked,
    styles: {
      root: {
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined   
      }
    }
  };

  // selection mode
  gridProps.selectionMode = SelectionMode.none;
  if (control.selection === 'single' || control.selection === 'multiple') {
    gridProps.selectionMode = control.selection === 'single' ? SelectionMode.single : SelectionMode.multiple;
    gridProps.selection = _selection;
    gridProps.selectionPreservedOnEmptyClick = control.preserveselection === 'true';
  }

  //console.log("RENDER:", items);

  // if (control.value) {
  //   dropdownProps.defaultSelectedKey = control.value;
  // }

  return <DetailsList {...gridProps} />;
})