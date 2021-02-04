import React, { useContext } from 'react';
import { WebSocketContext } from '../WebSocket';
import { useDispatch, shallowEqual, useSelector } from 'react-redux'
import { changeProps } from '../slices/pageSlice'
import { ShimmeredDetailsList, IShimmeredDetailsListProps, IColumn, ColumnActionsMode, SelectionMode, Selection } from '@fluentui/react';
import { IControlProps, defaultPixels } from './IControlProps'
import { ControlsList } from './ControlsList'

export const Grid = React.memo<IControlProps>(({control, parentDisabled}) => {

  //console.log(`render Dropdown: ${control.i}`);
  let disabled = (control.disabled === 'true') || parentDisabled;

  const ws = useContext(WebSocketContext);

  const dispatch = useDispatch();
  
  //console.log("GRID - START REDNER");

  let columns: IColumn[] = [];
  let items: any = null;

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
    ws.pageEventFromWeb(control.i, 'itemInvoke', item.i)
  }

  function cloneControls(controls:any[], item:any): any[] {
    return controls.map((c:any) => {
      let clone: any = {}
      Object.getOwnPropertyNames(c).forEach(p => {
        let val = c[p]
        if (typeof val === 'string' && val.startsWith('{') && val.endsWith('}')) {
          const fieldName = val.substring(1, val.length - 1).toLowerCase()
          val = item[fieldName]
          // set binding redirect
          if (p === "value") {
            clone["f"] = `${item.i}|${fieldName}`
          }
        }
        clone[p] = val
      })
      clone.children = cloneControls(c.children, item)
      return clone
    })
  }

  const _columnOnRender = (item?: any, index?: number, column?: IColumn) => {
    //console.log("Column render:", column?.fieldName, item);
    let template = cloneControls((column as any).template, item);
    return <ControlsList controls={template} parentDisabled={disabled} />
  }

  columns = useSelector<any, IColumn[]>((state: any) => {

    function getTemplateControls(state: any, parent: any): any {
      return parent.c.map((childId: any) =>
          state.page.controls[childId]).map((c:any) => {
            return { ...c, children: getTemplateControls(state, c)};
          });
    }    

    return control.c.map((childId: any) => state.page.controls[childId])
      .filter((c: any) => c.t === 'columns').map((columns: any) =>
        columns.c.map((childId: any) => state.page.controls[childId]))
        .reduce((acc: any, columns: any) => ([...acc, ...columns])).map((cc: any) => {

          const template = getTemplateControls(state, cc);

          return {
            key: cc.i,
            name: cc.name,
            iconName: cc.icon,
            isIconOnly: cc.icononly === 'true',
            fieldName: cc.fieldname ? cc.fieldname.toLowerCase() : undefined,
            sortField: cc.sortfield ? cc.sortfield.toLowerCase() : cc.fieldname ? cc.fieldname.toLowerCase() : undefined,
            isResizable: cc.resizable === 'true',
            isSortable: cc.sortable,
            isSorted: cc.sorted === 'true' || cc.sorted === 'asc' || cc.sorted === 'desc',
            isSortedDescending: cc.sorted === 'desc',
            //isPadded: true,
            minWidth: cc.minwidth ? parseInt(cc.minwidth) : undefined,
            maxWidth: cc.maxwidth ? parseInt(cc.maxwidth) : undefined,
            onClick: cc.onclick === 'true',
            onColumnClick: _onColumnClick,
            columnActionsMode: cc.onclick === 'true' ||
              (cc.sortable !== undefined && cc.sortable !== 'false') ? ColumnActionsMode.clickable : ColumnActionsMode.disabled,
            template: template,
            onRender: template.length > 0 ? _columnOnRender : undefined
          }
        });
  }, shallowEqual);

  //console.log("columns:", columns);

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
    const key = (sortColumn as any).sortField;
    items = items.slice(0).sort((a: any, b: any) => {
      if ((sortColumn as any).isSortable === 'number') {
        return (sortColumn.isSortedDescending ? parseFloat(a[key]) < parseFloat(b[key]) : parseFloat(a[key]) > parseFloat(b[key])) ? 1 : -1;
      } else {
        return (sortColumn.isSortedDescending ? a[key] < b[key] : a[key] > b[key]) ? 1 : -1;
      }
    })
  }

  const shimmerLines = control.shimmerlines ? parseInt(control.shimmerlines) : 0;

  const gridProps: IShimmeredDetailsListProps = {
    columns: columns,
    items: items,
    enableShimmer: items.length === 0 && shimmerLines > 0,
    shimmerLines: shimmerLines,
    setKey: "set",
    compact: control.compact === 'true',
    isHeaderVisible: control.headervisible === 'false' ? false : true,
    onItemInvoked: _onItemInvoked,
    styles: {
      root: {
        //width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? defaultPixels(control.height) : undefined,
        padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined   
      }
    }
  };

  const _selection = new Selection({
    onSelectionChanged: () => {
      const ids = _selection.getSelection().map((elem: any) => (elem.i)).join(' ');
      console.log("onSelectionChanged:", ids);
      ws.pageEventFromWeb(control.i, 'select', ids)
    },
    getKey: (item: any) => item.i.toString()
  });

  // _selection.setItems(items);

  // selection mode
  gridProps.selectionMode = SelectionMode.none;
  if (control.selection === 'single' || control.selection === 'multiple') {
    gridProps.selectionMode = control.selection === 'single' ? SelectionMode.single : SelectionMode.multiple;
    gridProps.selection = _selection;
    gridProps.selectionPreservedOnEmptyClick = control.preserveselection === 'true';
  }

  //console.log("END RENDER GRID");
  // _selection.setChangeEvents(true, false);

  // <div style={{width: control.width !== undefined ? control.width : 'auto'}}>
  return <ShimmeredDetailsList {...gridProps} />;
})