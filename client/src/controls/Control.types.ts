export interface IControlProps {
    control: any;
    parentDisabled: boolean;
}

export interface IControlsListProps {
    controls: any;
    parentDisabled: boolean;
}

export interface IPageProps {
    pageName: string;
    control: any;
}

export interface ISigninOptions {
    gitHubEnabled: boolean;
    gitHubGroupScope: boolean;
    azureEnabled: boolean;
    azureGroupScope: boolean;
    googleEnabled: boolean;
    googleGroupScope: boolean;    
}

export interface ISigninProps {
    signinOptions: ISigninOptions;
    onDismiss?: () => any;
}