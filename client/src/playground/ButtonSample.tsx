import React from 'react';
import { Text, ContextualMenu, IContextualMenuItem } from '@fluentui/react';
import { MenuButton } from '@fluentui/react-button';
import { Stack } from '@fluentui/react';

let menu = {
  items: [
    { key: 'a', text: 'aaaaa' }
  ]
};

export const ButtonSample: React.FunctionComponent = () => {

  // onItemClick?: (ev?: React.MouseEvent<HTMLElement> | React.KeyboardEvent<HTMLElement>, item?: IContextualMenuItem) => boolean | void;

  const handleClick = (ev?: React.MouseEvent<HTMLElement> | React.KeyboardEvent<HTMLElement>, item?: IContextualMenuItem) => {
    console.log(item);
    return true;
  }

  return (
    <div>
      <Stack horizontalAlign="stretch">
        <Stack horizontal horizontalAlign='space-between'>
          <Text>Left</Text>
          <Stack horizontal>
            <MenuButton primary menu={<ContextualMenu {...menu} onItemClick={handleClick} />}>A</MenuButton>
          </Stack>
        </Stack>
      </Stack>
    </div>
  );
};
