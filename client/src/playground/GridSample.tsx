import React from 'react';
import { DetailsList, IColumn, DetailsListLayoutMode, Selection, Link } from '@fluentui/react';

export interface IDocument {
  key: string;
  name: string;
  iconName: string;
}

let items: IDocument[] = [
  { key: 'item1', name: 'Item 1', iconName: 'Icon 1' },
  { key: 'item2', name: 'Item 2', iconName: 'Icon 2' },
  { key: 'item3', name: 'Item 3', iconName: 'Icon 3' }
];

export const GridSample: React.FunctionComponent = () => {

  // onItemClick?: (ev?: React.MouseEvent<HTMLElement> | React.KeyboardEvent<HTMLElement>, item?: IContextualMenuItem) => boolean | void;

  // const handleClick = (ev?: React.MouseEvent<HTMLElement> | React.KeyboardEvent<HTMLElement>, item?: IContextualMenuItem) => {
  //   console.log(item);
  //   return true;
  // }

  const _onItemInvoked = (item: any) => {
    alert(`Item invoked: ${item.name}`);
  }

  const _onColumnClick = (ev: React.MouseEvent<HTMLElement>, column: IColumn): void => {
    console.log(column)
  }

  let _lastResizedColumnTimeout: any = null;  

  const _onColumnResize = (column?: IColumn, newWidth?: number, columnIndex?: number) => {
    if (_lastResizedColumnTimeout != null) {
      clearTimeout(_lastResizedColumnTimeout);
    }
    _lastResizedColumnTimeout = setTimeout(() => {
      console.log(column);
    }, 500);
  }

  const _selection = new Selection({
    onSelectionChanged: () => {
      console.log(_selection.getSelection());
    },
  });

  const columns: IColumn[] = [
    {
      key: 'column1',
      name: 'File Type',
      iconName: 'Page',
      isIconOnly: true,
      fieldName: 'name',
      isResizable: true,
      isSorted: false,
      minWidth: 80,
      maxWidth: 200,
      onColumnClick: _onColumnClick
    },
    {
      key: 'column2',
      name: 'Name',
      fieldName: 'name',
      minWidth: 210,
      isRowHeader: true,
      isResizable: true,
      isSorted: true,
      isSortedDescending: false,
      sortAscendingAriaLabel: 'Sorted A to Z',
      sortDescendingAriaLabel: 'Sorted Z to A',
      onColumnClick: _onColumnClick,
      data: 'string',
      isPadded: true,
    },
    {
      key: 'column3',
      name: 'Action',
      fieldName: 'name',
      minWidth: 30,
      maxWidth: 30,
      onColumnClick: _onColumnClick,
      data: 'string',
      isPadded: true,
      onRender: (item: IDocument) => {
        return <Link href={'https://' + item.name}>{item.name}</Link>;
      },
    }
  ];

  return (
    <div>
      <DetailsList
        items={items}
        compact={true}
        columns={columns}
        onColumnResize={_onColumnResize}
        selection={_selection}
        selectionPreservedOnEmptyClick={true}
        //selectionMode={SelectionMode.none}
        //getKey={this._getKey}
        setKey="key"
        layoutMode={DetailsListLayoutMode.justified}
        isHeaderVisible={true}
        onItemInvoked={_onItemInvoked}
      />
    </div>
  );
};
