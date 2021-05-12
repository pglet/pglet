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

export interface ILoginOptions {
    gitHubEnabled: boolean;
    gitHubGroupScope: boolean;
    azureEnabled: boolean;
    azureGroupScope: boolean;
}

export interface ILoginProps {
    loginOptions: ILoginOptions
}