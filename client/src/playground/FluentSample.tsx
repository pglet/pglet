import React from 'react';
import { Text, Link, FontWeights, PrimaryButton, DefaultButton, TextField } from '@fluentui/react';
import { Stack, IStackProps, ProgressIndicator, Icon } from '@fluentui/react';
//import { SharedColors, FontSizes } from '@uifabric/fluent-theme';

const boldStyle = { root: { fontWeight: FontWeights.semibold } };

const MyIcon = () => <Icon iconName="CompassNW" styles={{ root: { fontSize: '36px'}}} />;

// const stackTokens = { childrenGap: 50 };
// const stackStyles: Partial<IStackStyles> = { root: { width: 800 } };
const columnProps: Partial<IStackProps> = {
  tokens: { childrenGap: 15 },
  styles: { root: { width: "50%" } },
};

const ButtonType = DefaultButton;

export const FluentSample: React.FunctionComponent = () => {
  return (
    <Stack
      horizontalAlign="center"
      verticalAlign="center"
      verticalFill
      styles={{
        root: {
          width: '100%',
          padding: '10px'
        }
      }}
      tokens={
        { childrenGap: 15 }
      }
    >
      <PrimaryButton text="Button 1" />
      <Text variant="xxLarge">Hello, you!</Text>
      <Stack horizontal styles={{ root: { width: '50%' } }}>
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
          <Text variant="large" styles={{ root: { width: "200px", height: "100px" }}}>For a guide on how to customize this project, check out the UI Fabric documentation.</Text>
          <Text variant="large" styles={boldStyle}>
            Essential <b><i>Links</i></b>
          </Text>
        </Stack>
        <Stack {...columnProps}>
          <TextField label="First name" styles={{ root: { width: "200px", height: "100px" }}} />
          <TextField label="Last name" />
          <Stack horizontal tokens={{ childrenGap: 10 }}>
            <ButtonType text="Button 2" styles={{ root: { padding: "30px", margin: "30px 0" }}} />
            <PrimaryButton text="Button 3" styles={{ root: { width: "200px", height: "100px" }}} />
          </Stack>
        </Stack>
      </Stack>
      <Text block variant="xSmall" styles={{ root: { textAlign: "left", width: "50%" } }}>{ "Copyright (c) 2020" }</Text>
      <MyIcon/>
      <ProgressIndicator label="Example title" description="Example description" percentComplete={0.5} styles={{ root: { width: "50%" } }} />
    </Stack>
  );
};
