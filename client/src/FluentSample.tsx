import React from 'react';
import { Text, Link, FontWeights, PrimaryButton, TextField } from 'office-ui-fabric-react';
import { Stack, IStackProps, IStackStyles } from 'office-ui-fabric-react/lib/Stack';

const boldStyle = { root: { fontWeight: FontWeights.semibold } };

const stackTokens = { childrenGap: 50 };
const stackStyles: Partial<IStackStyles> = { root: { width: 800 } };
const columnProps: Partial<IStackProps> = {
  tokens: { childrenGap: 15 },
  styles: { root: { width: "50%" } },
};

export const FluentSample: React.FunctionComponent = () => {
  return (
    <Stack
      horizontal
      horizontalAlign="center"
      verticalAlign="center"
      verticalFill
      styles={{
        root: {
          width: '100%',
          padding: '10px'
        }
      }}
      gap={15}
    >
      <Stack horizontal styles={{ root: { width: '100%'}}}>
        <Stack {...columnProps}>
          <Link href="https://developer.microsoft.com/en-us/fabric">Docs</Link>
          <Link href="https://stackoverflow.com/questions/tagged/office-ui-fabric">Stack Overflow</Link>
          <Link href="https://github.com/officeDev/office-ui-fabric-react/">Github</Link>
          <Link href="https://twitter.com/officeuifabric">Twitter</Link>
          {/* <img
        src="https://raw.githubusercontent.com/Microsoft/just/master/packages/just-stack-uifabric/template/src/components/fabric.png"
        alt="logo"
      /> */}
          <Text variant="xxLarge" styles={boldStyle}>
            Welcome to Your UI Fabric App
      </Text>
          <Text variant="large">For a guide on how to customize this project, check out the UI Fabric documentation.</Text>
          <Text variant="large" styles={boldStyle}>
            Essential Links
      </Text>
        </Stack>
        <Stack {...columnProps}>
          <TextField label="First name" />
          <TextField label="Last name" />
          <Stack horizontal gap="10">
            <PrimaryButton text="Button 2" />
            <PrimaryButton text="Button 3" />
          </Stack>
        </Stack>
      </Stack>
    </Stack>
  );
};
