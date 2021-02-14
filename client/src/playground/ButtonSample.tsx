import React from 'react';
import { Text, ContextualMenu, IContextualMenuItem } from '@fluentui/react';
import { MenuButton } from '@fluentui/react-button';
import { Stack } from '@fluentui/react';
import {Prism as SyntaxHighlighter} from 'react-syntax-highlighter'
import {dark} from 'react-syntax-highlighter/dist/esm/styles/prism'

const ReactMarkdown = require('react-markdown')
const gfm = require('remark-gfm')


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

  const markdown = `

A paragraph with *emphasis* and **strong importance**.

> A block quote with ~strikethrough~ and a URL: https://reactjs.org.
  
* Lists
* [ ] todo
* [x] done
  
A table:
  
| a | b |
| - | - |

## HTML block below

<blockquote>
  This blockquote will change based on the HTML settings above.
</blockquote>

## How about some code?

\`\`\`js
var React = require('react');
var Markdown = require('react-markdown');

React.render(
  <Markdown source="# Your markdown here" />,
  document.getElementById('content')
);
\`\`\`

Pretty neat, eh?`

const renderers = {
  code: (code: any) => {
    return <SyntaxHighlighter style={dark} language={code.language} children={code.value} />
  }
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
      <Stack>
        <Text>
          <ReactMarkdown plugins={[gfm]} renderers={renderers} children={markdown} />
        </Text>
      </Stack>
    </div>
  );
};
