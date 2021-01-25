import React from 'react';
import { Text, ContextualMenu } from '@fluentui/react';
import { MenuButton } from '@fluentui/react-button';
import { Stack } from '@fluentui/react';

let menu = {
  items: [
    { key: 'a', text: 'aaaaa' }
  ]
};

export const ButtonSample: React.FunctionComponent = () => {
  return (
    <div>
      <Stack horizontalAlign="stretch">
        <Stack horizontal horizontalAlign='space-between'>
          <Text>Left</Text>
          <Stack horizontal>
            <MenuButton primary menu={<ContextualMenu {...menu} />}>A</MenuButton>
          </Stack>
        </Stack>
      </Stack>
    </div>
  );
};
