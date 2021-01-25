import React from 'react';
import { BrowserRouter as Router, Switch, Route } from "react-router-dom";
import { PageLanding } from './controls/PageLanding'
import { AccountLanding } from './controls/AccountLanding';
import { FluentSample } from './playground/FluentSample';
import { ButtonSample } from './playground/ButtonSample';

export const App: React.FunctionComponent = () => {
  return (
    <Router>
      <Switch>
        <Route path="/:accountName/:pageName" children={<PageLanding />} />
        <Route path="/sample" children={<FluentSample />} />
        <Route path="/button-sample" children={<ButtonSample />} />
        <Route path="/:accountName" children={<AccountLanding />} />
        <Route path="/" children={<PageLanding />} />
      </Switch>
    </Router>
  );
};
