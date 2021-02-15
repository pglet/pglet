import React from 'react';
import { Text } from '@fluentui/react';
import { Stack } from '@fluentui/react';
//import { Light as SyntaxHighlighter } from 'react-syntax-highlighter';
//import { github } from 'react-syntax-highlighter/dist/esm/styles/hljs';

import ReactMarkdown from 'react-markdown'
import gfm from 'remark-gfm'

// import bash from 'react-syntax-highlighter/dist/esm/languages/hljs/bash';
// import powershell from 'react-syntax-highlighter/dist/esm/languages/hljs/powershell';
// import python from 'react-syntax-highlighter/dist/esm/languages/hljs/python';
// import js from 'react-syntax-highlighter/dist/esm/languages/hljs/javascript';

// SyntaxHighlighter.registerLanguage('javascript', js);
// SyntaxHighlighter.registerLanguage('powershell', powershell);
// SyntaxHighlighter.registerLanguage('python', python);
// SyntaxHighlighter.registerLanguage('bash', bash);

export const ButtonSample: React.FunctionComponent = () => {

  const markdown = `

A paragraph with *emphasis* and **strong importance**.

> A block quote with ~strikethrough~ and a URL: https://reactjs.org.
  
* Lists
* [ ] todo
* [x] done
  
A table:
  
| a | b |
| - | - |

# Heading 1

## Heading 2

### Heading 3

#### Heading 4

##### Heading 5

<blockquote>
  This blockquote will change based on the HTML settings above.
</blockquote>

## How about some code?

\`\`\`javascript
var React = require('react');
var Markdown = require('react-markdown');

React.render(
  <Markdown source="# Your markdown here" />,
  document.getElementById('content')
);
\`\`\`

\`\`\`powershell
for($i = 20; $i -lt 41; $i++) {
  Invoke-Pglet "setf rdp y=$i"
  Start-Sleep -ms 100
}
\`\`\`

Pretty neat, eh?`

// const renderers = {
//   code: (code: any) => {
//     return <SyntaxHighlighter style={github} language={code.language} children={code.value} />
//   }
// }

  return (
    <div>
      <Stack horizontalAlign="stretch">
        <Stack horizontal horizontalAlign='space-between'>
          <Text>Left</Text>
          <Stack horizontal>
          </Stack>
        </Stack>
      </Stack>
      <Stack>
        <Text>
          <ReactMarkdown plugins={[gfm]} children={markdown} />
        </Text>    
      </Stack>
    </div>
  );
};
