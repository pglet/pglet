import React from 'react';
import { Sticky, StickyPositionType, ScrollablePane, mergeStyles, TextField, Stack, Text, DetailsList, IColumn, DetailsListLayoutMode, Selection, Link } from '@fluentui/react';

export interface IDocument {
  key: string;
  name: string;
  iconName: string;
}

let items: IDocument[] = [];

for(let i = 0; i < 1000; i++) {
  items.push({ key: `item${i}`, name: `Item ${i}`, iconName: `Icon ${i}` });
}

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
      onRender: (item: IDocument) => {
        return <Link href={'https://' + item.name}>{item.name}</Link>;
      },
    }
  ];

  // const className1 = mergeStyles({
  //   width: '30%'
  // });

  const className2 = mergeStyles({
    //width: '50%'
  });  

  // const className = mergeStyles(blueBackgroundClassName, {
  //   width: 'auto',
  //   selectors: {
  //     '& .ms-ViewPort': {
  //       backgroundColor: 'red',
  //       //width: '100%'
  //     }
  //   }
  // });

  // const classNames = mergeStyleSets({
  //   table1: {
  //     'test': {
  //       width: '100'
  //     },
  //     margin: 'auto',
  //   }
  // });

  return (
    <div>
      <ScrollablePane>
      <Stack horizontal horizontalAlign="stretch">
        <Stack.Item grow={1}>
          <Text>Left menu</Text>
        </Stack.Item>
        <Stack.Item grow={2}>
          <TextField styles={{ root: { width: '100%' } }} />
        </Stack.Item>
      </Stack>        
        <Stack horizontal horizontalAlign="stretch">
          <Stack.Item grow>
            <Sticky stickyPosition={StickyPositionType.Header}>
              <Text>Nav menu</Text>
            </Sticky>
          </Stack.Item>
          <Stack.Item grow className={className2}>
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
            {/* <Text>Center</Text> */}
          </Stack.Item>
        </Stack>
      </ScrollablePane>
    </div>
  );
};
