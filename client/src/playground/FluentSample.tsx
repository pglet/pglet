import React from 'react';
import { Text, Link, FontWeights, TextField, FontIcon, ContextualMenu } from '@fluentui/react';
import { Button, MenuButton } from '@fluentui/react-button';
import { Stack, IStackProps, ProgressIndicator, mergeStyles } from '@fluentui/react';
import { DatePicker } from '@fluentui/react-date-time';

const boldStyle = { root: { fontWeight: FontWeights.semibold } };

const iconClass = mergeStyles({
  fontSize: 100,
  height: 100,
  width: 100,
  margin: 0,
  color: 'salmon'
});

const MyIcon = () => <FontIcon iconName="Upload" className={iconClass} />;

// const stackTokens = { childrenGap: 50 };
// const stackStyles: Partial<IStackStyles> = { root: { width: 800 } };
const columnProps: Partial<IStackProps> = {
  tokens: { childrenGap: 15 },
  styles: { root: { width: "50%" } },
};

let menu = {
  items: [
    { key: 'a', text: 'aaaaa' }
  ]
};

export const FluentSample: React.FunctionComponent = () => {
  return (
    <div>
      <Stack horizontalAlign="start">
        <DatePicker label="Select date" isRequired={true} placeholder="Select a date..." />
      </Stack>
      <Stack horizontalAlign="stretch">
        <Stack horizontal horizontalAlign='space-between'>
          <Text>Left</Text>
          <Stack horizontal>
            <Text><MenuButton primary menu={<ContextualMenu {...menu} />}>A</MenuButton></Text>
            <Text><Button content="B"/></Text>
          </Stack>
        </Stack>
      </Stack>
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
        <Button primary content="Button 1" />
        <Stack horizontal>
          <Text variant="large" styles={{
            root: {
              whiteSpace: 'pre'
            }
          }}>1{"\n"}2{"\n"}3{"\n"}</Text>
          <Text variant="xSmall" styles={{
            root: {
              width: '40px',
              height: '40px',
              backgroundColor: 'blue',
              verticalAlign: 'middle',
              textAlign: 'center',
              display: 'inline-flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: '#fff',
              borderColor: 'black',
              borderWidth: '1px',
              borderStyle: 'solid',
              borderRadius: '20px'
            }
          }}>2</Text>
          <Text variant="xSmall" styles={{
            root: {
              width: '40px',
              height: '40px',
              backgroundColor: 'blue',
              verticalAlign: 'middle',
              textAlign: 'center',
              display: 'inline-flex',
              alignItems: 'center',
              justifyContent: 'center',
              color: '#fff',
              borderColor: 'black',
              borderWidth: '1px',
              borderStyle: 'solid',
              borderRadius: '20px'
            }
          }}>3</Text>
        </Stack>

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
            <Text variant="large" styles={{ root: { width: "200px", height: "100px" } }}>For a guide on how to customize this project, check out the UI Fabric documentation.</Text>
            <Text variant="large" styles={boldStyle}>
              Essential <b><i>Links</i></b>
            </Text>
          </Stack>
          <Stack {...columnProps}>
            <TextField label="First name" styles={{ root: { width: "200px", height: "100px" } }} />
            <TextField label="Last name" />
            <Stack horizontal tokens={{ childrenGap: 10 }}>
              <Button content="Button 2" iconProps={{ iconName: "Installation" }} />
              <Button content="Button 3" iconProps={{ iconName: "Filter" }} styles={{ root: { width: "200px", height: "100px" } }} />
            </Stack>
          </Stack>
        </Stack>
        <Text block variant="xSmall" styles={{ root: { textAlign: "left", width: "50%" } }}>{"Copyright (c) 2020"}</Text>
        <MyIcon />
        <ProgressIndicator label="Example title" description="Example description" percentComplete={0.5} styles={{ root: { width: "50%" } }} />
      </Stack>
    </div>
  );
};
